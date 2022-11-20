package errnie

import (
	"log"
)

type ConsoleLogger struct {
	Logger
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (logger *ConsoleLogger) Error(msgs ...any) {
	log.Println(msgs...)
}

func (logger *ConsoleLogger) Warning(msgs ...any) {
	log.Println(msgs...)
}

func (logger *ConsoleLogger) Info(msgs ...any) {
	log.Println(msgs...)
}

func (logger *ConsoleLogger) Debug(msgs ...any) {
	log.Println(msgs...)
}
