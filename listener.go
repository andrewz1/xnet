package xnet

import (
	"net"
	"sync"
	"syscall"
)

type Listener struct {
	net.Listener
	sync.RWMutex
	closed bool
	fd     int
	proxy  bool
}

func (ls *Listener) Close() error {
	ls.Lock()
	if ls.closed {
		ls.Unlock()
		return nil
	}
	fd := ls.fd
	ls.closed = true
	ls.fd = unknownFD
	ls.Unlock()
	if fd != unknownFD {
		syscall.Shutdown(fd, syscall.SHUT_RDWR)
	}
	return ls.Listener.Close()
}

func (ls *Listener) AcceptXConn() (xc *Conn, err error) {
	var (
		nc net.Conn
		rc syscall.RawConn
	)

	if nc, err = ls.Listener.Accept(); err != nil {
		return
	}
	setLinger(nc)
	xc = &Conn{
		Conn: nc,
		fd:   unknownFD,
	}
	if rc = getSyscallConn(nc); rc == nil {
		return
	}
	if err = rc.Control(func(pfd uintptr) { xc.fd = int(pfd) }); err != nil {
		xc.Close()
		xc = nil
		return
	}
	if xc.fd < 0 {
		return
	}
	if ls.proxy {
		if err = setProxyFd(xc.fd); err != nil {
			xc.Close()
			xc = nil
			return
		}
	}
	return
}

func (ls *Listener) Accept() (net.Conn, error) {
	return ls.AcceptXConn()
}

func (ls *Listener) GetFD() (fd int) {
	ls.RLock()
	fd = ls.fd
	ls.RUnlock()
	return
}

func (ls *Listener) IsClosed() (rv bool) {
	ls.RLock()
	rv = ls.closed
	ls.RUnlock()
	return
}
