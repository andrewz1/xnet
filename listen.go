package xnet

import (
	"net"
)

func Listen(network, address string) (net.Listener, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		la, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return ListenTCP(network, la)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}

func ListenPacket(network, address string) (net.PacketConn, error) {
	switch network {
	case "udp", "udp4", "udp6":
		la, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return ListenUDP(network, la)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}
