package main_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// 使用反射来处理channel
// 动态处理两个 chan 的情形。因为这样的方式可以动态处理 case 数据
// 所以，你可以传入几百几千几万的 chan，这就解决了不能动态处理 n 个 chan 的问题。
// 首先，createCases 函数分别为每个 chan 生成了 recv case 和 send case，并返回一个 reflect.SelectCase 数组。
// 然后，通过一个循环的 for 循环执行 reflect.Select，这个方法会从 cases 中选择一个 case 执行。
// 第一次肯定是 send case，因为此时 chan 还没有元素，recv 还不可用。等 chan 中有了数据以后，recv case 就可以被选择了。这样，你就可以处理不定数量的 chan 了。

func channel_4() {
	var ch1 = make(chan int, 10)
	var ch2 = make(chan int, 10)

	// 创建SelectCase
	var cases = createCases(ch1, ch2)
	ch1 <- 10
	ch2 <- 10 

	go func () {
		time.Sleep(1 * time.Second)
		close(ch1)
		close(ch2)
	}()

	for len(cases) > 0 {
		chosen, resv, ok := reflect.Select(cases)
		if !ok {
			// 此通道关闭并且它的缓冲队列中为空
			cases = append(cases[:chosen], cases[chosen+1:]...)
			fmt.Printf("%d closed\n", chosen)
			continue
		}
		
		time.Sleep(200 * time.Millisecond)
		fmt.Println(chosen, resv, ok)
	}
}

func createCases(chs ...chan int) []reflect.SelectCase {
	var cases []reflect.SelectCase

	// 创建receive case, 与监听channel数据一样
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	return cases
}

func Test_channel_4(t *testing.T) {
	channel_4()
}
