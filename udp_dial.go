package xnet

import (
	"context"
	"net"
	"time"
)

func fixUdpNet(network string) (string, error) {
	switch network {
	case "udp", "udp4", "udp6":
		return network, nil
	case "":
		return "udp", nil
	default:
		return "", net.UnknownNetworkError(network)
	}
}

func DialUDPContext(ctx context.Context, network string, laddr, raddr *net.UDPAddr) (*net.UDPConn, error) {
	var err error
	if network, err = fixUdpNet(network); err != nil {
		return nil, err
	}
	dl := newDialer(laddr, 0)
	var c net.Conn
	if c, err = dl.DialContext(ctx, network, raddr.String()); err != nil {
		return nil, err
	}
	return tuneUDP(c.(*net.UDPConn)), nil
}

func DialUDPTimeout(network string, laddr, raddr *net.UDPAddr, tmo time.Duration) (*net.UDPConn, error) {
	var err error
	if network, err = fixUdpNet(network); err != nil {
		return nil, err
	}
	dl := newDialer(laddr, tmo)
	var c net.Conn
	if c, err = dl.Dial(network, raddr.String()); err != nil {
		return nil, err
	}
	return tuneUDP(c.(*net.UDPConn)), nil
}

func DialUDP(network string, laddr, raddr *net.UDPAddr) (*net.UDPConn, error) {
	return DialUDPTimeout(network, laddr, raddr, 0)
}
