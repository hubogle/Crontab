package router

import (
	"github.com/hashicorp/consul/api"
	"github.com/hubogle/Crontab/app/master/initialize"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/util/db/mysql"
	"go.uber.org/zap"
	"net/http"
)

// resource 路由需要的资源
type resource struct {
	mux          core.Mux
	db           mysql.Repo         // mysql 数据库
	consulClient *api.Client        // consul 客户端
	logger       *zap.SugaredLogger // 日志
}

// NewHTTPServer 对路由初始化，以及一些上下文资源封装
func NewHTTPServer() core.Mux {
	var (
		r            *resource
		dbRepo       mysql.Repo
		err          error
		mux          core.Mux
		consulClient *api.Client
	)
	r = new(resource)
	dbRepo, err = initialize.InitMySQL() // 初始化 Mysql 数据库
	if err != nil {
		zap.S().Fatal("new db err", zap.Error(err))
	}
	consulClient, err = initialize.InitConsul() // 初始化 Consul 数据库
	if err != nil {
		zap.S().Fatal("new consul err", zap.Error(err))
	}
	mux, err = core.New()
	system := mux.Group("/")
	{
		system.GET("/ping", func(c core.Context) {
			c.JSON(http.StatusOK, "pong")
		})
	}
	r.db = dbRepo
	r.consulClient = consulClient
	r.logger = zap.S()
	r.mux = mux
	setUserRouter(r)
	return mux
}
