package job

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"github.com/hubogle/Crontab/app/master/common"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/repository/dal/model"
	"github.com/hubogle/Crontab/app/master/repository/dal/query"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

type UpdateJobData struct {
	Id       int32  // 任务 id
	Name     string // 任务名称
	Status   int32  // 任务状态
	Command  string // shell 命令
	CronExpr string // cron 表达式
}

func (s *service) Update(ctx core.Context, jobData *UpdateJobData) (model.Job, error) {
	q := query.Use(s.db.GetDb())
	kv := s.consulClient.KV()
	t := q.Job
	do := t.WithContext(context.Background())
	params := map[string]interface{}{
		"name":     jobData.Name,
		"status":   jobData.Status,
		"command":  jobData.Command,
		"cronExpr": jobData.CronExpr,
	}
	_, err := do.Where(t.ID.Eq(jobData.Id)).Updates(params)
	if err != nil {
		return model.Job{}, err
	}
	job, _ := do.Where(t.ID.Eq(jobData.Id)).First()
	stopKey := common.JOB_SAVE_DIR + strconv.Itoa(int(jobData.Id))
	killKey := common.JOB_KILLER_DIR + strconv.Itoa(int(jobData.Id))
	switch job.Status {
	case common.Kill: // 停止这一次运行的任务
		kv.Put(&api.KVPair{Key: killKey, Value: uuid.NewV4().Bytes()}, nil)
	case common.Canceled: // 取消掉以后任务
		kv.Delete(stopKey, nil)
	// case common.Running: // 该状态由程序运行
	// case common.Success: // 任务正常完成
	// case common.Pending: // 任务等待中
	default:
		key := common.JOB_SAVE_DIR + strconv.Itoa(int(job.ID)) // 根据任务创好的 ID 作为 key
		jobApi := common.Job{
			Id:       int(job.ID),
			Name:     job.Name,
			Command:  job.Command,
			CronExpr: job.CronExpr,
		}
		if value, err := json.Marshal(jobApi); err == nil {
			_, err = kv.Put(&api.KVPair{Key: key, Value: value}, nil)
			if err != nil {
				return *job, err
			}
		}
	}
	return *job, err
}
