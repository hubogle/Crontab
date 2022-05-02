package job

import (
	"github.com/hashicorp/consul/api"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/repository/dal/model"
	"github.com/hubogle/Crontab/util/db/mysql"
)

var _ Service = (*service)(nil)

type Service interface {
	i()
	Create(ctx core.Context, jobData *CreateJobData) (id int32, err error)
	Delete(ctx core.Context, deleteData *DeleteJobData) error
	List(ctx core.Context, listData *ListJobData) ([]*model.Job, int, error)
	Detail(ctx core.Context, listData *DetailJobData) ([]*model.Job, error)
}

type service struct {
	db           mysql.Repo
	consulClient *api.Client
}

func (s *service) i() {}

func New(db mysql.Repo, client *api.Client) Service {
	return &service{
		db:           db,
		consulClient: client,
	}
}
