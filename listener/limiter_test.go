package listener

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/lone-cat/websocket"
)

func TestLimiter(t *testing.T) {
	var connCount uint8 = 10
	lim := NewLimiter(connCount, &websocket.TwoStageSemaphore{}, &websocket.LoggerMock{Srv: `test`})
	chFrom := make(chan net.Conn)
	chTo := make(chan net.Conn)

	var c uint8
	var lastConn net.Conn

	go func() {
		for conn := range chTo {
			lastConn = conn
			c++
			t.Logf(`conn number %d (%s) accepted`, c, lastConn.RemoteAddr().String())
		}
	}()

	go func() {
		var i uint8
		for i = 0; i < connCount+2; i++ {
			chFrom <- &netConnMock{netAddrMock{fmt.Sprint(i)}}
			time.Sleep(time.Second / 20)
		}
	}()

	err := lim.StartAsync(chFrom, chTo)

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 1)

	c2 := c

	if c != connCount {
		t.Errorf(`taken %d connections instead of limit %d`, c, connCount)
	}

	lastConn.Close()

	chFrom <- &netConnMock{netAddrMock{fmt.Sprint(22)}}

	time.Sleep(time.Second * 1)

	if c != c2+1 {
		t.Errorf(`dropped connection does not free place %d`, c)
	}

}
