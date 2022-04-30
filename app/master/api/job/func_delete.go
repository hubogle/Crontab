package job

import (
	"github.com/hubogle/Crontab/app/master/code"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/services/job"
	"github.com/hubogle/Crontab/util/validation"
	"net/http"
)

type deleteRequest struct {
	Id int32 `form:"id" binding:"required"` // 任务ID
}

type deleteResponse struct {
	Delete bool `json:"delete"` // 任务ID
}

func (h *handler) Delete() core.HandlerFunc {
	return func(c core.Context) {
		var (
			err      error
			request  *deleteRequest
			response *deleteResponse
		)
		request = new(deleteRequest)
		response = new(deleteResponse)
		if err = c.ShouldBindJSON(request); err != nil {
			c.JSONError(core.Error(
				http.StatusBadRequest,
				code.JobDeleteError,
				validation.Error(err),
			))
			return
		}
		deleteData := &job.DeleteJobData{Id: request.Id}
		err = h.jobServer.Delete(c, deleteData)
		if err == nil {
			response.Delete = true
		}
		c.JSON(http.StatusOK, response)
	}
}
