package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/hubogle/Crontab/app/master/router"
	"net/http"
)

func Router() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	apiGroup := Router.Group("/v1")
	router.InitUserRouter(apiGroup)
	return Router
}
