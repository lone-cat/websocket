package listener

import (
	"errors"
	"net"
	"time"
)

type Debouncer struct {
	duration      time.Duration
	stopSemaphore StopSemaphoreI
	logger        LoggerI
}

func NewDebouncer(duration time.Duration, stopSemaphore StopSemaphoreI, logger LoggerI) *Debouncer {
	d := &Debouncer{
		duration:      duration,
		stopSemaphore: stopSemaphore,
		logger:        logger,
	}

	return d
}

func (d *Debouncer) StartAsync(chanFrom <-chan net.Conn, chanTo chan<- net.Conn) error {
	if !d.stopSemaphore.Start() {
		return errors.New(`debouncer not stopped yet`)
	}

	d.logger.Info(`starting debouncer...`)

	d.stopSemaphore.Start()

	go d.start(chanFrom, chanTo)

	d.logger.Info(`debouncer started`)

	return nil
}

func (d *Debouncer) start(chanFrom <-chan net.Conn, chanTo chan<- net.Conn) {
	timer := time.NewTimer(time.Nanosecond)

	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
		d.stopSemaphore.FinishStopping()
		d.logger.Info(`debouncer stopped`)
	}()

	var conn net.Conn
	var ok bool

loop:
	for {
		// wait for incoming connection to debounce drop or stop signal or debounce period expire
		select {
		case conn, ok = <-chanFrom:
			if !ok {
				break loop
			}
			d.logger.Info(`connection from ` + conn.RemoteAddr().String() + ` came before debounce period expired. dropping...`)
			d.closeConn(conn)
		case <-d.stopSemaphore.GetStoppingChannel():
			break loop
		case <-timer.C:
			timer.Reset(d.duration)
			// wait for incoming connection or exit signal
			select {
			case conn, ok = <-chanFrom:
				if !ok {
					break loop
				}
				// pass connection to output channel or drop
				select {
				case chanTo <- conn:
				default:
					d.logger.Info(`connection from ` + conn.RemoteAddr().String() + ` was not taken. dropping...`)
					d.closeConn(conn)
				}
			case <-d.stopSemaphore.GetStoppingChannel():
				break loop
			}
		}
	}
}

func (d *Debouncer) closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		d.logger.Error(err)
	}
}

func (d *Debouncer) StopAsync() <-chan struct{} {
	d.logger.Info(`stopping debouncer...`)
	d.stopSemaphore.StartStopping()
	return d.stopSemaphore.GetStoppedChannel()
}

func (d *Debouncer) StopSync() {
	<-d.StopAsync()
}
