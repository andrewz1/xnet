package xnet

import (
	"net"
)

type TCPListener struct {
	*net.TCPListener
}

func ListenTCP(network string, laddr *net.TCPAddr) (*TCPListener, error) {
	lc, err := net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}
	l := &TCPListener{lc}
	return l, nil
}

func (l *TCPListener) AcceptTCP() (*net.TCPConn, error) {
	tc, err := l.TCPListener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	return tuneTCP(tc)
}

func (l *TCPListener) Accept() (net.Conn, error) {
	return l.AcceptTCP()
}
