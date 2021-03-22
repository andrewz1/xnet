package xnet

import (
	"syscall"

	"golang.org/x/sys/unix"
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

func (r *rawConn) setReuse() (err error) {
	if err = syscall.SetsockoptInt(r.fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return
	}
	if err = syscall.SetsockoptInt(r.fd, syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		return
	}
	return
}

func (r *rawConn) setNoDelay(network string) (err error) {
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
	if rc == nil {
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
