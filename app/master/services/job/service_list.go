package job

import (
	"context"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/repository/dal/model"
	"github.com/hubogle/Crontab/app/master/repository/dal/query"
)

type ListJobData struct {
	Offset int
	Limit  int
}

func (s *service) List(ctx core.Context, listData *ListJobData) ([]*model.Job, int, error) {
	t := query.Use(s.db.GetDb()).Job
	do := t.WithContext(context.Background())
	result, count, err := do.Where(t.IsDelete.Is(false)).FindByPage(listData.Offset, listData.Limit)
	if err != nil {
		return nil, 0, err
	}
	return result, int(count), nil
}
