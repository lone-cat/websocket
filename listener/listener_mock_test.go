package listener

import (
	"net"
	"time"
)

type listenerMock struct {
	stop bool
}

func (lm *listenerMock) StartAsync(resultChan chan<- net.Conn) error {
	lm.stop = false
	go func() {
		dur := time.Second * 3
		timer := time.NewTimer(dur)

		defer func() {
			if !timer.Stop() {
				<-timer.C
			}
		}()

		for !lm.stop {
			connMock := &netConnMock{}
			<-timer.C
			timer.Reset(dur)
			resultChan <- connMock
		}
	}()

	return nil
}

func (lm *listenerMock) StopSync() {
	lm.stop = true
}
