package xnet

import (
	"sync"
)

type SafeChan struct {
	sync.RWMutex
	closed bool
	C      chan interface{}
}

func NewSafeChan(n int) *SafeChan {
	return &SafeChan{
		C: make(chan interface{}, n),
	}
}

func (c *SafeChan) Close() {
	c.Lock()
	defer c.Unlock()
	if !c.closed {
		c.closed = true
		close(c.C)
	}
}

func (c *SafeChan) Send(v interface{}) bool {
	c.RLock()
	defer c.RUnlock()
	if c.closed {
		return false
	}
	c.C <- v
	return true
}
