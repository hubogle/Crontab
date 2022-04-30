package job

import (
	"context"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/repository/dal/query"
)

type DeleteJobData struct {
	Id int32
}

func (s *service) Delete(ctx core.Context, deleteData *DeleteJobData) error {
	var err error
	t := query.Use(s.db.GetDb()).Job
	do := t.WithContext(context.Background())
	_, err = do.Where(t.ID.Eq(deleteData.Id)).Update(t.IsDelete, true)
	if err != nil {
		return err
	}
	return err
}
