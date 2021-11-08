package listener

import "net"

type ConnectionProviderI interface {
	StartAsync(resultChan chan<- net.Conn) error
	StopSync()
}

type ConnectionMiddlewareI interface {
	StartAsync(chanFrom <-chan net.Conn, chanTo chan<- net.Conn) error
	StopAsync() <-chan struct{}
	StopSync()
}
