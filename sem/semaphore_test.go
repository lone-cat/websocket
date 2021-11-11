package sem

import (
	"testing"

	"github.com/lone-cat/websocket/listener"
)

func TestImplementation(t *testing.T) {
	var sem listener.StopSemaphoreI
	sem = &TwoStage{}
	sem.FinishStopping()
}
