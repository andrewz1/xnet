package xnet

import (
	"net"
	"sync"
	"syscall"
)

type Conn struct {
	net.Conn
	sync.RWMutex
	closed bool
	fd     int
}

func (cn *Conn) Close() error {
	cn.Lock()
	if cn.closed {
		cn.Unlock()
		return nil
	}
	fd := cn.fd
	cn.closed = true
	cn.fd = unknownFD
	cn.Unlock()
	if fd != unknownFD {
		syscall.Shutdown(fd, syscall.SHUT_RDWR)
	}
	return cn.Conn.Close()
}

func (cn *Conn) GetFD() (fd int) {
	cn.RLock()
	fd = cn.fd
	cn.RUnlock()
	return
}

func (cn *Conn) IsClosed() (rv bool) {
	cn.RLock()
	rv = cn.closed
	cn.RUnlock()
	return
}
