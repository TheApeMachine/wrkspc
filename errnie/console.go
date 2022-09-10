package errnie

import "sigs.k8s.io/kind/pkg/log"

type Console struct{}

func NewConsole() Logger {
	return NewLogger(Console{})
}

func (logger Console) Print() Error {
	return NewError(nil)
}

func (logger Console) Error(message string)                      {}
func (logger Console) Errorf(format string, args ...interface{}) {}
func (logger Console) Warn(message string)                       {}
func (logger Console) Warnf(format string, args ...interface{})  {}
func (logger Console) V(level log.Level) log.InfoLogger {
	return InfoLogger{}
}

type InfoLogger struct{}

func (info InfoLogger) Info(message string)                      {}
func (info InfoLogger) Infof(format string, args ...interface{}) {}
func (info InfoLogger) Enabled() bool                            { return true }
