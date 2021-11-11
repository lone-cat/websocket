package listener

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/lone-cat/websocket/mock"
	"github.com/lone-cat/websocket/sem"
)

func TestListener(t *testing.T) {
	const port = 8085
	l := NewListener(true, port, &sem.TwoStage{}, &mock.Logger{Srv: `test`})

	connChan := make(chan net.Conn)
	err := l.StartAsync(connChan)
	defer l.StopSync()
	if err != nil {
		t.Fatal(err)
	}

	var c uint8

	var conn net.Conn

	go func() {
		for conn := range connChan {
			c++
			t.Logf(`conn number %d (%s) accepted`, c, conn.RemoteAddr().String())
		}
	}()

	time.Sleep(time.Millisecond * 10)
	conn, err = net.Dial(`tcp`, `:`+strconv.Itoa(port))
	time.Sleep(time.Millisecond * 10)

	if err != nil {
		t.Fatal(err)
	}

	if c < 1 {
		t.Error(`connection 1 was not accepted`)
	}

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 10)
	conn, err = net.Dial(`tcp`, `:`+strconv.Itoa(port))
	time.Sleep(time.Millisecond * 10)

	if err != nil {
		t.Fatal(err)
	}

	if c < 2 {
		t.Error(`connection 2 was not accepted`)
	}

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}
