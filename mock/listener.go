package mock

import (
	"net"
	"runtime"
	"strconv"
	"time"
)

type Listener struct {
	Stop     bool
	Interval time.Duration
}

func (lm *Listener) StartAsync(resultChan chan<- net.Conn) error {
	lm.Stop = false
	if lm.Interval < 1 {
		lm.Interval = time.Millisecond * 100
	}
	go func() {
		timer := time.NewTimer(time.Nanosecond)

		defer func() {
			if !timer.Stop() {
				<-timer.C
			}
		}()

		var i = 0
		for !lm.Stop {
			i++
			connMock := &NetConn{&NetAddr{strconv.Itoa(i)}}
			<-timer.C
			timer.Reset(lm.Interval)
			select {
			case resultChan <- connMock:

			default:
				runtime.Gosched()
			}
		}
	}()

	return nil
}

func (lm *Listener) StopSync() {
	lm.Stop = true
	time.Sleep(time.Millisecond)
}
