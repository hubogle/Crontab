package job

import (
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/services/job"
	"github.com/hubogle/Crontab/util/db/mysql"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()
	Create() core.HandlerFunc
}

func (h *handler) i() {}

// handler 抽离每个接口需要的资源
type handler struct {
	repo      mysql.Repo
	jobServer job.Service
}

func New(db mysql.Repo) Handler {
	return &handler{repo: db, jobServer: job.New(db)}
}
