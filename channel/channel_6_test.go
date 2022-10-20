package main_test

import (
	"reflect"
)

// channel 任务编排场景

// 1. Or-Done 模式, 实现某个任务完成之后的信号通知
// 以下的例子是模拟一个请求发送到多个节点，只要有一个节点任务结束即返回。
func orDone(chs ...<-chan interface{}) <-chan interface{} {
	if len(chs) == 0 {
		return nil
	} else if len(chs) == 1 {
		return chs[0] // 只有一个 channel
	}

	orDoneChan := make(chan interface{})
	go func() {
		// 使用反射来处理 channel
		var cases []reflect.SelectCase
		for _, c := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		reflect.Select(cases) // 使用反射监听 channel
		close(orDoneChan)     // 关闭channel，传递节点已经结束的信号
	}()
	return orDoneChan
}

// 2. 扇入模式，监听多个 channel，将收到的所有值都发送到一个 channel 中
func fanInReflect(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	// 先使用反射监听多个 channel
	go func() {
		var cases []reflect.SelectCase
		for _, c := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		chosen, resv, resvOk := reflect.Select(cases) // 监听
		for len(cases) > 0 {
			if !resvOk { // channel 已经关闭
				cases = append(cases[:chosen], cases[chosen+1:]...)
				continue
			}
			out <- resv.Interface()
		}
	}()
	return out
}

// 3. 扇出模式，将一个 channel 中的值发送到多个 channel 中。
func fanOutReflect(in <-chan interface{}) []chan interface{} {
	numberOfOut := 10
	out := make([]chan interface{}, numberOfOut)
	go func() {
		for {
			select {
			case v, ok := <-in:
				if !ok { // 监听的 channel 关闭
					for _, c := range out {
						close(c)
					}
					return
				}
				for _, c := range out {
					c <- v
				}
			}
		}
	}()
	return out
}

// map-reduce
func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{})
	if in == nil { // 处理异常
		close(out)
		return out
	}

	go func() {
		for v := range in {
			out <- fn(v)
		}
		close(out) // 如果执行到这里，那么说明 in 被关闭。
	}()

	return out
}
