package listener

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/lone-cat/websocket/mock"
	"github.com/lone-cat/websocket/sem"
)

func TestDebouncer(t *testing.T) {

	const (
		connCount = 4
		timeout   = time.Second / 2
		interval  = timeout / connCount
	)

	deb := NewDebouncer(interval, &sem.TwoStage{}, &mock.Logger{Srv: `test`})
	chFrom := make(chan net.Conn)
	chTo := make(chan net.Conn)

	var c uint8

	go func() {
		for _ = range chTo {
			c++
		}
	}()

	go func() {
		var i int
		timer := time.NewTimer(time.Nanosecond)
		for {
			<-timer.C
			timer.Reset(time.Millisecond * 10)
			i++
			chFrom <- &mock.NetConn{Addr: &mock.NetAddr{Addr: fmt.Sprint(i)}}
		}
	}()

	time.Sleep(time.Millisecond * 10)

	err := deb.StartAsync(chFrom, chTo)

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(timeout)

	deb.StopSync()

	if c != connCount {
		t.Errorf(`taken %d connections instead of limit %d`, c, connCount)
	}

}
