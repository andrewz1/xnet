package xnet

import (
	"context"
	"net"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const unknownFD int = -1

func setReuse(fd int) (err error) {
	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return
	}
	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		return
	}
	return
}

func setNoDelay(fd int) (err error) {
	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
	return
}

func ListenCtx(ctx context.Context, network, address string) (*Listener, error) {
	fd := unknownFD
	lc := &net.ListenConfig{
		Control: func(network, address string, rc syscall.RawConn) error {
			if rc != nil {
				err := rc.Control(func(pfd uintptr) {
					fd = int(pfd)
				})
				if err != nil {
					return err
				}
			}
			if fd == unknownFD {
				return nil
			}
			if err := setReuse(fd); err != nil {
				return err
			}
			switch network {
			case "tcp":
			case "tcp4":
			case "tcp6":
			default:
				return nil
			}
			if err := setNoDelay(fd); err != nil {
				return err
			}
			return nil
		},
	}
	nl, err := lc.Listen(ctx, network, address)
	if err != nil {
		return nil, err
	}
	return &Listener{Listener: nl, fd: fd}, nil
}

func Listen(network, address string) (*Listener, error) {
	return ListenCtx(context.Background(), network, address)
}

func ListenPacketCtx(ctx context.Context, network, address string) (*PacketConn, error) {
	fd := unknownFD
	lc := &net.ListenConfig{
		Control: func(network, address string, rc syscall.RawConn) error {
			if rc != nil {
				err := rc.Control(func(pfd uintptr) {
					fd = int(pfd)
				})
				if err != nil {
					return err
				}
			}
			if fd == unknownFD {
				return nil
			}
			if err := setReuse(fd); err != nil {
				return err
			}
			switch network {
			case "tcp":
			case "tcp4":
			case "tcp6":
			default:
				return nil
			}
			if err := setNoDelay(fd); err != nil {
				return err
			}
			return nil
		},
	}
	pl, err := lc.ListenPacket(ctx, network, address)
	if err != nil {
		return nil, err
	}
	return &PacketConn{PacketConn: pl, fd: fd}, nil
}

func ListenPacket(network, address string) (*PacketConn, error) {
	return ListenPacketCtx(context.Background(), network, address)
}

func xdial(ctx context.Context, tmo time.Duration, network, address string) (*Conn, error) {
	fd := unknownFD
	dl := &net.Dialer{
		Timeout: tmo,
		Control: func(network, address string, rc syscall.RawConn) error {
			if rc != nil {
				err := rc.Control(func(pfd uintptr) {
					fd = int(pfd)
				})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	nc, err := dl.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	setLinger(nc)
	return &Conn{Conn: nc, fd: fd}, nil
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
