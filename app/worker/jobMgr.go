package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hubogle/Crontab/app/worker/common"
	"github.com/hubogle/Crontab/app/worker/config"
)

type JobMgr struct {
	client *api.Client
	kv     *api.KV
}

var (
	GJobMgr *JobMgr
)

// watchJobs 监听 Job 变化，监听 Job 的新增、修改、删除情况
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		// queryMeta *api.QueryMeta
		// plan     *watch.Plan
		oldKvMap map[string]*api.KVPair // 用于存储第一次，或后续更新的 Kv
		newKvMap map[string]*api.KVPair
		job      *common.Job
		jobEvent *common.JobEvent
	)
	newKvMap = make(map[string]*api.KVPair)
	wp, _ := watch.Parse(map[string]interface{}{
		"type":   "keyprefix", // 监听一个 Tree
		"prefix": config.JOB_SAVE_DIR,
	})
	wp.Handler = func(u uint64, i interface{}) {
		if oldKvMap == nil {
			oldKvMap = make(map[string]*api.KVPair)
			for _, v := range i.(api.KVPairs) {
				if job, err = common.UnpackJob(v.Value); err == nil {
					jobEvent = common.NewJobEvent(config.JOB_EVENT_SAVE, job)
					oldKvMap[v.Key] = v
					GScheduler.PushJobEvent(jobEvent) // 发送 Save 事件
				}
			}
		} else {
			newKvMap = make(map[string]*api.KVPair)
			for _, v := range i.(api.KVPairs) {
				newKvMap[v.Key] = v
				if _, ok := oldKvMap[v.Key]; !ok { // 新增 Job
					if job, err = common.UnpackJob(v.Value); err == nil {
						jobEvent = common.NewJobEvent(config.JOB_EVENT_SAVE, job)
						oldKvMap[v.Key] = v
						GScheduler.PushJobEvent(jobEvent) // 发送 Save 事件
					}
				} else if oldKvMap[v.Key].ModifyIndex != v.ModifyIndex { // 修改 Job
					if job, err = common.UnpackJob(v.Value); err == nil {
						jobEvent = common.NewJobEvent(config.JOB_EVENT_SAVE, job)
						oldKvMap[v.Key] = v
						GScheduler.PushJobEvent(jobEvent) // 发送 Save 事件
					}
				}
			}
			for k, _ := range oldKvMap {
				if _, ok := newKvMap[k]; !ok { // 删除 Job
					jobName := common.ExtractJobName(k)
					job = &common.Job{
						Name: jobName,
					}
					jobEvent = common.NewJobEvent(config.JOB_EVENT_DELETE, job)
					GScheduler.PushJobEvent(jobEvent) // 发送 Delete 事件
				}
			}
			oldKvMap = newKvMap
		}
	}
	go func() {
		err := wp.Run("127.0.0.1:8500")
		if err != nil {
			fmt.Println(err)
		}
	}()
	return
}

// watchKiller 监听 Job Kill 通知
func (jobMgr *JobMgr) watchKiller() (err error) {
	// 监听 /cron/killer 目录任务是否停止
	return
}

// NewJobLock 创建任务分布式锁
func (jobMgr *JobMgr) NewJobLock(jobName string) (jobLock *JobLock) {
	jobLock = NewJobLock(jobName, jobMgr.client)
	return
}

// InitJobMgr 初始化 Job 监听管理
func InitJobMgr() (err error) {
	var (
		config *api.Config
		client *api.Client
		kv     *api.KV
	)
	config = &api.Config{
		Address: "127.0.0.1:8500",
		Scheme:  "http",
	}

	client, err = api.NewClient(config)
	if err != nil {
		panic(err)
	}
	kv = client.KV()
	GJobMgr = &JobMgr{
		client: client,
		kv:     kv,
	}
	// 启动任务监听
	if err = GJobMgr.watchJobs(); err != nil {
		return
	}
	// 启动任务监听 killer
	if err = GJobMgr.watchKiller(); err != nil {
		return
	}
	return
}
