package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hubogle/Crontab/app/worker/config"
	"github.com/hubogle/Crontab/app/worker/initialize"
	utils "github.com/hubogle/Crontab/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitRegister() // consul 服务注册
	initialize.InitManager()  // Job Manager
}
func main() {
	var (
		port int
		err  error
	)
	router := gin.Default()
	cfg := config.GetConfig()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	port = cfg.App.Port
	if cfg.App.RunMode == gin.DebugMode {
		if port, err = utils.GetFreePort(); err != nil {
			panic(err)
		}
	}
	src := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			cfg.App.Host,
			port,
		),
		Handler: router,
	}
	grpcAddress := fmt.Sprintf("%s:%d",
		cfg.Grpc.Host,
		cfg.Grpc.Port,
	)
	go func() {
		if err = src.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err.Error())
		}
	}()
	gRpcConn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer gRpcConn.Close()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := src.Shutdown(ctx); err != nil {
		initialize.UnRegister() // 注销 consul 服务
		// TODO 当前执行的 Job 需要主动 kill 掉
		// TODO 当前执行 Job 的 Lock 可以等待超时后自动释放，也可主动删除 key
		fmt.Printf("Server Shutdown: %v\n", err)
	}
	fmt.Printf("Server exiting\n")
}
