package xnet

import (
	"net"
	"time"
)

func tuneTCP(c *net.TCPConn) (*net.TCPConn, error) {
	_ = c.SetLinger(5)
	_ = c.SetKeepAlive(true)
	_ = c.SetKeepAlivePeriod(10 * time.Second)
	_ = c.SetNoDelay(true)
	return c, nil
}
