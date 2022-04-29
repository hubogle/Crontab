package job

import (
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/util/db/mysql"
)

var _ Service = (*service)(nil)

type Service interface {
	i()
	Create(ctx core.Context, jobData *CreateJobData) (id int32, err error)
}

type service struct {
	db mysql.Repo
}

func (s *service) i() {}

func New(db mysql.Repo) Service {
	return &service{
		db: db,
	}
}
