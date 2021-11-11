package mock

import (
	"net"
	"testing"
)

func TestAddrMockAndConnectionMock(t *testing.T) {
	var a net.Addr
	a = &NetAddr{}
	var c net.Conn
	c = &NetConn{a}
	c.Close()
}
