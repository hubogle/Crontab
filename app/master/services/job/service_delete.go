package job

import (
	"context"
	"github.com/hashicorp/consul/api"
	"github.com/hubogle/Crontab/app/master/config"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/repository/dal/query"
	"strconv"
)

type DeleteJobData struct {
	Id int32
}

func (s *service) Delete(ctx core.Context, deleteData *DeleteJobData) error {
	var (
		err error
		kv  *api.KV
	)
	t := query.Use(s.db.GetDb()).Job
	do := t.WithContext(context.Background())
	_, err = do.Where(t.ID.Eq(deleteData.Id)).Update(t.IsDelete, true)
	if err != nil {
		return err
	}
	kv = s.consulClient.KV()
	key := config.JOB_SAVE_DIR + strconv.Itoa(int(deleteData.Id))
	_, err = kv.Delete(key, nil)
	return err
}
