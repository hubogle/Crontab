package job

import (
	"github.com/hubogle/Crontab/app/master/code"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/services/job"
	"github.com/hubogle/Crontab/util/validation"
	"net/http"
	"time"
)

type listRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit" binding:"required"`
}
type jobData struct {
	Name       string `json:"name"`
	Status     int32  `json:"status"`
	PlanTime   string `json:"planTime"`
	NextTime   string `json:"nextTime"`
	CreateTime string `json:"createTIme"`
}

type listResponse struct {
	Count  int       `json:"count"`
	Offset int       `json:"offset"`
	Limit  int       `json:"limit"`
	Data   []jobData `json:"data"`
}

func (h *handler) List() core.HandlerFunc {
	return func(c core.Context) {
		var (
			err      error
			request  *listRequest
			response *listResponse
		)
		request = new(listRequest)
		response = new(listResponse)
		if err = c.ShouldBindQuery(request); err != nil {
			c.JSONError(core.Error(
				http.StatusBadRequest,
				code.JobListError,
				validation.Error(err)))
			return
		}
		listData := &job.ListJobData{
			Offset: request.Offset,
			Limit:  request.Limit,
		}
		result, count, err := h.jobServer.List(c, listData)
		if err != nil {
			c.JSONError(core.Error(
				http.StatusBadRequest,
				code.JobListError,
				validation.Error(err)))
			return
		}
		response.Data = make([]jobData, 0, len(result))
		response.Offset = request.Offset
		response.Limit = request.Limit
		response.Count = count
		for _, v := range result {
			data := jobData{
				Name:       v.Name,
				Status:     v.Status,
				PlanTime:   v.PlanTime.Format(time.RFC3339),
				NextTime:   v.NextTime.Format(time.RFC3339),
				CreateTime: v.Created.Format(time.RFC3339),
			}
			response.Data = append(response.Data, data)
		}
		c.JSON(http.StatusOK, response)
	}
}
