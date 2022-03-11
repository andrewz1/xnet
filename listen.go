package xnet

import (
	"context"
	"net"
	"time"
)

type ListenConfig struct {
	net.ListenConfig
	ctx context.Context
}

func NewListenConfig(ctx context.Context) *ListenConfig {
	return &ListenConfig{
		ListenConfig: net.ListenConfig{
			KeepAlive: 10 * time.Second,
		},
		ctx: ctx,
	}
}

func (c *ListenConfig) Listen(network, address string) (net.Listener, error) {
	return c.ListenConfig.Listen(c.ctx, network, address)
}

func (c *ListenConfig) ListenPacket(network, address string) (net.PacketConn, error) {
	return c.ListenConfig.ListenPacket(c.ctx, network, address)
}

func (c *ListenConfig) ListenTCP(network string, laddr *net.TCPAddr) (*TCPListener, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
	case "":
		network = "tcp"
	default:
		return nil, net.UnknownNetworkError(network)
	}
	if laddr == nil {
		laddr = &net.TCPAddr{}
	}
	l, err := c.Listen(network, laddr.String())
	if err != nil {
		return nil, err
	}
	return &TCPListener{TCPListener: l.(*net.TCPListener)}, nil
}

func (c *ListenConfig) ListenUDP(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
	switch network {
	case "udp", "udp4", "udp6":
	case "":
		network = "udp"
	default:
		return nil, net.UnknownNetworkError(network)
	}
	if laddr == nil {
		laddr = &net.UDPAddr{}
	}
	cn, err := c.ListenPacket(network, laddr.String())
	if err != nil {
		return nil, err
	}
	return tuneUDP(cn.(*net.UDPConn)), nil
}

func ListenContext(ctx context.Context, network, address string) (net.Listener, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return NewListenConfig(ctx).Listen(network, address)
	default:
		return nil, net.UnknownNetworkError(network)
	}

}

func ListenPacketContext(ctx context.Context, network, address string) (net.PacketConn, error) {
	switch network {
	case "udp", "udp4", "udp6":
		return NewListenConfig(ctx).ListenPacket(network, address)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}

func Listen(network, address string) (net.Listener, error) {
	return ListenContext(context.Background(), network, address)
}

func ListenPacket(network, address string) (net.PacketConn, error) {
	return ListenPacketContext(context.Background(), network, address)
}
