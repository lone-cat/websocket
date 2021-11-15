package acceptor

import (
	"errors"
	"io"
	"net"
	"time"
)

var (
	SetDeadLineError = errors.New(`deadline`)
	TimeoutError     = errors.New(`timeout`)
	EOF              = errors.New(`eof`)
	UnknownError     = errors.New(`unknown`)
)

type AdvancedNetConn struct {
	id      string
	timeout time.Duration
	net.Conn
}

func (anc *AdvancedNetConn) Id() string {
	return anc.id
}

func (anc *AdvancedNetConn) NonBlockRead(buf *[]byte) (result []byte, err error) {
	err = anc.SetReadDeadline(time.Now().Add(anc.timeout))
	if err != nil {
		err = SetDeadLineError
		return
	}
	result, err = anc.read(buf)
	if err != nil {
		err = ConvertNetError(err)
	}

	return
}

func (anc *AdvancedNetConn) read(buf *[]byte) (result []byte, err error) {
	var n int
	n, err = anc.Read(*buf)
	if n > 0 {
		result = make([]byte, n)
		copy(result, *buf)
	}
	return
}

func ConvertNetConnToAdvanced(conn net.Conn, idGeneratorFunc func() string) *AdvancedNetConn {
	return &AdvancedNetConn{
		id:      idGeneratorFunc(),
		timeout: time.Millisecond,
		Conn:    conn,
	}
}

func ConvertNetError(err error) error {
	if err == nil {
		return nil
	}

	if err == io.EOF {
		return EOF
	}

	switch errT := err.(type) {
	case net.Error:
		if errT.Timeout() {
			return TimeoutError
		}
	default:
	}

	return UnknownError
}
