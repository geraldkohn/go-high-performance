package main_test

import (
	"fmt"
	"testing"
	"time"
)

// 测试 for range 遍历 channel 的时候，channel 关闭，for-range 是否会退出
func channel_7() {
	ch := make(chan interface{}, 1)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		// close(ch)
	}()

	go func() {
		for v := range ch {
			fmt.Println(v)
		}
		fmt.Println("退出")
	}()

	select {
	case <- time.After(2 * time.Second):
		return
	}
}

func Test7(t *testing.T) {
	channel_7()
}
