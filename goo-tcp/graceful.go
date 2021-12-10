package goo

import (
	"fmt"
	"github.com/facebookgo/grace/gracenet"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type Graceful struct {
	addr    string
	net     *gracenet.Net
	handler func(net.Conn)
	wg      sync.WaitGroup
}

func New(addr string, handler func(net.Conn)) *Graceful {
	return &Graceful{addr: addr, handler: handler, net: &gracenet.Net{}}
}

func (g *Graceful) Serve() {
	addr, err := net.ResolveTCPAddr("tcp", g.addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	l, err := g.net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	quit := make(chan struct{})

	goo_utils.AsyncFunc(g.killPPID)
	goo_utils.AsyncFunc(g.storePID)
	goo_utils.AsyncFunc(g.handleSignal(l, quit))

	for {
		conn, err := l.Accept()
		if err != nil {
			goo_log.WithTag("goo-tcp").Error(err)
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			continue
		}
		g.wg.Add(1)
		goo_utils.AsyncFunc(func() {
			defer g.wg.Done()
			defer conn.Close()
			go g.handler(conn)
			<-quit
		})
	}

	g.wg.Wait()
}

func (g *Graceful) handleSignal(l *net.TCPListener, quit chan struct{}) func() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	return func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				signal.Stop(ch)
				l.Close()
				close(quit)
				return
			case syscall.SIGUSR1, syscall.SIGUSR2:
				if _, err := g.net.StartProcess(); err != nil {
					goo_log.WithTag("goo-tcp").Error(err)
				}
			}
		}
	}
}

func (g *Graceful) storePID() {
	pid := fmt.Sprintf("%d", os.Getpid())
	ioutil.WriteFile(".pid", []byte(pid), 0644)
	log.Println(fmt.Sprintf("server is running, address=%s, pid=%s", g.addr, pid))
}

func (g *Graceful) killPPID() {
	inherit := os.Getenv("LISTEN_FDS") != ""
	if !inherit {
		return
	}
	ppid := os.Getppid()
	if ppid == 1 {
		return
	}
	syscall.Kill(ppid, syscall.SIGTERM)
}
