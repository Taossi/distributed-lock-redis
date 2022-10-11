package main

import (
	"sync"
	"taossi/distributed-lock/redis"
	"time"

	"github.com/google/uuid"
)

/**
测试：模拟我在工作中遇到过的业务场景: 多个服务实例中多个线程同时进行发起http请求，需要进行过滤限频，使得同个id只发送一次http
	 这里简单模拟向redis里写"hello"，确保200个线程做到只写入一次(打印hello)
	 实际业务场景中会更复杂，可能会包含http请求，因此可能会存在请求超时等使得锁失效的场景。
**/

func main() {
	redis.CreateClient()
	key := "test"
	redis.DelHello(key)

	w := sync.WaitGroup{}
	for i := 0; i < 200; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			ReportHello(key)
		}()
	}
	w.Wait()
}

func ReportHello(key string) {
	lock := redis.NewRedisLock("test1", uuid.New().String(), time.Second*5)
	locked, err := lock.TryLock() // lock
	if err != nil || !locked {
		return
	}
	defer lock.UnLock() // unlock

	value, _ := redis.GetHello(key)

	if value != "hello" {
		_ = redis.SetHello(key) // 限频 只set一次
		return
	}
	return
}
