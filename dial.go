package xnet

import (
	"context"
	"net"
	"time"
)

func newDialer(laddr net.Addr, tmo time.Duration) *net.Dialer {
	return &net.Dialer{
		Timeout:       tmo,
		LocalAddr:     laddr,
		FallbackDelay: 100 * time.Millisecond,
		KeepAlive:     10 * time.Second,
	}
}

func DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		ra, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return DialTCPContext(ctx, network, &net.TCPAddr{}, ra)
	case "udp", "udp4", "udp6":
		ra, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return DialUDPContext(ctx, network, &net.UDPAddr{}, ra)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}

func DialTimeout(network, address string, tmo time.Duration) (net.Conn, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		ra, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return DialTCPTimeout(network, &net.TCPAddr{}, ra, tmo)
	case "udp", "udp4", "udp6":
		ra, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return DialUDPTimeout(network, &net.UDPAddr{}, ra, tmo)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}

func Dial(network, address string) (net.Conn, error) {
	return DialTimeout(network, address, 0)
}
