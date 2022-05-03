package job

import (
	"github.com/hubogle/Crontab/app/master/code"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/repository/dal/model"
	"github.com/hubogle/Crontab/app/master/services/job"
	"github.com/hubogle/Crontab/util/validation"
	"net/http"
)

type updateRequest struct {
	ID       int32  `form:"id" binding:"required"`        // 任务ID
	Name     string `form:"name"`                         // 任务名称
	Status   int32  `form:"status" binding:"min=0,max=5"` // 任务状态
	Command  string `form:"command"`                      // shell 命令
	CronExpr string `form:"cronExpr"`                     // cron 表达式
}

type updateResponse struct {
	Id int32 `json:"id"`
}

// Update 修改 Job 任务的状态 结束任务、停止任务
func (h *handler) Update() core.HandlerFunc {
	return func(c core.Context) {
		var (
			err     error
			request *updateRequest
			jobInfo model.Job
		)
		request = new(updateRequest)
		if err = c.ShouldBindJSON(request); err != nil {
			c.JSONError(core.Error(
				http.StatusBadRequest,
				code.JobUpdateError,
				validation.Error(err),
			))
			return
		}
		updateData := &job.UpdateJobData{
			Id:       request.ID,
			Name:     request.Name,
			Status:   request.Status,
			Command:  request.Command,
			CronExpr: request.CronExpr,
		}
		jobInfo, err = h.jobServer.Update(c, updateData)
		if err != nil {
			c.JSONError(core.Error(
				http.StatusInternalServerError,
				code.JobUpdateError,
				err.Error(),
			))
			return
		}
		c.JSON(http.StatusOK, jobInfo)
	}
}
