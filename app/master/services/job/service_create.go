package job

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"github.com/hubogle/Crontab/app/master/common"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/repository/dal/model"
	"github.com/hubogle/Crontab/app/master/repository/dal/query"
	"github.com/hubogle/Crontab/app/worker/config"
	"strconv"
)

type CreateJobData struct {
	Name     string // 任务名称
	Status   int    // 任务状态
	Command  string // shell 命令
	CronExpr string // cron 表达式
}

func (s *service) Create(ctx core.Context, jobData *CreateJobData) (id int32, err error) {
	var (
		value []byte
	)
	q := query.Use(s.db.GetDb())
	kv := s.consulClient.KV()
	t := q.Job
	do := t.WithContext(context.Background())
	job := &model.Job{
		Name:     jobData.Name,
		Status:   0,
		Command:  jobData.Command,
		CronExpr: jobData.CronExpr,
	}
	err = do.Create(job) // 创建数据库对象，任务对象
	jobApi := common.Job{
		Name:     jobData.Name,
		Command:  jobData.Command,
		CronExpr: jobData.CronExpr,
	}
	key := config.JOB_SAVE_DIR + strconv.Itoa(int(job.ID)) // 根据任务创好的 ID 作为 key
	if value, err = json.Marshal(jobApi); err == nil {
		_, err = kv.Put(&api.KVPair{Key: key, Value: value}, nil)
		if err != nil {
			return 0, err
		}
	}
	if err != nil {
		return 0, err
	}
	return job.ID, nil
}
