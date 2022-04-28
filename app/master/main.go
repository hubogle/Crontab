package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hubogle/Crontab/app/master/global"
	"github.com/hubogle/Crontab/app/master/initialize"
	"net/http"
)

func init() {
	initialize.InitConfig() // 初始化读取配置文件
}
func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	if err := r.Run(fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)); err != nil {
		panic(err)
	}
}
