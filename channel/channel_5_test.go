package main_test

import (
	"fmt"
	"testing"
	"time"
)

func channel_5() {
	ch := make(chan struct{})
	go func() {
		close(ch)
	}()

	for {
		time.Sleep(1 * time.Second)
		select {
		case v, ok := <-ch:
			fmt.Printf("value: %v; closed or not(false is closed): %v\n", v, ok)
		}
	}
}

func Test_channel_5(t *testing.T) {
	channel_5()
}
