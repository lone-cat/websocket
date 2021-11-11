package mock

type Logger struct {
	Srv string
}

func (l *Logger) Error(args ...interface{}) {
	/*fmt.Print(`[` + l.Srv + `] `)
	fmt.Println(args...)*/
}

func (l *Logger) Info(args ...interface{}) {
	/*fmt.Print(`[` + l.Srv + `] `)
	fmt.Println(args...)*/
}
