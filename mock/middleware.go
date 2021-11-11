package mock

import (
	"errors"
	"net"
	"time"
)

// Debouncer and Limiter mock
type ConnMiddleware struct {
	stop chan struct{}
}

func (cmm *ConnMiddleware) init() {
	if cmm.stop == nil {
		cmm.stop = make(chan struct{})
		close(cmm.stop)
	}
}

func (cmm *ConnMiddleware) StartAsync(chanFrom <-chan net.Conn, chanTo chan<- net.Conn) error {
	cmm.init()
	select {
	case _, _ = <-cmm.stop:
		cmm.stop = make(chan struct{})
	default:
		return errors.New(`can't start`)
	}

	go func() {
		var conn net.Conn

		defer func() {
			//fmt.Println(`stopped middleware`)
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

func (cmm *ConnMiddleware) StopAsync() <-chan struct{} {
	cmm.init()
	select {
	case _, _ = <-cmm.stop:
	default:
		close(cmm.stop)
	}
	return cmm.stop
}

func (cmm *ConnMiddleware) StopSync() {
	<-cmm.StopAsync()
	time.Sleep(time.Millisecond)
}
