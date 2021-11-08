package listener

import (
	"errors"
	"net"
)

type Limiter struct {
	connectionSlots chan struct{}
	stopSemaphore   StopSemaphoreI
	logger          LoggerI
}

func NewLimiter(connectionsLimit uint8, stopSemaphore StopSemaphoreI, logger LoggerI) *Limiter {
	return &Limiter{
		connectionSlots: make(chan struct{}, connectionsLimit),
		stopSemaphore:   stopSemaphore,
		logger:          logger,
	}
}

func (l *Limiter) StartAsync(chanFrom <-chan net.Conn, chanTo chan<- net.Conn) error {
	if !l.stopSemaphore.Start() {
		return errors.New(`limiter not stopped yet`)
	}

	l.logger.Info(`starting limiter...`)

	l.stopSemaphore.Start()

	go l.limit(chanFrom, chanTo)

	l.logger.Info(`limiter started`)

	return nil
}

func (l *Limiter) limit(chanFrom <-chan net.Conn, chanTo chan<- net.Conn) {

	defer func() {
		l.stopSemaphore.FinishStopping()
		l.logger.Info(`limiter stopped`)
	}()

	var conn net.Conn
	var ok bool

loop:
	for {
		select {
		case conn, ok = <-chanFrom:
			if !ok {
				break loop
			}
			select {
			case l.connectionSlots <- struct{}{}:
				conn = addFunctionBeforeClose(conn, l.releaseConnectionSlot)
				select {
				case chanTo <- conn:
					// successfully passed
				default:
					l.logger.Info(`connection from ` + conn.RemoteAddr().String() + ` was not taken. dropping...`)
					l.closeConn(conn)
				}
			default:
				l.logger.Info(`connection from ` + conn.RemoteAddr().String() + ` exceeds the limit. dropping...`)
				l.closeConn(conn)
			}
		case <-l.stopSemaphore.GetStoppingChannel():
			break loop
		}
	}
}

func (l *Limiter) releaseConnectionSlot() {
	select {
	case <-l.connectionSlots:
	default:
		l.logger.Error(`could not release slot!`)
	}
}

func (l *Limiter) closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		l.logger.Error(err)
	}
}

func (l *Limiter) StopAsync() <-chan struct{} {
	l.logger.Info(`stopping limiter...`)
	l.stopSemaphore.StartStopping()
	return l.stopSemaphore.GetStoppedChannel()
}

func (l *Limiter) StopSync() {
	ch := l.StopAsync()
	<-ch
}
