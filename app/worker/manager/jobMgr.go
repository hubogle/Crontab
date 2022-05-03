package manager

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hubogle/Crontab/app/worker/common"
	"github.com/hubogle/Crontab/app/worker/config"
	"github.com/hubogle/Crontab/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type JobMgr struct {
	client    *api.Client
	kv        *api.KV
	rpcClient proto.JobClient
	timeOut   time.Duration
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
					jobId := common.ExtractJobId(k)
					job = &common.Job{
						Id: jobId,
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
	var (
		jobEvent *common.JobEvent
		jobId    int
		job      *common.Job
	)
	// 监听/cron/killer目录
	go func() {
		wp, _ := watch.Parse(map[string]interface{}{
			"type":   "keyprefix", // 监听一个 Tree
			"prefix": config.JOB_KILLER_DIR,
		})
		// newKvMap := make(map[string]*api.KVPair)
		wp.Handler = func(u uint64, i interface{}) {
			for _, v := range i.(api.KVPairs) {
				jobId = common.ExtractKillerJobId(v.Key)
				job = &common.Job{
					Id: jobId,
				}
				jobEvent = common.NewJobEvent(config.JOB_EVENT_KILLER, job)
				GScheduler.PushJobEvent(jobEvent) // 发送 Kill 事件
			}
		}
		err := wp.Run("127.0.0.1:8500")
		if err != nil {
			fmt.Println(err)
		}
	}()
	return
}

// NewJobLock 创建任务分布式锁
func (jobMgr *JobMgr) NewJobLock(jobId int) (jobLock *JobLock) {
	jobLock = NewJobLock(jobId, jobMgr.client)
	return
}

// UpdateJob 调用 gRPC 更改数据库状态
func (jobMgr *JobMgr) UpdateJob(jobId int, status int32, planTime int64, nextTime int64) {
	ctx, cancel := context.WithTimeout(context.Background(), jobMgr.timeOut*time.Second)
	defer cancel()
	_, err := jobMgr.rpcClient.UpdateJob(ctx, &proto.UpdateJobInfo{
		JobId:    uint32(jobId),
		Status:   uint32(status),
		PlanTime: uint64(planTime),
		NextTime: uint64(nextTime),
	})
	if err != nil {
		zap.S().Error("RPC 修改失败:", err.Error())
	}
}

// DeleteKey 删除 Key
func (jobMgr *JobMgr) DeleteKey(jobId int, event int) error {
	var (
		err error
		key string
	)
	switch event {
	case config.JOB_EVENT_SAVE:
		key = fmt.Sprintf("%s%d", config.JOB_SAVE_DIR, jobId)
	case config.JOB_EVENT_DELETE:
		key = fmt.Sprintf("%s%d", config.JOB_SAVE_DIR, jobId)
	case config.JOB_EVENT_KILLER:
		key = fmt.Sprintf("%s%d", config.JOB_KILLER_DIR, jobId)
	case config.JOB_EVENT_LOCK:
		key = fmt.Sprintf("%s%d", config.JOB_KILLER_DIR, jobId)
	default:
	}
	_, err = jobMgr.kv.Delete(key, nil)
	return err
}

// InitJobMgr 初始化 Job 监听管理
func InitJobMgr() (err error) {
	var (
		client    *api.Client
		kv        *api.KV
		conn      *grpc.ClientConn
		address   string
		rpcClient proto.JobClient
	)
	cfg := config.GetConfig()
	address = fmt.Sprintf("%s:%d", cfg.Consul.Host, cfg.Consul.Port)
	apiConfig := &api.Config{
		Address: address,
		Scheme:  "http",
	}
	client, err = api.NewClient(apiConfig) // 链接 Consul
	if err != nil {
		panic(err)
	}
	address = fmt.Sprintf("%s:%d", cfg.Grpc.Host, cfg.Grpc.Port)
	conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	rpcClient = proto.NewJobClient(conn)
	if err != nil {
		zap.S().Panic("RPC 链接失败:", err.Error())
	}
	kv = client.KV()
	GJobMgr = &JobMgr{
		client:    client,
		kv:        kv,
		rpcClient: rpcClient,
		timeOut:   cfg.Grpc.MaxTimeConn,
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
