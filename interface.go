package websocket

import (
	"io"
	"net"
	"time"
)

type ConnectionWithIdI interface {
	GetId() string
	ConnectionI
}

type ConnectionI interface {
	net.Conn
}

type ReaderWithDeadlineI interface {
	SetReadDeadline(time.Time) error
	ReaderI
}

type ReaderI interface {
	io.Reader
}

type WriterI interface {
	io.Writer
}

type CloserI interface {
	io.Closer
}
