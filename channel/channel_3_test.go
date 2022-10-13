package main_test

import (
	"fmt"
	"testing"
	"time"
)

// 用channel来实现锁
// 一种方式是先初始化一个 capacity 等于 1 的 Channel，然后再放入一个元素。这个元素就代表锁，谁取得了这个元素，就相当于获取了这把锁。
// 另一种方式是，先初始化一个 capacity 等于 1 的 Channel，它的“空槽”代表锁，谁能成功地把元素发送到这个 Channel，谁就获取了这把锁。
// 这里实现第一种方式

type mutex struct {
	ch chan struct{}
}

// 初始化锁
func NewMutex() *mutex {
	mu := &mutex{ch: make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

// 获得锁
func (mu *mutex) Lock() {
	<-mu.ch
}

// 释放锁
func (mu *mutex) UnLock() {
	select {
	case mu.ch <- struct{}{}:
	default:
		panic("释放了没有被lock的锁")
	}
}

// 尝试获取锁
func (mu *mutex) TryLock() bool {
	select {
	case <-mu.ch:
		return true
	default:
		return false
	}
}

// 锁是否已经被持有
func (mu *mutex) IsLocked() bool {
	return len(mu.ch) == 0
}

// 加入持有锁不得超时的设置, 返回加锁是否成功
func (mu *mutex) LockTimeout(duration time.Duration) bool {
	timer := time.NewTimer(duration)
	select {
	case <-timer.C:
		return false
	case <-mu.ch:
		<-mu.ch
		timer.Stop()
		return true
	}
}

func Test_channel_3(t *testing.T) {
	m := NewMutex()
	ok := m.TryLock()
	fmt.Printf("locked v %v\n", ok)
	ok = m.LockTimeout(1 * time.Second)
	fmt.Printf("locked v %v\n", ok)
	m.UnLock()
}
