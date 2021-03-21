package xnet

import (
	"syscall"
)

const unknownFD int = -1

type rawConn struct {
	fd int
}

func newRawConn() *rawConn {
	return &rawConn{
		fd: unknownFD,
	}
}

func (r *rawConn) isInit() bool {
	return r != nil
}

func (r *rawConn) isOk() bool {
	return r.isInit() && r.fd >= 0
}

func (r *rawConn) setReuse() (err error) {
	if !r.isOk() {
		return
	}
	if err = syscall.SetsockoptInt(r.fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return
	}
	if err = syscall.SetsockoptInt(r.fd, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1); err != nil {
		return
	}
	return
}

func (r *rawConn) setNoDelay(network string) (err error) {
	if !r.isOk() {
		return
	}
	switch network {
	case "tcp":
	case "tcp4":
	case "tcp6":
	default:
		return
	}
	err = syscall.SetsockoptInt(r.fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
	return
}

func (r *rawConn) ctrlBase(rc syscall.RawConn) (ok bool, err error) {
	if !r.isInit() || rc == nil {
		return
	}
	if err = rc.Control(func(pfd uintptr) { r.fd = int(pfd) }); err != nil {
		return
	}
	if r.fd == unknownFD {
		return
	}
	ok = true
	return
}
