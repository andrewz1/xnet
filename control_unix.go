package xnet

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func Control(network, address string, c syscall.RawConn) (err error) {
	c.Control(func(fd uintptr) {
		if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
			return
		}
		if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
			return
		}
	})
	return
}
