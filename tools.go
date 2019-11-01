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
	rc, err := v.SyscallConn()
	if err != nil {
		return nil
	}
	return rc
}

func setLinger(c interface{}) {
	if v := haveLingerCall(c); v != nil {
		v.SetLinger(0)
	}
}

func GetFD(c interface{}) int {
	v, ok := c.(haveFD)
	if !ok {
		return unknownFD
	}
	return v.GetFD()
}

func HaveFD(c interface{}) bool {
	_, ok := c.(haveFD)
	return ok
}
