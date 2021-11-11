package listener

import (
	"testing"
	"time"

	"github.com/lone-cat/websocket/mock"
	"github.com/lone-cat/websocket/sem"
)

func TestConnectionProvider(t *testing.T) {
	const (
		connCount = 5
		timeout   = time.Millisecond * 500
		interval  = timeout / connCount
	)

	l := &mock.Listener{Interval: interval}
	m := []ConnectionMiddlewareI{&mock.ConnMiddleware{}, &mock.ConnMiddleware{}}
	testConnectionProvider(timeout-interval/2, l, connCount, m, t, 1)
	time.Sleep(interval)

	m = []ConnectionMiddlewareI{}
	testConnectionProvider(timeout-interval/2, l, connCount, m, t, 2)
	time.Sleep(interval)

	l = &mock.Listener{Interval: time.Millisecond * 10}
	m = []ConnectionMiddlewareI{NewDebouncer(interval, &sem.TwoStage{}, &mock.Logger{Srv: `test`})}
	testConnectionProvider(timeout-interval/2, l, connCount, m, t, 3)
	time.Sleep(interval)

	//l = &mock.Listener{Interval: time.Millisecond * 10}
	m = []ConnectionMiddlewareI{
		NewDebouncer(interval/2, &sem.TwoStage{}, &mock.Logger{Srv: `test`}),
		NewLimiter(connCount, &sem.TwoStage{}, &mock.Logger{Srv: `test`}),
	}
	testConnectionProvider(timeout-interval/2, l, connCount, m, t, 4)
}

func testConnectionProvider(timeout time.Duration, l ConnectionProviderI, connCount int, middlewares []ConnectionMiddlewareI, t *testing.T, iter int) {
	cp := NewConnectionProvider(l, middlewares...)

	resultChan := cp.GetResultChan()
	c := 0
	go func() {
		for range resultChan {
			//fmt.Printf("iter %d: connection accepted %s \n", iter, conn.RemoteAddr().String())
			c++
		}
	}()

	time.Sleep(time.Millisecond * 10)

	cp.Start()

	time.Sleep(timeout)

	cp.Stop()

	if c != connCount {
		t.Errorf("iter %d: taken %d connections instead of %d", iter, c, connCount)
	}
}
