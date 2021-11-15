package listener

import (
	"errors"
	"net"
)

type Listener struct {
	addr          *net.TCPAddr
	listener      *net.TCPListener
	stopSemaphore StopSemaphoreI
	logger        LoggerI
}

func NewListener(
	local bool,
	port uint16,
	semaphore StopSemaphoreI,
	logger LoggerI,
) *Listener {
	var ip4 net.IP = nil
	if local {
		ip4 = net.IPv4(127, 0, 0, 1)
	}
	addr := net.TCPAddr{IP: ip4, Port: int(port)}

	l := &Listener{
		addr:          &addr,
		logger:        logger,
		stopSemaphore: semaphore,
	}

	return l
}

func (l *Listener) StartAsync(resultChan chan<- net.Conn) error {
	if !l.stopSemaphore.Start() {
		return errors.New(`listener not stopped yet`)
	}

	l.logger.Info(`starting listener...`)

	var err error
	l.listener, err = net.ListenTCP(`tcp`, l.addr)
	if err != nil {
		return err
	}

	l.logger.Info(`listener started on ` + l.addr.String())

	go l.startListening(resultChan)

	return nil
}

func (l *Listener) startListening(resultChan chan<- net.Conn) {
	var err error
	var conn *net.TCPConn

	defer func() {
		l.stopSemaphore.FinishStopping()
		l.logger.Info(`listener stopped`)
	}()

	for !l.stopSemaphore.IsStopping() {
		conn, err = l.listener.AcceptTCP()
		if l.stopSemaphore.IsStopping() {
			break
		}
		if err != nil || conn == nil {
			l.logger.Error(err)
			continue
		}
		select {
		case resultChan <- conn:
			l.logger.Info(`incoming connection from ` + conn.RemoteAddr().String())
		default:
			l.logger.Info(`connection from ` + conn.RemoteAddr().String() + ` was not taken. dropping...`)
			err = conn.Close()
			if err != nil {
				l.logger.Error(err)
			}
		}
	}
}

func (l *Listener) StopAsync() <-chan struct{} {
	l.logger.Info(`stopping listener...`)
	l.stopSemaphore.StartStopping()
	_ = l.listener.Close()
	return l.stopSemaphore.GetStoppedChannel()
}

func (l *Listener) StopSync() {
	<-l.StopAsync()
}
