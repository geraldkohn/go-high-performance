package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type config struct {
	Addr string
	Count int32
}

func loadNewConfig() config {
	return config{
		Addr: "ShangHai",
		Count: rand.Int31(),
	}
}

func loadNewConfigGoroutine() {
	// ... 业务代码
	var configAtomic atomic.Value
	configAtomic.Store(loadNewConfig())
	var cond = sync.NewCond(&sync.Mutex{})
	go func () {
		for {
			time.Sleep(1 * time.Second)
			configAtomic.Store(loadNewConfig())
			cond.Broadcast() // 通知配置已经变更
		}
	}()

	go func () {
		for {
			cond.L.Lock()
			cond.Wait()
			c := configAtomic.Load().(config)
			fmt.Printf("new config: %+v\n", c)
			cond.L.Unlock()
		}
	}()

	select{}
}

// func main() {
// 	loadNewConfigGoroutine()
// }