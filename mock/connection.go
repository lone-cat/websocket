package mock

import (
	"net"
	"time"
)

// net.Addr mock
type NetAddr struct {
	Addr string
}

func (nam *NetAddr) Network() string {
	return `tcp`
}

func (nam *NetAddr) String() string {
	return nam.Addr
}

// net.Conn mock
type NetConn struct {
	Addr net.Addr
}

func (ncm *NetConn) Read(buf []byte) (int, error) {
	return 0, nil
}

func (ncm *NetConn) Write(buf []byte) (int, error) {
	return 0, nil
}

func (ncm *NetConn) Close() error {
	return nil
}

func (ncm *NetConn) LocalAddr() net.Addr {
	return ncm.Addr
}

func (ncm *NetConn) RemoteAddr() net.Addr {
	return ncm.Addr
}

func (ncm *NetConn) SetDeadline(t time.Time) error {
	return nil
}

func (ncm *NetConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (ncm *NetConn) SetWriteDeadline(t time.Time) error {
	return nil
}
