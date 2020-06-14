package xnet

import "syscall"

func getSyscallConn(c interface{}) syscall.RawConn {
	var (
		v  syscall.Conn
		ok bool
	)
	if v, ok = c.(syscall.Conn); !ok || v == nil {
		return nil

	}
	if rc, err := v.SyscallConn(); err == nil {
		return rc
	}
	return nil
}

func setLinger(c interface{}) {
	if v, ok := c.(interface{ SetLinger(int) error }); ok && v != nil {
		v.SetLinger(0)
	}
}

func GetFD(c interface{}) int {
	if v, ok := c.(interface{ GetFD() int }); ok && v != nil {
		return v.GetFD()
	}
	return unknownFD
}

func HaveFD(c interface{}) bool {
	v, ok := c.(interface{ GetFD() int })
	return ok && v != nil
}
