package listener

import (
	"net"
	"time"
)

type netAddrMock struct {
	addr string
}

func (nam *netAddrMock) Network() string {
	return `tcp`
}

func (nam *netAddrMock) String() string {
	return nam.addr
}

type netConnMock struct {
	addr netAddrMock
}

func (ncm *netConnMock) Read(buf []byte) (int, error) {
	return 0, nil
}

func (ncm *netConnMock) Write(buf []byte) (int, error) {
	return 0, nil
}

func (ncm *netConnMock) Close() error {
	return nil
}

func (ncm *netConnMock) LocalAddr() net.Addr {
	return &ncm.addr
}

func (ncm *netConnMock) RemoteAddr() net.Addr {
	return &ncm.addr
}

func (ncm *netConnMock) SetDeadline(t time.Time) error {
	return nil
}

func (ncm *netConnMock) SetReadDeadline(t time.Time) error {
	return nil
}

func (ncm *netConnMock) SetWriteDeadline(t time.Time) error {
	return nil
}
