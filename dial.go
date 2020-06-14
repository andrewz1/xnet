package xnet

import (
	"context"
	"net"
	"syscall"
	"time"
)

func (r *rawConn) ctrlDial(network, address string, rc syscall.RawConn) (err error) {
	_, err = r.ctrlBase(rc)
	return
}

func xdial(ctx context.Context, tmo time.Duration, network, address string) (xc *Conn, err error) {
	r := newRawConn(false)
	dl := &net.Dialer{
		Timeout: tmo,
		Control: r.ctrlDial,
	}
	var nc net.Conn
	if nc, err = dl.DialContext(ctx, network, address); err != nil {
		return
	}
	setLinger(nc)
	xc = &Conn{
		Conn: nc,
		fd:   r.fd,
	}
	return
}

func DialCtx(ctx context.Context, network, address string) (*Conn, error) {
	return xdial(ctx, time.Duration(0), network, address)
}

func Dial(network, address string) (*Conn, error) {
	return xdial(context.Background(), time.Duration(0), network, address)
}

func DialTimeout(network, address string, timeout time.Duration) (*Conn, error) {
	return xdial(context.Background(), timeout, network, address)
}
