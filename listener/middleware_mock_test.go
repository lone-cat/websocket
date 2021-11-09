package listener

import (
	"errors"
	"fmt"
	"net"
	"time"
)

// Debouncer and Limiter mock
type connMiddlewareMock struct {
	stop chan struct{}
}

func (cmm *connMiddlewareMock) StartAsync(chanFrom <-chan net.Conn, chanTo chan<- net.Conn) error {
	if cmm.stop == nil {
		cmm.stop = make(chan struct{})
		close(cmm.stop)
	}
	select {
	case _, _ = <-cmm.stop:
		cmm.stop = make(chan struct{})
	default:
		return errors.New(`can't start`)
	}

	go func() {
		var conn net.Conn

		defer func() {
			fmt.Println(`stopped middleware`)
		}()

	loop:
		for {
			select {
			case <-cmm.stop:
				break loop
			case conn = <-chanFrom:
				select {
				case chanTo <- conn:

				case <-cmm.stop:
					break loop
				}
			}

		}
	}()

	return nil
}

func (cmm *connMiddlewareMock) StopAsync() <-chan struct{} {
	close(cmm.stop)
	return cmm.stop
}

func (cmm *connMiddlewareMock) StopSync() {
	<-cmm.StopAsync()
	time.Sleep(time.Second)
}
