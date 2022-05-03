package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"github.com/hubogle/Crontab/app/worker/config"
	"strconv"
	"strings"
	"time"
)

// UnpackJob 从存储的 value 值反序列化 job
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)
	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}

// ExtractJobId 从 etcd 的 key 中提取任务名称
func ExtractJobId(jobKey string) int {
	id, _ := strconv.Atoi(strings.TrimPrefix(jobKey, config.JOB_SAVE_DIR))
	return id
}

// ExtractKillerJobId 监听杀死任务名称
func ExtractKillerJobId(killerKey string) int {
	id, _ := strconv.Atoi(strings.TrimPrefix(killerKey, config.JOB_KILLER_DIR))
	return id
}

// BuildJobExecuteInfo 构造任务执行信息
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime, // 计划执行时间
		RealTime: time.Now(),               // 真实执行时间
		NextTime: jobSchedulePlan.NextTime, // 下次执行时间
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}

// BuildJobSchedulePlan 解析 Cron 表达式，生成调度执行对象
func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil { // 解析 JOB 的 Cron 表达式
		return
	}
	jobSchedulePlan = &JobSchedulePlan{ // 任务调度计划对象
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}
