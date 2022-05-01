package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hubogle/Crontab/app/worker/executor"
	"github.com/hubogle/Crontab/app/worker/initialize"
	"github.com/hubogle/Crontab/app/worker/scheduler"
)

func init() {
	initialize.InitLogger()
	initialize.InitConfig()
	// initialize.InitRegister()
}
func main() {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// Job 执行器
	if err := executor.InitExecutor(); err != nil {
		panic(err)
	}
	// Job 调度器
	if err := scheduler.InitScheduler(); err != nil {
		panic(err)
	}
	if err := InitJobMgr(); err != nil {
		panic(err)
	}
	router.Run("127.0.0.1:8080")
}
