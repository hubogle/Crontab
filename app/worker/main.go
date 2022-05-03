package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hubogle/Crontab/app/worker/config"
	"github.com/hubogle/Crontab/app/worker/initialize"
	"github.com/hubogle/Crontab/app/worker/manager"
	utils "github.com/hubogle/Crontab/util"
)

func init() {
	initialize.InitLogger()
	initialize.InitConfig()
	// initialize.InitRegister()
}
func main() {
	var (
		port int
		err  error
	)
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// Job 执行器
	if err := manager.InitExecutor(); err != nil {
		panic(err)
	}
	// Job 调度器
	if err := manager.InitScheduler(); err != nil {
		panic(err)
	}
	// Job 管理器
	if err := manager.InitJobMgr(); err != nil {
		panic(err)
	}
	if config.GetConfig().App.RunMode == gin.DebugMode {
		if port, err = utils.GetFreePort(); err != nil {
			panic(err)
		}
	} else {
		port = config.GetConfig().App.Port
	}
	addres := fmt.Sprintf("%s:%d", config.GetConfig().App.Host, port)
	router.Run(addres)
	// TODO 任务停止后续的清理工作
	// TODO 需要停止所有对 key 的上锁操作
}
