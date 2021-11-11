package websocket

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/lone-cat/websocket/listener"
	"github.com/lone-cat/websocket/mock"
	"github.com/lone-cat/websocket/sem"
)

func aTestWebsocket(t *testing.T) {
	lis := listener.Factory{}.CreateConnectionProvider(
		false,
		8080,
		&sem.TwoStage{},
		&mock.Logger{Srv: `listener`},
		time.Second,
		&sem.TwoStage{},
		&mock.Logger{Srv: `debouncer`},
		128,
		&sem.TwoStage{},
		&mock.Logger{Srv: `limiter`},
	)

	err := lis.Start()
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGSEGV,
		syscall.SIGQUIT)

	a := <-sigChan
	log.Println(`Stopping server...`, a)
	lis.Stop()
}
