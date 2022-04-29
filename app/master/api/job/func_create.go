package job

import (
	"github.com/hubogle/Crontab/app/master/code"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/services/job"
	"github.com/hubogle/Crontab/util/validation"
	"net/http"
)

type createRequest struct {
	Name     string `form:"name" binding:"required"`     // 任务名称
	Status   int    `form:"status" binding:"required"`   // 任务状态
	Command  string `form:"command" binding:"required"`  // shell 命令
	CronExpr string `form:"cronExpr" binding:"required"` // cron 表达式
}

type createResponse struct {
	Id int32 `json:"id"`
}

func (h *handler) Create() core.HandlerFunc {
	return func(c core.Context) {
		var (
			err      error
			request  *createRequest
			response *createResponse
			id       int32
		)
		request = new(createRequest)
		response = new(createResponse)
		if err = c.ShouldBindJSON(request); err != nil {
			c.JSONError(core.Error(
				http.StatusBadRequest,
				code.JobCreateError,
				validation.Error(err),
			))
			return
		}
		createData := &job.CreateJobData{
			Name:     request.Name,
			Status:   request.Status,
			Command:  request.Command,
			CronExpr: request.CronExpr,
		}
		id, err = h.jobServer.Create(c, createData)
		if err != nil {
			c.JSONError(core.Error(
				http.StatusInternalServerError,
				code.JobCreateError,
				err.Error(),
			))
			return
		}
		response.Id = id
		c.JSON(http.StatusOK, response)
	}
}
