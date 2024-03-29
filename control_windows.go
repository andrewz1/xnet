package xnet

import (
	"syscall"

	"golang.org/x/sys/windows"
)

func Control(network, address string, c syscall.RawConn) (err error) {
	c.Control(func(fd uintptr) {
		err = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
	})
	return
}
