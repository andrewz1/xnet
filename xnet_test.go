package xnet

import (
	"net"
	"syscall"
	"testing"
)

func TestListen(t *testing.T) {
	lst, err := Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("Listen: %s", err)
	}
	lst.Close()
}

func TestListenPacket(t *testing.T) {
	lst, err := ListenPacket("udp", ":8080")
	if err != nil {
		t.Fatalf("ListenPacket: %s", err)
	}
	lst.Close()
}

func TestInterface(t *testing.T) {
	c, err := net.Dial("tcp", "www.google.com:443")
	if err != nil {
		t.Fatalf("Dial: %s", err)
	}
	v, ok := c.(syscall.Conn)
	if !ok {
		t.Fatalf("Interface convart fail")
	}
	rc, err := v.SyscallConn()
	if err != nil {
		t.Fatalf("SyscallConn: %s", err)
	}
	t.Logf("%T", rc)
}
