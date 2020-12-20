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

func makeLocalAddr(network, laStr string) (la net.Addr) {
	if len(network) < 3 {
		return
	}
	switch network[0:3] {
	case "udp":
		la, _ = net.ResolveUDPAddr(network, laStr)
	case "tcp":
		la, _ = net.ResolveTCPAddr(network, laStr)
	}
	return
}

func xdial(ctx context.Context, tmo time.Duration, network, la, address string) (xc *Conn, err error) {
	r := newRawConn(false)
	dl := &net.Dialer{
		Timeout:   tmo,
		LocalAddr: makeLocalAddr(network, la),
		Control:   r.ctrlDial,
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
	return xdial(ctx, time.Duration(0), network, "", address)
}

func Dial(network, address string) (*Conn, error) {
	return xdial(context.Background(), time.Duration(0), network, "", address)
}

func DialTimeout(network, address string, timeout time.Duration) (*Conn, error) {
	return xdial(context.Background(), timeout, network, "", address)
}

func Dial2Ctx(ctx context.Context, network, la, address string) (*Conn, error) {
	return xdial(ctx, time.Duration(0), network, la, address)
}

func Dial2(network, la, address string) (*Conn, error) {
	return xdial(context.Background(), time.Duration(0), network, la, address)
}

func Dial2Timeout(network, la, address string, timeout time.Duration) (*Conn, error) {
	return xdial(context.Background(), timeout, network, la, address)
}
