package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// CAS操作修改共享变量时候不需要对共享变量加锁
// 而是通过类似乐观锁的方式进行检查
// 本质还是不断的占用CPU资源换取加锁带来的开销

// 使用CAS来实现计数器

var (
	counter int32
	wg      sync.WaitGroup
)

func casStart() {
	threadNum := 100000
	wg.Add(threadNum)
	for i := 0; i < threadNum; i++ {
		go incCounter(i)
	}
	wg.Wait()
	fmt.Println(counter)
}

func incCounter(threadNum int) {
	defer wg.Done()
	repeatTimes := 0

	for {
		// 原子操作
		old := counter
		// 本质上是通过类似于乐观锁的形式来保证并发安全
		ok := atomic.CompareAndSwapInt32(&counter, old, old+1)
		// 如果*counter的值等于old, 则将*counter的值换成old+1, 放回true
		// 如果*counter的值不等于old, 则返回false
		// 以上两种情况均为原子操作
		if ok {
			break
		} else {
			repeatTimes++
		}
	}
	if repeatTimes == 0 {
		return
	}
	fmt.Printf("thread: %d repeats %d times\n", threadNum, repeatTimes)
}

func main() {
	casStart()
}
