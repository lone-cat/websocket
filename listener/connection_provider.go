package listener

import (
	"net"
	"runtime"
)

type ConnectionProvider struct {
	listener    ConnectionProviderI
	middlewares []ConnectionMiddlewareI
	channels    []chan net.Conn
}

func NewConnectionProvider(listener ConnectionProviderI, middlewares ...ConnectionMiddlewareI) *ConnectionProvider {
	channels := make([]chan net.Conn, len(middlewares)+1)

	for i := range channels {
		channels[i] = make(chan net.Conn)
	}

	return &ConnectionProvider{
		listener:    listener,
		middlewares: middlewares,
		channels:    channels,
	}
}

func (cp *ConnectionProvider) GetResultChan() (resultChan <-chan net.Conn) {
	return cp.channels[len(cp.channels)-1]
}

func (cp *ConnectionProvider) Start() (err error) {
	for i := len(cp.middlewares) - 1; i >= 0; i-- {
		err = cp.middlewares[i].StartAsync(cp.channels[i], cp.channels[i+1])
		if err != nil {
			return
		}
		runtime.Gosched()
	}

	err = cp.listener.StartAsync(cp.channels[0])
	if err != nil {
	}
	runtime.Gosched()

	return
}

func (cp *ConnectionProvider) Stop() {
	cp.listener.StopSync()
	for _, middleware := range cp.middlewares {
		middleware.StopSync()
	}
}
