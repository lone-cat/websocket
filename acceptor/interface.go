package acceptor

import "net"

type LoggerI interface {
	Info(...interface{})
	Error(...interface{})
}

type StopSemaphoreI interface {
	IsStopping() bool
	IsStopped() bool
	Start() bool
	StartStopping()
	FinishStopping()
	WaitTillStopped()
	GetStoppingChannel() chan struct{}
	GetStoppedChannel() chan struct{}
}

type NonBlockingReaderI interface {
	NonBlockRead(buf *[]byte) (result []byte, err error)
}

type IdentificatedI interface {
	Id() string
}

type SyscalledI interface {
	Init(func()) error
	Resume() error
}

type AdvancedNetConnI interface {
	IdentificatedI
	NonBlockingReaderI
	net.Conn
}

type SyscallConnectionI interface {
	AdvancedNetConnI
	SyscalledI
}
