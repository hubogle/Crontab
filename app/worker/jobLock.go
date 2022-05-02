package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/hubogle/Crontab/app/worker/config"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

// 基于 Consul Session 实现分布式锁

// JobLock 锁
type JobLock struct {
	client    *api.Client // consul 客户端
	kv        *api.KV
	session   *api.Session  // 创建 Session 对象
	jobName   string        // 任务名称
	isLocked  bool          // 是否上锁成功
	sessionId string        // 生成 Session 的唯一标识
	lockKey   string        // 锁的key
	stopCh    chan struct{} // 停止 Session 续期
}

func NewJobLock(jobName string, client *api.Client) *JobLock {
	return &JobLock{
		client:  client,
		jobName: jobName,
		kv:      client.KV(),
		session: client.Session(),
	}
}
func (j *JobLock) TryLock() (err error) {
	var (
		locKey     string
		stopCh     chan struct{}
		createId   string // 创建 Session 后的唯一标识
		isAcquired bool   // 是否抢锁成功
	)
	locKey = config.JOB_LOCK_DIR + j.jobName
	createId, _, err = j.client.Session().Create(&api.SessionEntry{
		Name:      uuid.NewV4().String(),
		Behavior:  "delete",
		TTL:       "20s",           // Session 过期时间
		LockDelay: 2 * time.Second, // 防止 KV 被释放后立马被获取到
	}, nil)
	isAcquired, _, err = j.kv.Acquire(&api.KVPair{
		Key:     locKey,
		Session: createId,
	}, nil)
	if !isAcquired || err != nil { // 未抢占到锁
		err = errors.New("lock already required")
		return
	} else {
		stopCh = make(chan struct{})
		go j.session.RenewPeriodic("20s", createId, nil, stopCh) // 定时给 Session 延长 TTL
		j.lockKey = locKey
		j.isLocked = true
		j.sessionId = createId
		j.stopCh = stopCh
		return nil
	}
}

// Unlock 如果抢到锁的话释放锁
func (j *JobLock) Unlock() {
	if j.isLocked == true {
		j.kv.Release(&api.KVPair{
			Key:     j.lockKey,
			Session: j.sessionId,
		}, nil) // 释放绑定 key 的锁
		if j.stopCh != nil {
			close(j.stopCh)
		}
		// j.kv.Delete(j.lockKey, nil) // 删除 Key
	}
}
