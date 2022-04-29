package initialize

import (
	"github.com/hubogle/Crontab/app/master/config"
	"github.com/hubogle/Crontab/util/db/mysql"
)

func InitMySQL() (mysql.Repo, error) {
	cfg := config.GetConfig().MySQL
	db := mysql.Mysql{
		Host:            cfg.Host,
		Port:            cfg.Port,
		User:            cfg.User,
		Pass:            cfg.Pass,
		DbName:          cfg.DBName,
		MaxIdleConn:     cfg.MaxIdleConn,
		MaxOpenConn:     cfg.MaxOpenConn,
		MaxLifetimeConn: cfg.MaxLifeTimeConn,
	}
	return db.Client()
}
