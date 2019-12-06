package xnet

import "syscall"

type setLingerIf interface {
	SetLinger(sec int) error
}

type haveFD interface {
	GetFD() int
}

func haveLingerCall(v interface{}) setLingerIf {
	if l, ok := v.(setLingerIf); ok {
		return l
	}
	return nil
}

func getSyscallConn(c interface{}) syscall.RawConn {
	v, ok := c.(syscall.Conn)
	if !ok {
		return nil
	}
	if rc, err := v.SyscallConn(); err == nil {
		return rc
	}
	return nil
}

func setLinger(c interface{}) {
	if v := haveLingerCall(c); v != nil {
		v.SetLinger(0)
	}
}

func GetFD(c interface{}) int {
	if v, ok := c.(haveFD); ok {
		return v.GetFD()
	}
	return unknownFD
}

func HaveFD(c interface{}) bool {
	_, ok := c.(haveFD)
	return ok
}
