package main

import (
	"fmt"
	"sync"
	"time"
)

/**
* 设计:
* Write-preferring：
* 写优先的设计意味着，如果已经有一个 writer 在等待请求锁的话，它会阻止新来的请求锁的 reader 获取到锁，所以优先保障 writer。
* 当然，如果有一些 reader 已经请求了锁的话，新请求的 writer 也会等待已经存在的 reader 都释放锁之后才能获取。
* 所以，写优先级设计中的优先权是针对新来的请求而言的。这种设计主要避免了 writer 的饥饿问题。

* 踩坑点:
* 1. 不可复制
* 2. 重入导致死锁: 重复加写锁, 加读锁不释放锁继续加写锁(这是常见的读写锁死锁场景, 写锁阻塞等待读锁释放, 读锁不可能被释放)
* 3. 非常隐蔽的死锁场景:
  当一个 writer 请求锁的时候，如果已经有一些活跃的 reader，它会等待这些活跃的 reader 完成，才有可能获取到锁，但是，如果之后活跃的 reader 再依赖新的 reader 的话，这些新的 reader 就会等待 writer 释放锁之后才能继续执行，这就形成了一个环形依赖： writer 依赖活跃的 reader -> 活跃的 reader 依赖新来的 reader -> 新来的 reader 依赖 writer。
*/

// 第三个场景的例子
func deadLockExample() {
	var mu sync.RWMutex

	// writer,稍微等待，然后制造一个调用Lock的场景
	go func() {
		time.Sleep(200 * time.Millisecond)
		mu.Lock()
		fmt.Println("Lock")
		time.Sleep(100 * time.Millisecond)
		mu.Unlock()
		fmt.Println("Unlock")
	}()

	go func() {
		factorial(&mu, 10) // 计算10的阶乘, 10!
	}()

	select {}
}

// 递归调用计算阶乘
func factorial(m *sync.RWMutex, n int) int {
	if n < 1 { // 阶乘退出条件
		return 0
	}
	fmt.Println("RLock")
	m.RLock()
	defer func() {
		fmt.Println("RUnlock")
		m.RUnlock()
	}()
	time.Sleep(100 * time.Millisecond)
	return factorial(m, n-1) * n // 递归调用
}

func main() {
	deadLockExample()
}

// 所以读锁在递归调用读锁之前, 必须解锁! 否则在递归调用前一个goroutine调用写锁会产生死锁.
// 总之, 尽量避免重入.
