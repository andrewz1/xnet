package xnet

import (
	"net"
	"sync"
	"syscall"
)

type Conn struct {
	net.Conn
	sync.Mutex
	closed bool
	fd     int
}

func (c *Conn) Close() (err error) {
	c.Lock()
	if c.closed {
		c.Unlock()
		return
	}
	if c.fd != unknownFD {
		_ = syscall.Shutdown(c.fd, syscall.SHUT_RDWR)
		c.fd = unknownFD
	}
	if err = c.Conn.Close(); err != nil {
		c.Unlock()
		return
	}
	c.closed = true
	c.Unlock()
	return
}

func (c *Conn) GetFD() (fd int) {
	c.Lock()
	fd = c.fd
	c.Unlock()
	return
}
