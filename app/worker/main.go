package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hubogle/Crontab/app/worker/initialize"
)

func init() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitRegister()
}
func main() {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// Job 执行器
	if err := InitExecutor(); err != nil {
		panic(err)
	}
	// Job 调度器
	if err := InitScheduler(); err != nil {
		panic(err)
	}
	// Job 管理器
	if err := InitJobMgr(); err != nil {
		panic(err)
	}
	router.Run("127.0.0.1:8080")
	// TODO 任务停止后续的清理工作
}
