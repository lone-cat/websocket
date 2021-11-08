package listener

import (
	"net"
)

type ConnectionProvider struct {
	listener             ConnectionProviderI
	debouncer            ConnectionMiddlewareI
	limiter              ConnectionMiddlewareI
	listenerResultsChan  chan net.Conn
	debouncerResultsChan chan net.Conn
	resultChan           chan net.Conn
}

func NewConnectionProvider(listener ConnectionProviderI, debouncer ConnectionMiddlewareI, limiter ConnectionMiddlewareI) *ConnectionProvider {
	return &ConnectionProvider{
		listener:             listener,
		debouncer:            debouncer,
		limiter:              limiter,
		listenerResultsChan:  make(chan net.Conn),
		debouncerResultsChan: make(chan net.Conn),
		resultChan:           make(chan net.Conn),
	}
}

func (cp *ConnectionProvider) GetResultChan() (resultChan <-chan net.Conn) {
	return cp.resultChan
}

func (cp *ConnectionProvider) Start() (err error) {
	err = cp.limiter.StartAsync(cp.debouncerResultsChan, cp.resultChan)
	if err != nil {
		return
	}

	err = cp.debouncer.StartAsync(cp.listenerResultsChan, cp.debouncerResultsChan)
	if err != nil {
		cp.limiter.StopSync()
		return
	}

	err = cp.listener.StartAsync(cp.listenerResultsChan)
	if err != nil {
		cp.debouncer.StopSync()
		cp.limiter.StopSync()
	}

	return
}

func (cp *ConnectionProvider) Stop() {
	cp.listener.StopSync()
	cp.debouncer.StopSync()
	cp.limiter.StopSync()
}
