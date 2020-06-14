package xnet

import (
	"net"
	"sync"
	"syscall"
)

type PacketConn struct {
	net.PacketConn
	sync.RWMutex
	closed bool
	fd     int
	proxy  bool
}

func (pc *PacketConn) Close() error {
	pc.Lock()
	if pc.closed {
		pc.Unlock()
		return nil
	}
	fd := pc.fd
	pc.closed = true
	pc.fd = unknownFD
	pc.Unlock()
	if fd != unknownFD {
		syscall.Shutdown(fd, syscall.SHUT_RDWR)
	}
	return pc.PacketConn.Close()
}

func (pc *PacketConn) GetFD() (fd int) {
	pc.RLock()
	fd = pc.fd
	pc.RUnlock()
	return
}

func (pc *PacketConn) IsClosed() (rv bool) {
	pc.RLock()
	rv = pc.closed
	pc.RUnlock()
	return
}
