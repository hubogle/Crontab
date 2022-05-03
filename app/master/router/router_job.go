package router

import (
	"github.com/hubogle/Crontab/app/master/api/job"
	"github.com/hubogle/Crontab/app/master/api/jobRPC"
	"github.com/hubogle/Crontab/proto"
)

// setJobRouter 设置 Gin 路由，以及接口需要的资源
func setJobRouter(r *resource) {
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

// setJobGRPC 设置 gRPC 资源
func setJobGRPC(r *resource) {
	jobGrpc := jobRPC.NewJob(r.db)
	r.mux.RegisterServer(jobGrpc, &proto.Job_ServiceDesc)
}
