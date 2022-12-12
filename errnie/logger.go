package errnie

import (
	"bytes"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

type LogLevel uint

const (
	ERROR LogLevel = iota
	WARNING
	INFO
	DEBUG
)

/*
Logger is an interface objects can implement if they want to act as
log output channels for errnie.
*/
type Logger interface {
	Error(...any)
	Warning(...any)
	Info(...any)
	Debug(...any)
	Inspect(...any)
}

func write(level string, msgs ...any) {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(level)

	for _, msg := range msgs {
		buf.WriteString(" ")
		fmt.Fprintf(buf, "%v", msg)
	}

	buf.WriteString("\n")
	ctx.log.Write(buf.Bytes())
}

/*
Warns is syntactic sugar to call the Warning method on
a Logger interface.
*/
func Warns(msgs ...any) {
	write("WARNING", msgs...)
}

/*
Informs is syntactic sugar to call the Info method on
a Logger interface.
*/
func Informs(msgs ...any) {
	write(" INFO  ", msgs...)
}

/*
Debugs is syntactic sugar to call the Debug method on
a Logger interface.
*/
func Debugs(msgs ...any) {
	write(" DEBUG ", msgs...)
}

/*
Inspects is syntactic sugar to dump the structure and values
of objects with arbitrary complexity to logger output channels.
*/
func Inspects(msgs ...any) {
	buf := bytes.NewBuffer([]byte{})

	for _, msg := range msgs {
		spew.Fprintf(buf, "%v", msg)
	}

	buf.WriteString("\n")
	ctx.log.Write(buf.Bytes())
}
