package listener

import "testing"

func TestConnectionProvider(t *testing.T) {
	l := &listenerMock{}
	m1 := &connMiddlewareMock{}
	m2 := &connMiddlewareMock{}
	cp := NewConnectionProvider(l, m1, m2)

	cp.Start()

	cp.Stop()
}
