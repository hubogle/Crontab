package executor

import (
	"fmt"
	"github.com/hubogle/Crontab/app/worker/common"
)

// Executor 任务执行器
type Executor struct {
	// Worker *Worker
}

var (
	GExecutor *Executor
)

// ExecuteJob 执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		fmt.Println(info.Job.Name)
	}()
}

// InitExecutor 初始化执行器
func InitExecutor() (err error) {
	GExecutor = &Executor{}
	return
}
