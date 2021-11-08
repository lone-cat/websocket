package listener

import (
	"net"
	"sync"
)

type connWithAdditionalFunctionOnClose struct {
	net.Conn
	mu          sync.Mutex
	closingFunc func()
}

func (cwafoc *connWithAdditionalFunctionOnClose) Close() error {
	cwafoc.mu.Lock()
	defer cwafoc.mu.Unlock()
	if cwafoc.closingFunc != nil {
		cwafoc.closingFunc()
		cwafoc.closingFunc = nil
	}
	return cwafoc.Conn.Close()
}

func addFunctionBeforeClose(conn net.Conn, closingFunc func()) net.Conn {
	return &connWithAdditionalFunctionOnClose{
		Conn:        conn,
		closingFunc: closingFunc,
	}
}
