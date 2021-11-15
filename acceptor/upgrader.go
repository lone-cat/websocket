package acceptor

import (
	"errors"
	"net"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
)

/*
var Header ws.HandshakeHeaderHTTP = ws.HandshakeHeaderHTTP(
	http.Header{
		"X-Go-Version": []string{runtime.Version()},
	})
*/

var upgrader = &ws.Upgrader{
	OnHost: func(host []byte) error {
		return nil
		/*fmt.Println(string(host))
		  endpoint := `localhost`
		  if string(host) == endpoint {
		  	return nil
		  }
		  return ws.RejectConnectionError(
		  	ws.RejectionStatus(403),
		  	ws.RejectionHeader(ws.HandshakeHeaderString(`X-Want-Host: `+endpoint+"\r\n")),
		  )
		*/
	},
	OnHeader: func(key, value []byte) error {
		if string(key) != "Cookie" {
			return nil
		}
		ok := httphead.ScanCookie(value, func(key, value []byte) bool {
			// Check session here or do some other stuff with cookies.
			// Maybe copy some values for future use.
			return true
		})
		if ok {
			return nil
		}
		return ws.RejectConnectionError(
			ws.RejectionReason("bad cookie"),
			ws.RejectionStatus(400),
		)
	},
	OnBeforeUpgrade: func() (ws.HandshakeHeader, error) {
		// first response can be header
		return nil, nil
	},
}

type Upgrader struct {
	stopSemaphore StopSemaphoreI
	logger        LoggerI
}

func NewUpgrader(stopSemaphore StopSemaphoreI, logger LoggerI) *Upgrader {
	return &Upgrader{
		stopSemaphore: stopSemaphore,
		logger:        logger,
	}
}

func (u *Upgrader) StartAsync(chanFrom <-chan net.Conn, chanTo chan<- net.Conn) error {
	if !u.stopSemaphore.Start() {
		return errors.New(`limiter not stopped yet`)
	}

	u.logger.Info(`starting limiter...`)

	u.stopSemaphore.Start()

	go u.start(chanFrom, chanTo)

	u.logger.Info(`limiter started`)

	return nil
}

func (u *Upgrader) start(chanFrom <-chan net.Conn, chanTo chan<- net.Conn) {
	var conn net.Conn
	var ok bool
loop:
	for {
		select {
		case conn, ok = <-chanFrom:
			if !ok {
				break loop
			}
			_, err := upgrader.Upgrade(conn)
			if err != nil {
				u.logger.Info(`error upgrading connection from ` + conn.RemoteAddr().String() + `. dropping...`)
				u.closeConn(conn)
				continue loop
			}
			select {
			case chanTo <- conn:
				// successfully upgraded
			case <-u.stopSemaphore.GetStoppingChannel():
				u.logger.Info(`connection from ` + conn.RemoteAddr().String() + ` dropped. service stopping...`)
				u.closeConn(conn)
				break loop
			}
		case <-u.stopSemaphore.GetStoppingChannel():
			break loop
		}
	}
}

func (u *Upgrader) closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		u.logger.Error(err)
	}
}

func (u *Upgrader) StopAsync() <-chan struct{} {
	u.logger.Info(`stopping limiter...`)
	u.stopSemaphore.StartStopping()
	return u.stopSemaphore.GetStoppedChannel()
}

func (u *Upgrader) StopSync() {
	<-u.StopAsync()
}
