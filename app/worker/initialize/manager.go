package initialize

import "github.com/hubogle/Crontab/app/worker/manager"

func InitManager() {
	// Job 执行器
	if err := manager.InitExecutor(); err != nil {
		panic(err)
	}
	// Job 调度器
	if err := manager.InitScheduler(); err != nil {
		panic(err)
	}
	// Job 管理器
	if err := manager.InitJobMgr(); err != nil {
		panic(err)
	}
}
