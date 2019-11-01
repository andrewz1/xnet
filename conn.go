package xnet

import (
	"net"
	"sync"
)

type Conn struct {
	net.Conn
	sync.Mutex
	closed bool
	fd     int
}

func (c *Conn) Close() error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil
	}
	//if c.fd != unknownFD {
	//	syscall.Shutdown(c.fd, syscall.SHUT_RDWR)
	//	c.fd = unknownFD
	//}
	if err := c.Conn.Close(); err != nil {
		return err
	}
	c.closed = true
	return nil
}

func (c *Conn) GetFD() (fd int) {
	c.Lock()
	fd = c.fd
	c.Unlock()
	return
}
