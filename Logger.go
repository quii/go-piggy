package go_piggy

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Info(...interface{})
}

type stdoutLogger struct {
	log *log.Logger
}

func (s *stdoutLogger) Info(msg ...interface{}) {
	s.log.Println(fmt.Sprintf("INFO %s", msg))
}

func NewStdoutLogger() Logger {
	log := log.New(os.Stdout, "go-piggy ", log.Ldate|log.Ltime)
	return &stdoutLogger{log}
}
