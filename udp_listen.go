package xnet

import (
	"net"
)

func ListenUDP(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
	uc, err := net.ListenUDP(network, laddr)
	if err != nil {
		return nil, err
	}
	return tuneUDP(uc), nil
}
