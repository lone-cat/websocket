package acceptor

import (
	"errors"
	"net"
)

type ConnConverter struct {
	idGeneratorFunc func() string
	stopSemaphore   StopSemaphoreI
	logger          LoggerI
}

func NewConnConverter(idGeneratorFunc func() string, stopSemaphore StopSemaphoreI, logger LoggerI) *ConnConverter {
	return &ConnConverter{
		idGeneratorFunc: idGeneratorFunc,
		stopSemaphore:   stopSemaphore,
		logger:          logger,
	}
}

func (cc *ConnConverter) StartAsync(chanFrom <-chan net.Conn, chanTo chan<- AdvancedNetConnI) error {
	if !cc.stopSemaphore.Start() {
		return errors.New(`ConnConverter not stopped yet`)
	}

	cc.logger.Info(`starting ConnConverter...`)

	cc.stopSemaphore.Start()

	go cc.start(chanFrom, chanTo)

	cc.logger.Info(`ConnConverter started`)

	return nil
}

func (cc *ConnConverter) start(chanFrom <-chan net.Conn, chanTo chan<- AdvancedNetConnI) {

	defer func() {
		cc.stopSemaphore.FinishStopping()
		cc.logger.Info(`ConnConverter stopped`)
	}()

	var conn net.Conn
	var advConn AdvancedNetConnI
	var ok bool

loop:
	for {
		select {
		case conn, ok = <-chanFrom:
			if !ok {
				break loop
			}
			advConn = ConvertNetConnToAdvanced(conn, cc.idGeneratorFunc)
			select {
			case chanTo <- advConn:
				// successfully passed
			case <-cc.stopSemaphore.GetStoppingChannel():
				cc.logger.Info(`connection from ` + conn.RemoteAddr().String() + ` dropped. ConnConverter stopping...`)
				cc.closeConn(conn)
			}
		case <-cc.stopSemaphore.GetStoppingChannel():
			break loop
		}
	}
}

func (cc *ConnConverter) closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		cc.logger.Error(err)
	}
}

func (cc *ConnConverter) StopAsync() <-chan struct{} {
	cc.logger.Info(`stopping ConnConverter...`)
	cc.stopSemaphore.StartStopping()
	return cc.stopSemaphore.GetStoppedChannel()
}

func (cc *ConnConverter) StopSync() {
	<-cc.StopAsync()
}
