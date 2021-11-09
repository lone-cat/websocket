package websocket

import "fmt"

type LoggerMock struct {
	Srv string
}

func (l *LoggerMock) Error(args ...interface{}) {
	fmt.Print(`[` + l.Srv + `] `)
	fmt.Println(args...)
}

func (l *LoggerMock) Info(args ...interface{}) {
	fmt.Print(`[` + l.Srv + `] `)
	fmt.Println(args...)
}
