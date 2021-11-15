package acceptor

import (
	"errors"

	"github.com/mailru/easygo/netpoll"
)

type Syscaller struct {
	poller        netpoll.Poller
	stopSemaphore StopSemaphoreI
	logger        LoggerI
}

func NewSyscaller(stopSemaphore StopSemaphoreI, logger LoggerI) *Syscaller {
	poller, err := netpoll.New(nil)
	if err != nil {
		panic(err)
	}

	return &Syscaller{
		poller:        poller,
		stopSemaphore: stopSemaphore,
		logger:        logger,
	}
}

func (s *Syscaller) StartAsync(chanFrom <-chan AdvancedNetConnI, chanTo chan<- SyscallConnectionI) error {
	if !s.stopSemaphore.Start() {
		return errors.New(`limiter not stopped yet`)
	}

	s.logger.Info(`starting limiter...`)

	s.stopSemaphore.Start()

	go s.start(chanFrom, chanTo)

	s.logger.Info(`limiter started`)

	return nil
}

func (s *Syscaller) start(chanFrom <-chan AdvancedNetConnI, chanTo chan<- SyscallConnectionI) {

	defer func() {
		s.stopSemaphore.FinishStopping()
		s.logger.Info(`ConnConverter stopped`)
	}()

	var advConn AdvancedNetConnI
	var syscallConn SyscallConnectionI
	var ok bool

loop:
	for {
		select {
		case advConn, ok = <-chanFrom:
			if !ok {
				break loop
			}
			desc, err := netpoll.Handle(advConn, netpoll.EventRead|netpoll.EventEdgeTriggered|netpoll.EventOneShot)
			if err != nil {
				s.logger.Error(err)
				continue loop
			}

			syscallConn = ConvertAdvancedConnToSyscall(advConn, desc, s.poller)
			select {
			case chanTo <- syscallConn:
				// successfully passed
			case <-s.stopSemaphore.GetStoppingChannel():
				s.logger.Info(`connection from ` + advConn.RemoteAddr().String() + ` dropped. ConnConverter stopping...`)
				s.closeConn(advConn)
			}
		case <-s.stopSemaphore.GetStoppingChannel():
			break loop
		}
	}
}

func (s *Syscaller) closeConn(conn AdvancedNetConnI) {
	err := conn.Close()
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Syscaller) StopAsync() <-chan struct{} {
	s.logger.Info(`stopping limiter...`)
	s.stopSemaphore.StartStopping()
	return s.stopSemaphore.GetStoppedChannel()
}

func (s *Syscaller) StopSync() {
	<-s.StopAsync()
}
