package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/artifactsign"
	"example.com/artifactsign/internal/api"
)

func main() {
	addr := flag.String("addr", ":8099", "HTTP 监听地址")
	web := flag.String("web", "web", "静态管理页目录")
	persist := flag.String("persist", "", "持久化目录（可选）")
	flag.Parse()

	opts := []artifactsign.Option{}
	if *persist != "" {
		opts = append(opts, artifactsign.WithPersistDir(*persist))
	}
	svc := artifactsign.New(opts...)
	defer svc.Close()
	if *persist != "" {
		_ = svc.LoadPersist()
	}

	srv := api.New(svc, api.Options{WebDir: *web, AllowCORS: true})
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("signd 监听 %s", *addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	_ = httpSrv.Close()
}
