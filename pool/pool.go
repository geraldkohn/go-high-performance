package pool

import (
	"errors"
	"io"
	"sync"
)

// 因为TCP三次握手的原因, 建立一个连接是成本比较高的行为. 所以在一个需要多次与特定实体交互的程序中, 就需要维持一个连接池, 里面有可以复用的连接可供使用.

type Pool struct {
	m       sync.Mutex                // 保证多个Goroutine访问的时候, close的线程安全
	res     chan io.Closer            // 连接存储的chan
	factory func() (io.Closer, error) //新建连接的工厂方法
	closed  bool                      // 连接关闭的标志
}

func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size is too small")
	}
	return &Pool{
		factory: fn,
		res:     make(chan io.Closer, size),
	}, nil
}
