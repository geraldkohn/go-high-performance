package main_test

import (
	"log"
	"testing"
	"time"
)

// 多生产者，多消费者模型（消息队列）

func producer(name string, ch chan int, isClosed chan struct{}) {
	for i := 0; i < 10; i++ {
		select {
		case <-isClosed: // 如果通道已经被其他生产者关闭了，就不能生产了
			return
		default: // 通道没有关闭，则继续生产
			log.Printf("生产者 %s 生产值 %d\n", name, i)
			ch <- i
		}
	}
	isClosed <- struct{}{} // 通知其他生产者通道要关闭了
	close(ch)
}

func consumer(name string, ch <-chan int) {
	for v := range ch {
		log.Printf("消费者 %s 接收到值 %d\n", name, v)
	}
	log.Println("channel 已经关闭")
}

func Test8(t *testing.T) {
	isClosed := make(chan struct{})
	ch := make(chan int, 5)
	go producer("1", ch, isClosed)
	go producer("2", ch, isClosed)
	go consumer("1", ch)
	go consumer("2", ch)
	select {
	case <-time.After(2 * time.Second):
		return
	}
}
