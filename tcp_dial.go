package xnet

import (
	"context"
	"net"
	"time"
)

func fixTcpNet(network string) (string, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return network, nil
	case "":
		return "tcp", nil
	default:
		return "", net.UnknownNetworkError(network)
	}
}

func DialTCPContext(ctx context.Context, network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	var err error
	if network, err = fixTcpNet(network); err != nil {
		return nil, err
	}
	dl := newDialer(laddr, 0)
	var c net.Conn
	if c, err = dl.DialContext(ctx, network, raddr.String()); err != nil {
		return nil, err
	}
	return tuneTCP(c.(*net.TCPConn))
}

func DialTCPTimeout(network string, laddr, raddr *net.TCPAddr, tmo time.Duration) (*net.TCPConn, error) {
	var err error
	if network, err = fixTcpNet(network); err != nil {
		return nil, err
	}
	dl := newDialer(laddr, tmo)
	var c net.Conn
	if c, err = dl.Dial(network, raddr.String()); err != nil {
		return nil, err
	}
	return tuneTCP(c.(*net.TCPConn))
}

func DialTCP(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	return DialTCPTimeout(network, laddr, raddr, 0)
}
