package main_test

import (
	"fmt"
	"testing"
	"time"
)

// 有一道经典的使用 Channel 进行任务编排的题，你可以尝试做一下：
// 有四个 goroutine，编号为 1、2、3、4。每秒钟会有一个 goroutine 打印出它自己的编号，要求你编写一个程序，让输出的编号总是按照 1、2、3、4、1、2、3、4、……的顺序打印出来。

func channel_2() {
	chanArr := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}

	for i := 0; i < 4; i++ {
		go func(i int) {
			for {
				<-chanArr[i] // 监听第i个channel, 没有值的时候阻塞住了
				fmt.Printf("I am %d\n", i)
				time.Sleep(1 * time.Second)
				// 往下一个应该被执行的goroutine监听的channel中传值, 让下一个goroutine由阻塞变为正在执行
				// 因为第3个goroutine应该唤醒第0个goroutine, 所以应该是一个类似的环路.
				chanArr[(i+1)%4] <- struct{}{}
			}
		}(i)
	}

	// 唤醒第0个goroutine
	chanArr[0] <- struct{}{}
	select {
	case <- time.After(12 * time.Second):
		return
	}
}

func Test_channel_2(t *testing.T) {
	channel_2()
}
