package job

import "github.com/hubogle/Crontab/util/db/mysql"

var _ Service = (*service)(nil)

type Service interface {
	i()
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
