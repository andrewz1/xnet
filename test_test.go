package xnet

import (
	"testing"
	"time"
)

func TestAccept(t *testing.T) {
	ln, err := MultiListen("tcp", ":12345", ":1234")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			cn, err := ln.Accept()
			if err != nil {
				t.Error("accept: ", err)
				return
			}
			t.Log("accepted")
			cn.Close()
		}
	}()
	time.Sleep(time.Second)
	go func() {
		cn, err := Dial("tcp", "127.0.0.1:12345")
		if err != nil {
			t.Error("dial: ", err)
			return
		}
		t.Log("connected")
		cn.Close()
	}()

	go func() {
		cn, err := Dial("tcp", "127.0.0.1:1234")
		if err != nil {
			t.Error("dial: ", err)
			return
		}
		t.Log("connected")
		cn.Close()
	}()

	time.Sleep(time.Second)
	t.Error(ln.Close())
	time.Sleep(time.Second)
	t.Error(ln.Close())
}
