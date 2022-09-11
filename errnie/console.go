package errnie

import (
	"fmt"
	"time"

	"github.com/theapemachine/wrkspc/tui"
	"sigs.k8s.io/kind/pkg/log"
)

type Console struct{}

func NewConsole() Logger {
	return NewLogger(Console{})
}

func (logger Console) Print(l, t, c, i string) {
	fmt.Println(
		tui.NewLabel(t).Print(),
		tui.NewColor(
			"MUTE", time.Now().Format("2006-01-02 15:04:05.000000"),
		).Print(),
		tui.NewIcon(i),
		tui.NewColor(c, l).Print(),
	)
}

func (logger Console) Error(message string) {
	Logs(message).With(ERROR)
}

func (logger Console) Errorf(format string, args ...interface{}) {
	Logs(fmt.Sprintf(format, args...)).With(ERROR)
}

func (logger Console) Warn(message string) {
	Logs(message).With(WARNING)
}

func (logger Console) Warnf(format string, args ...interface{}) {
	Logs(fmt.Sprintf(format, args...)).With(WARNING)
}

func (logger Console) V(level log.Level) log.InfoLogger {
	return InfoLogger{}
}

type InfoLogger struct{}

func (info InfoLogger) Info(message string) {
	Logs(message).With(INFO)
}

func (info InfoLogger) Infof(format string, args ...interface{}) {
	Logs(fmt.Sprintf(format, args...)).With(INFO)
}

func (info InfoLogger) Enabled() bool {
	return true
}
