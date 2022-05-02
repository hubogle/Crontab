package main

import (
	"fmt"
	"github.com/hubogle/Crontab/app/worker/common"
	"math/rand"
	"os/exec"
	"time"
)

// Executor 任务执行器
type Executor struct{}

var (
	GExecutor *Executor
)

// ExecuteJob 执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			result *common.JobExecuteResult
			cmd    *exec.Cmd
			output []byte
		)
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		jobLock := GJobMgr.NewJobLock(info.Job.Name) // 分布式锁
		result.StartTime = time.Now()                // 任务开始时间
		// 随机睡眠 0~1s ，确保多个节点，能够同时执行任务竞争均匀
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		err := jobLock.TryLock()
		defer jobLock.Unlock()
		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			// TODO 临时先执行 SHELL 命令
			result.StartTime = time.Now() // 任务真正执行时间
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)
			output, err = cmd.CombinedOutput()
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}
		// 返回调度器给 Scheduler，然后从调度器中删除执行完成的任务
		fmt.Println(result.StartTime, string(result.Output))
		GScheduler.PushJobResult(result)
	}()
}

// InitExecutor 初始化任务执行器
func InitExecutor() (err error) {
	GExecutor = &Executor{}
	return
}
