package xnet

import (
	"context"
	"net"
	"syscall"
)

func (r *rawConn) ctrlListen(network, address string, rc syscall.RawConn) (err error) {
	var ok bool
	if ok, err = r.ctrlBase(rc); err != nil {
		return
	}
	if !ok {
		return
	}
	if err = r.setReuse(); err != nil {
		return
	}
	if err = r.setNoDelay(network); err != nil {
		return
	}
	return
}

func ListenCtx(ctx context.Context, network, address string) (xl *Listener, err error) {
	r := newRawConn()
	lc := &net.ListenConfig{
		Control: r.ctrlListen,
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
		Control: r.ctrlListen,
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
