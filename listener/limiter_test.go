package listener

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/lone-cat/websocket/mock"
	"github.com/lone-cat/websocket/sem"
)

func TestLimiter(t *testing.T) {
	const interval = time.Millisecond * 10
	const connCount = 5

	lim := NewLimiter(connCount, &sem.TwoStage{}, &mock.Logger{Srv: `test`})
	chFrom := make(chan net.Conn)
	chTo := make(chan net.Conn)

	var c int
	var lastConn net.Conn

	go func() {
		for conn := range chTo {
			lastConn = conn
			c++
			t.Logf(`conn number %d (%s) accepted`, c, lastConn.RemoteAddr().String())
		}
	}()

	go func() {
		for i := 0; i < connCount+2; i++ {
			chFrom <- &mock.NetConn{Addr: &mock.NetAddr{Addr: fmt.Sprint(i)}}
			time.Sleep(interval)
		}
	}()

	err := lim.StartAsync(chFrom, chTo)
	defer lim.StopSync()

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(interval * (connCount + 1))

	c2 := c

	if c != connCount {
		t.Errorf(`taken %d connections instead of limit %d`, c, connCount)
	}

	lastConn.Close()

	chFrom <- &mock.NetConn{Addr: &mock.NetAddr{Addr: fmt.Sprint(22)}}

	time.Sleep(interval)

	if c != c2+1 {
		t.Errorf(`dropped connection does not free place %d`, c)
	}

}
