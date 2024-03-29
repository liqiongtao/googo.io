package goo_grpc

import (
	"fmt"
	"github.com/facebookgo/grace/gracenet"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Server struct {
	conf Config

	*gracenet.Net
	*grpc.Server

	lis net.Listener

	pprof *PProf
}

func New(conf Config, opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	serverOptions := append(opts.ServerOptions, []grpc.ServerOption{
		// 单向拦截 - 链式
		grpc.ChainUnaryInterceptor(
			serverUnaryInterceptorLog(),
			serverUnaryInterceptorRecovery(),
			serverUnaryInterceptorAuth(opts.AuthFunc),
		),
		// 流式拦截 - 链式
		grpc.ChainStreamInterceptor(
			serverStreamInterceptorLog(),
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

	// 随机端口
	addr := s.conf.Addr
	if addr == "" {
		addr = "0.0.0.0:0"
	}
	if !strings.Contains(addr, ":") {
		addr += ":0"
	}

	s.lis, err = s.Net.Listen("tcp", addr)
	if err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
		return
	}

	// 服务注册
	if defaultServerOptions.Register2Etcd {
		if cli := defaultServerOptions.EtcdClient; cli != nil {
			cli.RegisterService(s.conf.ServiceName, s.lis.Addr().String())
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				goo_log.WithTag("goo-grpc").Error(r)
			}
		}()

		if err = s.Server.Serve(s.lis); err != nil {
			goo_log.WithTag("goo-grpc").Error(err)
		}
	}()

	s.storePID()
	s.handleSignal()

	return
}

func (s *Server) handleSignal() {
	ch := make(chan os.Signal)

	signal.Notify(ch, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)

	for sig := range ch {
		switch sig {
		case syscall.SIGUSR1: // kill -USR1
			s.pprofStart()

		case syscall.SIGUSR2: // kill -USR2
			s.pprofStop()

		case syscall.SIGHUP: // kill -1
			s.gracefulReStart()
			return

		case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL: // kill -9 or ctrl+c
			s.gracefulStop()
			return
		}
	}
}

// 开启分析监控
func (s *Server) pprofStart() {
	if s.pprof == nil {
		s.pprof = newPProf()
	}
	s.pprof.start()
	log.Println("pprof running")
}

// 停止分析监控
func (s *Server) pprofStop() {
	if s.pprof != nil {
		s.pprof.stop()
		log.Println("pprof stopped, dump files:", s.pprof.cpuFile, s.pprof.memoryFile,
			s.pprof.goroutineFile, s.pprof.mutexFile, s.pprof.blockFile)
	}
	s.pprof = nil
}

// 平滑重启
func (s *Server) gracefulReStart() {
	if _, err := s.Net.StartProcess(); err != nil {
		goo_log.WithTag("goo-grpc").Error(err)
	}
}

// 平滑退出
func (s *Server) gracefulStop() {
	s.Server.GracefulStop()
}

func (s *Server) storePID() {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(".pid", []byte(pid), 0644); err != nil {
		log.Println(fmt.Sprintf("server store pid err: %s", err.Error()))
		return
	}
	log.Println(fmt.Sprintf("server is running, address=%s, pid=%s", s.lis.Addr(), pid))
}
