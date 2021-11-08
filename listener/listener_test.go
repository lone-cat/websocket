package listener

import (
	"testing"
)

type loggerMock struct{}

func (lm *loggerMock) Info(...interface{}) {

}

func (lm *loggerMock) Error(...interface{}) {

}

type semMock struct{}

func (lm *semMock) IsStopping() bool                  { return true }
func (lm *semMock) IsStopped() bool                   { return true }
func (lm *semMock) Start() bool                       { return true }
func (lm *semMock) StartStopping()                    {}
func (lm *semMock) FinishStopping()                   {}
func (lm *semMock) WaitTillStopped()                  {}
func (lm *semMock) GetStoppingChannel() chan struct{} { return nil }
func (lm *semMock) GetStoppedChannel() chan struct{}  { return nil }

func TestNewListener(t *testing.T) {
	var l interface{}
	l = NewListener(true, 123, &semMock{}, &loggerMock{})
	switch l.(type) {
	case *Listener:

	default:
		t.Error(`NewListener fails creation!`)
	}
}
