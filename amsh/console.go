package amsh

import (
	"os"

	"github.com/containerd/console"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Console is a thin wrapper around the terminal itself.

I am Consolio! I need TTY for my bunngghooole...
*/
type Console struct {
	current  console.Console
	disposer *twoface.Disposer
}

/*
NewConsole constructs a Console wrapper and returns a reference to itself.
This will put the terminal in `raw` mode, and it will remain this way until Cleanup is called which
will set it back to `cooked` mode. This is important, otherwise you leave the terminal in a state
where it shows unusual behavior.
*/
func NewConsole(disposer *twoface.Disposer) *Console {
	errnie.Traces()
	current := console.Current()
	errnie.Handles(current.SetRaw()).With(errnie.KILL)

	return &Console{
		current:  current,
		disposer: disposer,
	}
}

/*
Resize the terminal if for some reason the output buffer does not match the current window size.
*/
func (terminal *Console) Resize() *Console {
	errnie.Traces()
	ws, err := terminal.current.Size()
	errnie.Handles(err).With(errnie.NOOP)
	terminal.current.Resize(ws)
	return terminal
}

/*
Cleanup will put the terminal back in `cooked` mode, which is the behavior you are used to.
*/
func (terminal *Console) Cleanup() {
	errnie.Traces()
	terminal.current.Reset()
}

/*
StdIn listens for keyboard events and pipes them back to the Buffer
*/
func (terminal *Console) StdIn() chan []byte {
	errnie.Traces()
	out := make(chan []byte)

	go func() {
		defer close(out)
		defer terminal.Cleanup()

		// We allocate this above the for loop so we only have to do it once.
		keyByte := make([]byte, 1)

		for {
			select {
			case <-terminal.disposer.Done():
				// With the defered statements, returning from the goroutine
				// is all the cleanup we need to do.
				return
			default:
				os.Stdin.Read(keyByte)
				out <- keyByte
			}
		}
	}()

	return out
}
