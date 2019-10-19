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

func (l *Listener) Close() error {
	l.Lock()
	defer l.Unlock()
	if l.closed {
		return nil
	}
	if l.fd != unknownFD {
		syscall.Shutdown(l.fd, syscall.SHUT_RDWR)
		l.fd = unknownFD
	}
	if err := l.Listener.Close(); err != nil {
		return err
	}
	l.closed = true
	return nil
}

func (l *Listener) Accept() (net.Conn, error) {
	nc, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	setLinger(nc)
	xc := &Conn{Conn: nc, fd: unknownFD}
	rc := getSyscallConn(nc)
	if rc == nil {
		return xc, nil
	}
	err = rc.Control(func(pfd uintptr) {
		xc.fd = int(pfd)
	})
	if err != nil {
		xc.Close()
		return nil, err
	}
	return xc, nil
}

func (l *Listener) GetFD() (fd int) {
	l.Lock()
	fd = l.fd
	l.Unlock()
	return
}
