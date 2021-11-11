package mock

import (
	"testing"

	"github.com/lone-cat/websocket/listener"
)

func TestListenerMock(t *testing.T) {
	var l listener.ConnectionProviderI
	l = &Listener{}

	l.StopSync()
}
