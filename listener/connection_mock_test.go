package listener

import (
	"net"
	"time"
)

type netConnMock struct{}

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
	return nil
}

func (ncm *netConnMock) RemoteAddr() net.Addr {
	return nil
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
