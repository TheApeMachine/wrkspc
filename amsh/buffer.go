package amsh

import (
	"os"

	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/hefner"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Buffer is an instance of the Ape Machine Shell.
*/
type Buffer struct {
	args     []string
	disposer *twoface.Disposer
	tty      *Console
	pipe     hefner.Pipe
}

/*
NewBuffer returns a pointer reference to a constructed Ape Machine Shell.
*/
func NewBuffer(
	args []string,
	disposer *twoface.Disposer,
	pipe hefner.Pipe,
) *Buffer {
	errnie.Traces()
	return &Buffer{
		args:     args,
		disposer: disposer,
		tty:      NewConsole(disposer).Resize(),
		pipe:     pipe,
	}
}

/*
Execute the Buffer and drop into the shell.
*/
func (buffer *Buffer) Execute() chan []byte {
	errnie.Traces()
	out := make(chan []byte)

	go func() {
		defer close(out)

		// Receive characters typed from the keyboard from a goroutine
		// running in the background to serve as the listener.
		for char := range buffer.tty.StdIn() {
			// Enable CTRL+C to exit and get back to `cooked` mode.
			if char[0] == 3 {
				os.Exit(0)
			}

			// Send the character byte out to the channel so a terminal ui
			// component can use it as a data channel.
			out <- char
		}
	}()

	return out
}
