package xnet

import (
	"net"
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
	rc := getSyscallConn(c)
	if rc == nil {
		t.Fatalf("getSyscallConn error")
	}
	t.Logf("%T", rc)
}
