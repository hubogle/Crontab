package manager

import (
	"fmt"
	"github.com/hubogle/Crontab/app/worker/common"
	"github.com/hubogle/Crontab/app/worker/config"
	"time"
)

// Scheduler 任务调度协程
type Scheduler struct {
	jobEventChan      chan *common.JobEvent           // consul 任务执行事件队列
	jobPlanTable      map[int]*common.JobSchedulePlan // 任务调度计划表
	jobExecutingTable map[int]*common.JobExecuteInfo  // 任务执行表，每个节点会抢到任务执行，通过返回的信息的 error 判断
	jobResultChan     chan *common.JobExecuteResult   // 任务结果队列
}

var (
	GScheduler *Scheduler
)

// scheduleLoop 任务调度协程，处理 Job 任务的所有调度
func (scheduler *Scheduler) scheduleLoop() {
	// 定时任务 common_jobEvent
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult     *common.JobExecuteResult
	)
	scheduleAfter = scheduler.TrySchedule() // 初始化一次

	scheduleTimer = time.NewTimer(scheduleAfter) // 下次调度时间唤醒执行

	for {
		select {
		case jobEvent = <-scheduler.jobEventChan: // 监听任务变化事件
			// 对内存中维护的任务列表增删改
			scheduler.handleJobEvent(jobEvent)
		case jobResult = <-scheduler.jobResultChan: // 监听任务执行结果
			scheduler.handleJobResult(jobResult)
		case <-scheduleTimer.C: // 最近的任务到期了
		}
		scheduleAfter = scheduler.TrySchedule() // 计算任务调度间隔
		scheduleTimer.Reset(scheduleAfter)      // 重置调度间隔
	}
}

// TrySchedule 计算任务调度时间，返回下一次调度时间
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = time.Second * 1
		return
	}
	now = time.Now()
	// 遍历所有任务
	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) { // 有任务超过执行时间
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) // 更新下次执行时间
		}
		// 统计下次任务执行时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	// 下次调度时间间隔 （最新执行时间-当前时间）
	scheduleAfter = (*nearTime).Sub(now)
	return
}

// TryStartJob 尝试执行任务
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)
	// 执行任务运行很久
	if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Id]; jobExecuting {
		fmt.Println("跳过执行，任务正在执行", jobPlan.Job.Name)
		return
	}
	// 构建执行状态信息
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)
	// 保存执行状态
	scheduler.jobExecutingTable[jobPlan.Job.Id] = jobExecuteInfo
	// 执行任务
	// fmt.Println("执行任务：", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	GExecutor.ExecuteJob(jobExecuteInfo) // 调度器真实执行任务
}

// handleJobEvent 处理任务事件，保存、删除、停止任务
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExecuteInfo  *common.JobExecuteInfo
		jobExisted      bool
		err             error
	)
	switch jobEvent.EventType {
	case config.JOB_EVENT_SAVE: // 解析执行任务的 Cron 表达式
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Id] = jobSchedulePlan
	case config.JOB_EVENT_DELETE: // 删除任务事件
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Id]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Id)
		}
	case config.JOB_EVENT_KILLER: // 停止任务事件
		// 取消掉Command执行, 判断任务是否在执行中
		if jobExecuteInfo, jobExisted = scheduler.jobExecutingTable[jobEvent.Job.Id]; jobExisted {
			jobExecuteInfo.CancelFunc() // 触发command杀死shell子进程, 任务得到退出
		}
		_ = GJobMgr.DeleteKey(jobEvent.Job.Id, config.JOB_EVENT_KILLER) // 删除掉 key 值
	}
}

// handleJobResult 处理任务获取到的结果
func (scheduler Scheduler) handleJobResult(result *common.JobExecuteResult) {
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Id)
}

// PushJobEvent 将 Job 任务 Event 放入到队列进行处理
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

// PushJobResult 回传任务执行结果
func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}

// InitScheduler 初始化 Job 任务调度器
func InitScheduler() (err error) {
	GScheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[int]*common.JobSchedulePlan, 1000),
		jobExecutingTable: make(map[int]*common.JobExecuteInfo, 1000),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}
	go GScheduler.scheduleLoop() // 启动协程调度
	return
}
