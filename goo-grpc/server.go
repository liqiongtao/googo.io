package goo_grpc

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/facebookgo/grace/gracenet"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_pprof "github.com/liqiongtao/googo.io/goo-pprof"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"github.com/liqiongtao/googo.io/goocontext"
	"google.golang.org/grpc"
)

type Server struct {
	conf Config
	opts serverOptions

	*gracenet.Net
	*grpc.Server

	lis net.Listener

	stopOnce  sync.Once
	hooksOnce sync.Once
}

// restarting 防止连发 SIGHUP 重复 StartProcess
var restarting int32

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
		//grpc.UnknownServiceHandler(func(srv interface{}, stream grpc.ServerStream) error {
		//	return nil
		//}),
	}...)

	return &Server{
		conf:   conf,
		opts:   opts,
		Net:    &gracenet.Net{},
		Server: grpc.NewServer(serverOptions...),
	}
}

func (s *Server) Serve() (err error) {
	defer func() {
		if r := recover(); r != nil {
			goo_log.WithTag("goo-grpc").Error(r)
		}
	}()

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

	s.lis, err = s.Net.Listen("tcp", s.conf.Addr)
	if err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return
	}

	// 服务注册
	goo_utils.AsyncFunc(func() {
		if !s.opts.Register2Etcd {
			return
		}

		cli := s.opts.EtcdClient
		if cli == nil {
			goo_log.WithTag("goo-grpc").Error("no etcd client")
			return
		}

		address := s.lis.Addr().String()
		if s.conf.ServiceEndpoint != "" {
			index := strings.LastIndex(address, ":")
			cli.RegisterService(s.conf.ServiceName, fmt.Sprintf("%s:%s", s.conf.ServiceEndpoint, address[index+1:]))
		} else {
			cli.RegisterService(s.conf.ServiceName, address)
		}
	})

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

	s.storePID()

	// 继承 listener 的子进程就绪后，延迟通知父进程优雅退出
	goocontext.NotifyParentExitAfter(300 * time.Millisecond)

	select {
	case <-goocontext.Root().Done():
	case serveErr := <-serveErrCh:
		// Serve 立刻失败时触发退出钩子，避免 AsyncFunc(Serve)+<-Root().Done() 假活
		err = serveErr
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
		<-goocontext.Root().Done()
	}

	return
}

func (s *Server) registerSignalHooks() {
	s.hooksOnce.Do(func() {
		goo_pprof.RegisterSignal()

		goocontext.OnRestart(func() {
			// 防止连发 SIGHUP 重复 fork；失败则允许重试
			if !atomic.CompareAndSwapInt32(&restarting, 0, 1) {
				return
			}
			if _, err := s.Net.StartProcess(); err != nil {
				atomic.StoreInt32(&restarting, 0)
				goo_log.WithTag("goo-grpc").Error(err)
				return
			}
			goo_log.WithTag("goo-grpc").Warn("服务重启")
			// handoff 失败时父进程仍存活，超时后允许再次热重启
			time.AfterFunc(10*time.Second, func() {
				atomic.StoreInt32(&restarting, 0)
			})
		})
		goocontext.OnExit(func() {
			goo_pprof.StopDefault()
			s.gracefulStop()
			goo_log.WithTag("goo-grpc").Warn("服务退出")
		})
	})
}

// 平滑退出
func (s *Server) gracefulStop() {
	s.stopOnce.Do(func() {
		s.Server.GracefulStop()
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
	goo_log.WithTag("goo-grpc").DebugF(fmt.Sprintf("server is running, address=%s, pid=%s", s.lis.Addr(), pid))
}
