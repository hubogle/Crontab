package job

import (
	"context"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/repository/dal/model"
	"github.com/hubogle/Crontab/app/master/repository/dal/query"
)

type CreateJobData struct {
	Name     string // 任务名称
	Status   int    // 任务状态
	Command  string // shell 命令
	CronExpr string // cron 表达式
}

func (s *service) Create(ctx core.Context, jobData *CreateJobData) (id int32, err error) {
	q := query.Use(s.db.GetDb())
	t := q.Job
	do := t.WithContext(context.Background())
	job := &model.Job{
		Name:     jobData.Name,
		Status:   1,
		Command:  jobData.Command,
		CronExpr: jobData.CronExpr,
	}
	err = do.Create(job)
	if err != nil {
		return 0, err
	}
	return job.ID, nil
}
