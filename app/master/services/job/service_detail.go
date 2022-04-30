package job

import (
	"context"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/repository/dal/model"
	"github.com/hubogle/Crontab/app/master/repository/dal/query"
)

type DetailJobData struct {
	JobId int32
}

func (s *service) Detail(ctx core.Context, listData *DetailJobData) ([]*model.Job, error) {
	t := query.Use(s.db.GetDb()).Job
	do := t.WithContext(context.Background())
	result, err := do.Where(t.IsDelete.Is(false), t.ID.Eq(listData.JobId)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}
