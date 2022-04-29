package router

import (
	"github.com/hubogle/Crontab/app/master/api/job"
)

func setUserRouter(r *resource) {
	jobGroup := r.mux.Group("job")
	{
		jobHandler := job.New(r.db)
		jobGroup.POST("/create", jobHandler.Create())
	}
}
