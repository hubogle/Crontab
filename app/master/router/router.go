package router

import (
	"github.com/hubogle/Crontab/app/master/initialize"
	"github.com/hubogle/Crontab/app/master/pkg/core"
	"github.com/hubogle/Crontab/util/db/mysql"
	"go.uber.org/zap"
	"net/http"
)

// resource 路由需要的资源
type resource struct {
	mux    core.Mux
	db     mysql.Repo         // mysql 数据库
	logger *zap.SugaredLogger // 日志
}

func NewHTTPServer() core.Mux {
	var (
		r      *resource
		dbRepo mysql.Repo
		err    error
	)
	r = new(resource)
	dbRepo, err = initialize.InitMySQL()
	if err != nil {
		zap.S().Fatal("new db err", zap.Error(err))
	}
	mux, err := core.New()
	system := mux.Group("/")
	{
		system.GET("/ping", func(c core.Context) {
			c.JSON(http.StatusOK, "pong")
		})
	}
	r.db = dbRepo
	r.logger = zap.S()
	r.mux = mux
	setUserRouter(r)
	return mux
}
