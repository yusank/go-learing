package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	sigChan := make(chan os.Signal)

	// 监听信号
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// 信号chan有数据 执行cancel
		<-sigChan
		cancel()
	}()

	var g *errgroup.Group
	g, ctx = errgroup.WithContext(ctx)

	// 实现一个简单的server group 统一启动和关闭
	sg := newServerGroup(g)
	sg.addServer(newServer(":8081"))
	sg.addServer(newServer(":8082"))
	sg.addServer(newServer(":8083"))
	sg.start()

	// 模拟某个 server 突然挂了
	g.Go(func() error {
		time.Sleep(time.Second * 5)
		return fmt.Errorf("custom error")
	})

	// 任意一个errGroup 的 Go 返回错误了 ctx 的cancel 方法会被执行
	// sigChan 收到退出信号也会执行 cancel
	//
	// http server 不会存在都执行完正常退出的情况（整个生命周期所有server都会阻塞）
	// 所以不需要通过 g.Wait() 来阻塞主线程
	<-ctx.Done()
	sg.stop()
}

type serverGroup struct {
	g        *errgroup.Group
	srvSlice []*server
}

func newServerGroup(g *errgroup.Group) *serverGroup {
	return &serverGroup{
		g:        g,
		srvSlice: make([]*server, 0),
	}
}

func (sg *serverGroup) addServer(s *server) {
	sg.srvSlice = append(sg.srvSlice, s)
}

func (sg *serverGroup) start() {
	for _, s := range sg.srvSlice {
		s.start(sg.g)
	}
}

func (sg *serverGroup) stop() {
	for _, s := range sg.srvSlice {
		s.stop()
	}
}

type server struct {
	addr string
	hs   *http.Server
}

func newServer(addr string) *server {
	return &server{
		addr: addr,
		hs: &http.Server{
			Addr: addr,
		},
	}
}

func (s *server) start(g *errgroup.Group) {
	g.Go(func() error {
		log.Println("listen:", s.addr)
		return s.hs.ListenAndServe()
	})
}

func (s *server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 假设500毫秒才能处理完退出
	time.Sleep(time.Millisecond * 500)
	log.Println("shutdown:", s.addr)
	s.hs.Shutdown(ctx)
}
