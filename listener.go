package xnet

import (
	"net"
	"sync"
	"syscall"
)

type Listener struct {
	net.Listener
	sync.Mutex
	closed bool
	fd     int
}

func (l *Listener) Close() (err error) {
	l.Lock()
	if l.closed {
		l.Unlock()
		return
	}
	if l.fd != unknownFD {
		_ = syscall.Shutdown(l.fd, syscall.SHUT_RDWR)
		l.fd = unknownFD
	}
	if err = l.Listener.Close(); err != nil {
		l.Unlock()
		return
	}
	l.closed = true
	l.Unlock()
	return
}

func (l *Listener) xAccept() (xc *Conn, err error) {
	var (
		nc net.Conn
		rc syscall.RawConn
	)

	if nc, err = l.Listener.Accept(); err != nil {
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
	err = rc.Control(func(pfd uintptr) {
		xc.fd = int(pfd)
	})
	if err != nil {
		xc.Close()
		xc = nil
		return
	}
	return
}

func (l *Listener) Accept() (net.Conn, error) {
	return l.xAccept()
}

func (l *Listener) GetFD() (fd int) {
	l.Lock()
	fd = l.fd
	l.Unlock()
	return
}
