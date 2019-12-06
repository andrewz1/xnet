package xnet

import (
	"net"
	"sync"
	"syscall"
)

type PacketConn struct {
	net.PacketConn
	sync.Mutex
	closed bool
	fd     int
}

func (pc *PacketConn) Close() (err error) {
	pc.Lock()
	if pc.closed {
		pc.Unlock()
		return
	}
	if pc.fd != unknownFD {
		syscall.Shutdown(pc.fd, syscall.SHUT_RDWR)
		pc.fd = unknownFD
	}
	if err = pc.PacketConn.Close(); err != nil {
		pc.Unlock()
		return
	}
	pc.closed = true
	pc.Unlock()
	return
}

func (pc *PacketConn) GetFD() (fd int) {
	pc.Lock()
	fd = pc.fd
	pc.Unlock()
	return
}
