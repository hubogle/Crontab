package common

import (
	"context"
	"github.com/gorhill/cronexpr"
	"time"
)

// Job 任务定义
type Job struct {
	Name     string `json:"name"`     // 任务名称
	Command  string `json:"command"`  // 任务命令
	CronExpr string `json:"cronExpr"` // cron表达式
}

// JobEvent 任务的变化事件有两种，更新 or 删除
type JobEvent struct {
	EventType int  `json:"eventType"` // 事件类型
	Job       *Job `json:"job"`       // 任务
}

// JobSchedulePlan 任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 // 任务信息
	Expr     *cronexpr.Expression // cron 表达式
	NextTime time.Time            // 下次调度时间
}

// JobExecuteInfo 任务执行状态
type JobExecuteInfo struct {
	Job        *Job               // 任务信息
	PlanTime   time.Time          // 理论上的调度时间
	RealTime   time.Time          // 实际的调度时间
	CancelCtx  context.Context    // 用于取消任务的上下文
	CancelFunc context.CancelFunc // 用于取消任务的函数
}

// JobExecuteResult 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo // 执行状态
	Output      []byte          // 脚本输出
	Err         error           // 脚本错误原因
	StartTime   time.Time       // 启动时间
	EndTime     time.Time       // 结束时间
}

// NewJobEvent 构建任务执行变化事件
func NewJobEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}
