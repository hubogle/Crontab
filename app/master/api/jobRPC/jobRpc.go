package jobRPC

import (
	"context"
	"github.com/hubogle/Crontab/app/master/repository/dal/query"
	"github.com/hubogle/Crontab/proto"
	"github.com/hubogle/Crontab/util/db/mysql"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Job struct {
	proto.UnimplementedJobServer
	db mysql.Repo
}

func NewJob(db mysql.Repo) *Job {
	return &Job{
		db: db,
	}
}

// UpdateJob 通过 gRPC 调用服务更新状态
func (j *Job) UpdateJob(ctx context.Context, jobInfo *proto.UpdateJobInfo) (*emptypb.Empty, error) {
	q := query.Use(j.db.GetDb())
	t := q.Job
	do := t.WithContext(ctx)
	params := map[string]interface{}{
		"status": jobInfo.Status,
	}
	if jobInfo.NextTime != 0 {
		params["nextTime"] = time.Unix(int64(jobInfo.NextTime), 0)
	}
	if jobInfo.PlanTime != 0 {
		params["planTIme"] = time.Unix(int64(jobInfo.PlanTime), 0)
	}
	_, err := do.Where(t.ID.Eq(int32(jobInfo.JobId))).Updates(params)
	return &emptypb.Empty{}, err
}
