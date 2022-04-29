package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hubogle/Crontab/app/master/config"
	"github.com/hubogle/Crontab/app/master/initialize"
	"github.com/hubogle/Crontab/app/master/router"
	"log"
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
	cfg := config.GetConfig().App
	gin.SetMode(cfg.RunMode)
	mux := router.NewHTTPServer()
	src := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Handler: mux,
	}
	go func() {
		if err = src.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
}
