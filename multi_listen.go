package xnet

import (
	"net"
	"sync/atomic"
)

const (
	chanLen = 100
)

type MultiListener struct {
	closed uint32
	ln     []net.Listener
	ch     *SafeChan
}

func (ml *MultiListener) Accept() (net.Conn, error) {
	for cn := range ml.ch.C {
		return cn.(net.Conn), nil
	}
	return nil, net.ErrClosed
}

func (ml *MultiListener) closeListeners() {
	if ml == nil || len(ml.ln) == 0 {
		return
	}
	for _, ln := range ml.ln {
		ln.Close()
	}
	ml.ch.Close()
	for v := range ml.ch.C {
		v.(net.Conn).Close()
	}
}

func (ml *MultiListener) Close() error {
	if !atomic.CompareAndSwapUint32(&ml.closed, 0, 1) {
		return net.ErrClosed
	}
	ml.closeListeners()
	return nil
}

func (ml *MultiListener) Addr() net.Addr {
	if atomic.LoadUint32(&ml.closed) != 0 {
		return nil
	}
	return ml.ln[0].Addr() // return first listener addr
}

func (ml *MultiListener) acceptOne(ln net.Listener) {
	for {
		cn, err := ln.Accept()
		if err != nil {
			return
		}
		if !ml.ch.Send(cn) {
			return
		}
	}
}

func MultiListen(network string, address ...string) (*MultiListener, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
	default:
		return nil, net.UnknownNetworkError(network)
	}
	if len(address) == 0 {
		return nil, net.InvalidAddrError("invalid listen address")
	}
	ml := &MultiListener{
		ch: NewSafeChan(chanLen),
	}
	var err error
	defer func() {
		if err != nil {
			ml.closeListeners()
		}
	}()
	var ln net.Listener
	for _, a := range address {
		if ln, err = Listen(network, a); err != nil {
			return nil, err
		}
		ml.ln = append(ml.ln, ln)
	}
	for _, l := range ml.ln {
		go ml.acceptOne(l)
	}
	return ml, nil
}

func AggregateListen(lns ...net.Listener) net.Listener {
	if len(lns) == 0 {
		return nil
	}
	ml := &MultiListener{
		ln: append([]net.Listener{}, lns...),
		ch: NewSafeChan(chanLen),
	}
	for _, l := range ml.ln {
		go ml.acceptOne(l)
	}
	return ml
}
