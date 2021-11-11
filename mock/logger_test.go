package mock

import (
	"testing"

	"github.com/lone-cat/websocket/listener"
)

func TestLoggerMock(t *testing.T) {
	var m listener.LoggerI
	m = &Logger{`test`}

	m.Info(`wow!`)
}
