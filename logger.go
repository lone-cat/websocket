package websocket

import "fmt"

type logger struct {
	srv string
}

func (l *logger) Error(args ...interface{}) {
	fmt.Print(`[` + l.srv + `]`)
	fmt.Println(args...)
}

func (l *logger) Info(args ...interface{}) {
	fmt.Print(`[` + l.srv + `]`)
	fmt.Println(args...)
}
