package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hubogle/Crontab/app/master/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	jobGroup := router.Group("job")
	{
		jobGroup.GET("/job/save", api.HandleJobCreate)
		jobGroup.GET("/job/list", api.HandleJobList)
		jobGroup.GET("/job/delete", api.HandleJobDelete)
	}
}
