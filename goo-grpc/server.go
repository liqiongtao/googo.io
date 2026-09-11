package goo_grpc

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_pprof "github.com/liqiongtao/googo.io/goo-pprof"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"github.com/liqiongtao/googo.io/goocontext"
	"google.golang.org/grpc"
)

type Server struct {
	conf Config
	opts serverOptions

	*grpc.Server
	upg *tableflip.Upgrader

	lis net.Listener

	stopOnce  sync.Once
	restartMu sync.Mutex
	hooksOnce sync.Once
}

func New(conf Config, opt ...ServerOption) *Server {
	opts := newDefaultServerOptions(conf)
	for _, o := range opt {
		o.apply(&opts)
	}

	serverOptions := append(opts.ServerOptions, []grpc.ServerOption{
		grpc.MaxRecvMsgSize(MaxRecvMsgSize),
		grpc.MaxSendMsgSize(MaxSendMsgSize),
		// 单向拦截 - 链式
		grpc.ChainUnaryInterceptor(
			serverUnaryInterceptorLog(opts.NoLogMethods),
			serverUnaryInterceptorRecovery(),
			serverUnaryInterceptorAuth(opts.AuthFunc),
		),
		// 流式拦截 - 链式
		grpc.ChainStreamInterceptor(
			serverStreamInterceptorLog(opts.NoLogMethods),
			serverStreamInterceptorRecovery(),
			serverStreamInterceptorAuth(opts.AuthFunc),
		),
		// todo:: 服务未找到
		//grpc.UnknownServiceHandler(func(srv any, stream grpc.ServerStream) error {
		//	return nil
		//}),
	}...)

	return &Server{
		conf:   conf,
		opts:   opts,
		Server: grpc.NewServer(serverOptions...),
	}
}

func (s *Server) Serve() (err error) {
	defer func() {
		if r := recover(); r != nil {
			goo_log.WithTag("goo-grpc").Error(r)
		}
	}()

	s.upg, err = tableflip.New(tableflip.Options{})
	if err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return
	}
	defer s.upg.Stop()

	// 本机内网IP
	if s.conf.ServiceEndpoint == "" || s.conf.Addr == "" {
		var localIp string
		localIp, err = goo_utils.LocalIP()
		if err != nil {
			goo_log.WithTag("goo-grpc").Error(err)
			return
		}

		if s.conf.ServiceEndpoint == "" {
			s.conf.ServiceEndpoint = localIp
		}
		if s.conf.Addr == "" {
			s.conf.Addr = fmt.Sprintf("%s:0", localIp)
		}
	}

	// 随机端口
	if !strings.Contains(s.conf.Addr, ":") {
		s.conf.Addr += ":0"
	}

	s.lis, err = s.upg.Listen("tcp", s.conf.Addr)
	if err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return
	}

	// 先注册钩子，再 Serve，缩小空钩子 cancel 窗口
	s.registerSignalHooks()

	serveErrCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				goo_log.WithTag("goo-grpc").Error(r)
			}
		}()

		if serveErr := s.Server.Serve(s.lis); serveErr != nil {
			goo_log.WithTag("goo-grpc").Error(serveErr)
			serveErrCh <- serveErr
		}
	}()

	// Ready 前完成首次 etcd 注册，避免服务发现空窗
	if s.opts.Register2Etcd {
		cli := s.opts.EtcdClient
		if cli == nil {
			err = fmt.Errorf("Register2Etcd enabled but etcd client is nil")
			goo_log.WithTag("goo-grpc").Error(err)
			s.gracefulStop()
			return
		}
		address := s.lis.Addr().String()
		regAddr := address
		if s.conf.ServiceEndpoint != "" {
			_, port, perr := net.SplitHostPort(address)
			if perr != nil {
				err = fmt.Errorf("parse listen addr %q: %w", address, perr)
				goo_log.WithTag("goo-grpc").Error(err)
				s.gracefulStop()
				return
			}
			regAddr = net.JoinHostPort(s.conf.ServiceEndpoint, port)
		}
		if err = cli.RegisterServiceTimeout(s.conf.ServiceName, regAddr, 30*time.Second); err != nil {
			goo_log.WithTag("goo-grpc").Error(err)
			s.gracefulStop()
			return
		}
	}

	s.storePID()

	// 通知父进程：本进程已就绪可接管
	if readyErr := s.upg.Ready(); readyErr != nil {
		goo_log.WithTag("goo-grpc").Error(readyErr)
	}

	// 父进程在子进程 Ready 后会收到 Exit，走统一优雅退出路径
	go func() {
		<-s.upg.Exit()
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()

	select {
	case <-goocontext.Root().Done():
	case serveErr := <-serveErrCh:
		// Serve 立刻失败时触发退出钩子，避免假活
		err = serveErr
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}
	goocontext.Wait()

	return
}

func (s *Server) registerSignalHooks() {
	s.hooksOnce.Do(func() {
		goo_pprof.RegisterSignal()

		goocontext.OnRestart(func() {
			if s.upg == nil {
				return
			}
			if !s.restartMu.TryLock() {
				return
			}
			defer s.restartMu.Unlock()
			if err := s.upg.Upgrade(); err != nil {
				goo_log.WithTag("goo-grpc").Error(err)
				return
			}
			goo_log.WithTag("goo-grpc").Warn("服务重启")
		})
		goocontext.OnExit(func() {
			goo_pprof.StopDefault()
			s.gracefulStop()
			goo_log.WithTag("goo-grpc").Warn("服务退出")
		})
	})
}

// 平滑退出（最长等 30s，超时强制 Stop，避免 OnExit/Wait 永久卡住）
func (s *Server) gracefulStop() {
	s.stopOnce.Do(func() {
		done := make(chan struct{})
		go func() {
			s.Server.GracefulStop()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(30 * time.Second):
			goo_log.WithTag("goo-grpc").Warn("GracefulStop 超时，强制 Stop")
			s.Server.Stop()
		}
		if s.lis != nil {
			if err := s.lis.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
				goo_log.WithTag("goo-grpc").Error(err)
			}
		}
	})
}

func (s *Server) storePID() {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := os.WriteFile(".pid", []byte(pid), 0644); err != nil {
		goo_log.WithTag("goo-grpc").Error(fmt.Sprintf("server store pid err: %s", err.Error()))
		return
	}
	goo_log.WithTag("goo-grpc").DebugF("server is running, address=%s, pid=%s", s.lis.Addr(), pid)
}
