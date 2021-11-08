package acceptor

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
