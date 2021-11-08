package websocket

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/lone-cat/websocket/listener"
)

func TestWebsocket(t *testing.T) {
	lis := listener.Factory{}.CreateConnectionProvider(
		false,
		8080,
		&TwoStageSemaphore{},
		&logger{`listener`},
		time.Second,
		&TwoStageSemaphore{},
		&logger{`debouncer`},
		128,
		&TwoStageSemaphore{},
		&logger{`limiter`},
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
