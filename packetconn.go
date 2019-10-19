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

func (pc *PacketConn) Close() error {
	pc.Lock()
	defer pc.Unlock()
	if pc.closed {
		return nil
	}
	if pc.fd != unknownFD {
		syscall.Shutdown(pc.fd, syscall.SHUT_RDWR)
		pc.fd = unknownFD
	}
	if err := pc.PacketConn.Close(); err != nil {
		return err
	}
	pc.closed = true
	return nil
}

func (pc *PacketConn) GetFD() (fd int) {
	pc.Lock()
	fd = pc.fd
	pc.Unlock()
	return
}
