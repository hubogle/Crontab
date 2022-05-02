package job

import (
	"github.com/hashicorp/consul/api"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/app/master/services/job"
	"github.com/hubogle/Crontab/util/db/mysql"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()
	Create() core.HandlerFunc
	Delete() core.HandlerFunc
	Detail() core.HandlerFunc
	List() core.HandlerFunc
}

func (h *handler) i() {}

// handler 抽离每个接口需要的资源
type handler struct {
	repo         mysql.Repo
	consulClient *api.Client
	jobServer    job.Service
}

func New(db mysql.Repo, client *api.Client) Handler {
	return &handler{
		repo:         db,
		consulClient: client,
		jobServer:    job.New(db, client),
	}
}
