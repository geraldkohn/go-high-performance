package main_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 测试channel关闭之后仍然能读取值.
// 这里要注意的是: channel不是并发安全的, 多个goroutine读取同一个channel可能导致值的丢失.

var ch chan interface{}
var mu sync.Mutex

func test_channel() {
	mu.Lock()
	i, isClosed := <-ch
	mu.Unlock()

	fmt.Println(i)
	if isClosed {
		fmt.Println("closed")
	} else {
		fmt.Println("not closed")
	}
}

func produce_channel() {
	ch = make(chan interface{}, 2)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	close(ch)
}

func Test_channel_1(t *testing.T) {
	go produce_channel()
	for i := 0; i < 4; i++ {
		go test_channel()
	}
	select {
	case <-time.After(10 * time.Second):
		return
	}
}
