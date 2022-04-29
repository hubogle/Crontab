package main

import (
	"fmt"
	"github.com/hubogle/Crontab/app/master/config"
	"github.com/hubogle/Crontab/app/master/initialize"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	initialize.InitLogger()
	initialize.InitConfig()
}
func main() {
	var err error
	cfg := config.GetConfig().App
	router := initialize.Router()

	src := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Handler: router,
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
