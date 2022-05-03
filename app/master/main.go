package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hubogle/Crontab/app/master/config"
	"github.com/hubogle/Crontab/app/master/initialize"
	"github.com/hubogle/Crontab/app/master/router"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	initialize.InitLogger()
	initialize.InitConfig()
	_ = initialize.InitTrans()
}
func main() {
	var err error
	cfg := config.GetConfig()
	gin.SetMode(cfg.App.RunMode)
	mux := router.NewHTTPServer()
	src := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			cfg.App.Host,
			cfg.App.Port,
		),
		Handler: mux,
	}
	address := fmt.Sprintf("%s:%d", cfg.Grpc.Host, cfg.Grpc.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		zap.S().Fatalf("gRPC failed to listen %s\n", err.Error())
	}
	go func() {
		if err = src.ListenAndServe(); err != nil {
			log.Fatalf("gin listen: %s\n", err.Error())
		}

	}()
	go func() {
		if err = mux.Serve(listen); err != nil {
			log.Fatalf("gRPC failed to serve %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
}
