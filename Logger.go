package go_piggy

import (
	"fmt"
	"log"
	"os"
	"testing"
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

type testLogger struct {
	log *testing.T
}

func (t *testLogger) Info(msg ...interface{}) {
	t.log.Log(msg)
}

func NewTestLogger(t *testing.T) Logger {
	return &testLogger{t}
}
