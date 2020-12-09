package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	fmt.Printf("1111")
	stop, cancel := start()
	group, _ := errgroup.WithContext(stop)
	group.Go(func() error {
		defer cancel()
		return httpServer(stop)
	})
	group.Go(func() (err error) {
		defer cancel()
		return signalServer(stop)
	})

	// 两个都挂了之后，运行收尾shutdown()
	if err := group.Wait(); err != nil {
		shutdown()
	}
}

func start() (context.Context, func()) {
	stop, cancel := context.WithCancel(context.Background())
	var once sync.Once
	c := func() {
		once.Do(cancel)
	}
	return stop, c
}

func httpServer(ctx context.Context) (err error) {
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	// 模拟一个失败的情况
	http.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "服务器挂了")
		server.Shutdown(context.Background())
	})
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()
	err = server.ListenAndServe()
	log.Printf("http server stop: %v", err)
	return nil
}

func signalServer(ctx context.Context) error {
	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-c:
		log.Println("接收到信号")
		return fmt.Errorf("关闭信号")
	case <-ctx.Done():
		log.Println("关闭所有信号")
		return nil
	}
}

func shutdown() {
	down, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	finish := make(chan struct{})
	go func() {
		time.Sleep(time.Microsecond * 200)
		<-finish
	}()
	select {
	case <-finish:
		return
	case <-down.Done():
		return
	}
}
