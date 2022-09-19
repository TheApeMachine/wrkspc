package errnie

/*
Op enumerates the operations errnie can perform on a error event.
*/
type Op func() (string, string, string)

var (
	// NOOP does not do anything.
	NOOP Op = writeLog(" NOOP  ", "HIGH", "ghost")
	// KILL exits the program with code 1.
	KILL Op = writeLog(" KILL  ", "HIGH", "skull")
	// SUCCESS ...
	SUCCESS Op = writeLog("SUCCESS", "HIGH", "thumu")
	// INFO ...
	INFO Op = writeLog(" INFO  ", "NORM", "badge")
	// DEBUG ...
	DEBUG Op = writeLog(" DEBUG ", "MUTE", "lbug")
	// WARNING ...
	WARNING Op = writeLog("WARNING", "NORM", "warn")
	// ERROR ...
	ERROR Op = writeLog(" ERROR ", "HIGH", "fire")
)

func writeLog(t, c, i string) func() (string, string, string) {
	return func() (string, string, string) {
		return t, c, i
	}
}

func Handles(err error) Error {
	t, c, i := ERROR()
	e := NewError(err)

	if l, ok := ambctx.loggers[0], e.Type != NIL; ok {
		l.Print(e.Error(), t, c, i)
	}

	return e
}
