package mock

import (
	"testing"

	"github.com/lone-cat/websocket/listener"
)

func TestMiddlewareMock(t *testing.T) {
	var m listener.ConnectionMiddlewareI
	m = &ConnMiddleware{}

	m.StopSync()
}
