package job

import (
	"github.com/hubogle/Crontab/app/master/code"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/services/job"
	"github.com/hubogle/Crontab/util/validation"
	"net/http"
	"time"
)

type detailRequest struct {
	ID int32 `uri:"id"`
}

type detailResponse struct {
	Name       string `json:"name"`
	Status     int32  `json:"status"`
	PlanTime   string `json:"planTime"`
	NextTime   string `json:"nextTime"`
	CreateTime string `json:"createTIme"`
}

func (h *handler) Detail() core.HandlerFunc {
	return func(c core.Context) {
		var (
			err      error
			request  *detailRequest
			response *detailResponse
		)
		request = new(detailRequest)
		response = new(detailResponse)
		if err = c.ShouldBindURI(&request); err != nil {
			c.JSONError(core.Error(
				http.StatusBadRequest,
				code.JobDetailError,
				validation.Error(err)))
			return
		}
		detailData := &job.DetailJobData{JobId: request.ID}
		result, err := h.jobServer.Detail(c, detailData)
		if err != nil || len(result) == 0 {
			c.JSONError(core.Error(
				http.StatusBadRequest,
				code.JobDetailError,
				validation.Error(err)))
			return
		}
		response.Name = result[0].Name
		response.Status = result[0].Status
		response.PlanTime = result[0].PlanTime.Format(time.RFC3339)
		response.NextTime = result[0].NextTime.Format(time.RFC3339)
		response.CreateTime = result[0].Created.Format(time.RFC3339)
		c.JSON(http.StatusOK, response)
	}
}
