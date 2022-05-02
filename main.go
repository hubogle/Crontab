package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}
	config := &api.Config{
		Address: "127.0.0.1:8500",
		Scheme:  "http",
	}

	client, _ := api.NewClient(config)
	// var ch <-chan struct{}
	ch := make(chan struct{})
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go tryLock(client, "mylock1", "session"+strconv.Itoa(i), wg, ch)
	}
	wg.Wait()

}

func tryLock(client *api.Client, key string, sessionName string, wg *sync.WaitGroup, ch chan struct{}) {
	defer wg.Done()
	createId, _, err := client.Session().Create(&api.SessionEntry{
		Name:     sessionName,
		TTL:      "10s",
		Behavior: "delete",
	}, nil)
	// fmt.Println(createId, w.RequestTime, err)

	// opts := &api.LockOptions{
	// 	Key: key,
	// 	// Value:      []byte("sender 1"),
	// 	// SessionTTL: "10s",
	// 	Session: createId,
	// 	// SessionOpts: &api.SessionEntry{
	// 	// 	Checks:   []string{"check1", "check2"},
	// 	// 	Behavior: "release",
	// 	// },
	// 	// SessionName: uuid.NewV4().String(),
	// }
	// lock, _ := client.LockOpts(opts)
	isAcquired, _, err := client.KV().Acquire(&api.KVPair{
		Key:     key,
		Value:   []byte("any value"),
		Session: createId,
	}, nil)

	if err != nil {
		fmt.Println(err)
		// 错误处理
	}
	if isAcquired {
		fmt.Println("获取成功", sessionName)
		fmt.Println(isAcquired)
		// go client.Session().Renew(createId, nil)
		ch1 := make(chan struct{})
		go client.Session().RenewPeriodic("20s", createId, nil, ch1)
		time.Sleep(5 * time.Second)
		// client.Session().Destroy(createId, nil)
		kv, _, _ := client.KV().Get(key, nil)
		kv1, _, _ := client.KV().Get("go-consul-test", nil)
		fmt.Println(kv.LockIndex, kv.Session)
		fmt.Println(kv1.LockIndex, kv1.Session)
		client.KV().Release(&api.KVPair{
			Key:     key,
			Value:   []byte("any value"),
			Session: createId,
		}, nil)
		close(ch1)
		fmt.Println("释放完")
	} else {
		fmt.Println("获取失败", sessionName)
	}
	// log.Println(sessionName, "try to get lock obj")
	// leaderCh, err := lock.Lock(ch)
	// if err != nil || leaderCh == nil {
	// 	log.Println("err", err, sessionName)
	// 	return
	// }
	// time.Sleep(time.Second * 5)
	// ch1 := make(chan struct{})
	// for i := 0; i < 10; i++ {
	// 	time.Sleep(time.Second * 1)
	// 	fmt.Println(createId)
	// client.Session().Renew(createId, nil)
	// go client.Session().RenewPeriodic("5s", createId, nil, ch1)
	// time.Sleep(time.Second * 5)
	// ch1 <- struct{}{}

	// }
	// close(ch)
	// log.Println(sessionName, "lock and sleep")
	// time.Sleep(3 * time.Second)
	// ch <- struct{}{}
	// fmt.Println(<-leaderCh)
	// leaderCh <- struct{}{}
	// err = lock.Unlock()
	// if err != nil {
	// 	log.Fatal("err", err)
	// }
	// log.Println(sessionName, "unlock")
}
