package router

import (
	"github.com/hubogle/Crontab/app/master/api/job"
)

func setUserRouter(r *resource) {
	jobGroup := r.mux.Group("job")
	{
		jobHandler := job.New(r.db, r.consulClient)
		jobGroup.GET("/detail/:id", jobHandler.Detail())
		jobGroup.GET("/list", jobHandler.List())
		jobGroup.POST("/create", jobHandler.Create())
		jobGroup.POST("/delete", jobHandler.Delete())
		jobGroup.POST("/update", jobHandler.Update())
	}
}
