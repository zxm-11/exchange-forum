package main

import (
	"context"
	"exchangeapp/config"
	"exchangeapp/consumers"
	"exchangeapp/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	config.InitConfig()
	r := router.SetupRouter()

	go config.ReconnectRabbitMQ() //启动重连RabbitMQ的goroutine

	go func() {
		prefetch := config.AppConfig.RabbitMQ.Likeprefetch
		interval := time.Duration(config.AppConfig.RabbitMQ.Liketaskinterval) * time.Second
		consumers.LikeConsumer(prefetch, interval)
	}()

	port := config.AppConfig.App.Port

	if port == "" {
		port = ":8080"
	}
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); //启动HTTP服务器,监听指定的端口并处理传入的HTTP请求
		err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)   //创建一个信号通道,用于接收操作系统信号
	signal.Notify(quit, os.Interrupt) //监听操作系统的中断信号,当用户按下Ctrl+C时会触发这个信号
	<-quit                            //阻塞主线程,直到接收到中断信号
	log.Println("Shutdown Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil { //关闭HTTP服务器,等待所有当前处理的请求完成
		log.Println("Server Shutdown:", err)
	}

	config.CloseRabbitMQ()
	log.Println("Server exiting")
}
