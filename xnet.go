package xnet

import (
	"syscall"

	"golang.org/x/sys/unix"
)

const unknownFD int = -1

type rawConn struct {
	fd    int
	proxy bool
}

func newRawConn(proxy bool) *rawConn {
	return &rawConn{
		fd:    unknownFD,
		proxy: proxy,
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
	if err = syscall.SetsockoptInt(r.fd, syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
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

func setProxyFd(fd int) error {
	return syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT, 1)
}

func (r *rawConn) setProxy() (err error) {
	if !r.isOk() {
		return
	}
	err = setProxyFd(r.fd)
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
