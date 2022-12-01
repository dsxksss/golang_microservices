package main

import (
	"context"
	"go_sc_small_server/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// 创建记录器对象
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// 初始化路由
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	// 创建多路复用器
	sm := http.NewServeMux()

	// 注册路由
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	// 服务配置
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	// 服务监听
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// 关闭服务后释放资源等功能的代码(不是很清楚)
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	tc, err := context.WithTimeout(context.Background(), 30*time.Second)
	if err != nil {
		l.Fatal(err)
	}
	s.Shutdown(tc)
}
