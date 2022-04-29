package initialize

import (
	"github.com/hubogle/Crontab/app/master/config"
	"github.com/hubogle/Crontab/util/validation"
)

func InitTrans() error {
	cfg := config.GetConfig().App
	return validation.LocalTrans(cfg.Language)
}
