package errnie

import "sigs.k8s.io/kind/pkg/log"

type Logger interface {
	Print() Error
	Error(string)
	Errorf(string, ...interface{})
	Warn(string)
	Warnf(string, ...interface{})
	V(log.Level) log.InfoLogger
}

type Log struct {
	Value string
}

func NewLog(val string) Log {
	return Log{
		Value: val,
	}
}

/*
Logs is a conveniance method to send values into the logging pipeline.
*/
func Logs(vals ...interface{}) AmbientContext {
	ambctx.logs = append(ambctx.logs, vals...)
	return ambctx
}

/*
NewLogger converts a struct type into a Logger interface type.
*/
func NewLogger(loggerType Logger) Logger {
	return loggerType
}
