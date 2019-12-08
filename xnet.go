package xnet

import (
	"context"
	"net"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const unknownFD int = -1

type rawConn struct {
	fd int
}

func newRawConn() *rawConn {
	return &rawConn{
		fd: unknownFD,
	}
}

func (r *rawConn) setReuse() (err error) {
	if r == nil {
		return
	}
	if err = syscall.SetsockoptInt(r.fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return
	}
	if err = syscall.SetsockoptInt(r.fd, syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		return
	}
	return
}

func (r *rawConn) setNoDelay() (err error) {
	if r == nil {
		return
	}
	err = syscall.SetsockoptInt(r.fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
	return
}

func (r *rawConn) ctrl(network, address string, rc syscall.RawConn) (err error) {
	if r == nil || rc == nil {
		return
	}
	if err = rc.Control(func(pfd uintptr) { r.fd = int(pfd) }); err != nil {
		return
	}
	if r.fd == unknownFD {
		return nil
	}
	if err = r.setReuse(); err != nil {
		return
	}
	switch network {
	case "tcp":
	case "tcp4":
	case "tcp6":
	default:
		return nil
	}
	if err = r.setNoDelay(); err != nil {
		return
	}
	return
}

func (r *rawConn) dctrl(network, address string, rc syscall.RawConn) (err error) {
	if r == nil || rc == nil {
		return
	}
	if err = rc.Control(func(pfd uintptr) { r.fd = int(pfd) }); err != nil {
		return
	}
	return
}

func ListenCtx(ctx context.Context, network, address string) (xl *Listener, err error) {
	r := newRawConn()
	lc := &net.ListenConfig{
		Control: r.ctrl,
	}
	var nl net.Listener
	if nl, err = lc.Listen(ctx, network, address); err != nil {
		return
	}
	xl = &Listener{
		Listener: nl,
		fd:       r.fd,
	}
	return
}

func Listen(network, address string) (*Listener, error) {
	return ListenCtx(context.Background(), network, address)
}

func ListenPacketCtx(ctx context.Context, network, address string) (xl *PacketConn, err error) {
	r := newRawConn()
	lc := &net.ListenConfig{
		Control: r.ctrl,
	}
	var pl net.PacketConn
	if pl, err = lc.ListenPacket(ctx, network, address); err != nil {
		return
	}
	xl = &PacketConn{
		PacketConn: pl,
		fd:         r.fd,
	}
	return
}

func ListenPacket(network, address string) (*PacketConn, error) {
	return ListenPacketCtx(context.Background(), network, address)
}

func xdial(ctx context.Context, tmo time.Duration, network, address string) (xc *Conn, err error) {
	r := newRawConn()
	dl := &net.Dialer{
		Timeout: tmo,
		Control: r.dctrl,
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
