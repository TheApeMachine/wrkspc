package eddie

import (
	"bytes"
	"io"
	"os"

	"github.com/containerd/console"
	"github.com/muesli/termenv"
	"github.com/theapemachine/wrkspc/errnie"
)

type Buffer struct {
	input     interface{}
	altscreen bool
	back      io.ReadWriter
	console   console.Console
	cursor    *Cursor
}

func NewBuffer(input *os.File) *Buffer {
	return &Buffer{
		input:     input,
		altscreen: true,
		back:      bytes.NewBuffer([]byte{}),
		cursor:    NewCursor(input),
	}
}

func (buffer *Buffer) Init() *Buffer {
	// If no file was given, open a tty for input.
	if buffer.input == nil {
		fh, err := os.Open("/dev/tty")
		errnie.Handles(err)
		buffer.input = fh
	}

	if fh, ok := buffer.input.(*os.File); ok {
		c, err := console.ConsoleFromFile(fh)
		errnie.Handles(err)
		buffer.console = c
	}

	if buffer.altscreen {
		termenv.AltScreen()
		termenv.ClearScreen()
		termenv.MoveCursor(0, 0)
	}

	return buffer
}

func (buffer *Buffer) Focus() {
	if buffer.console != nil {
		// Set the terminal raw mode and hide the cursor.
		errnie.Handles(buffer.console.SetRaw())
		buffer.cursor.Hide()

		// Make sure to restore the terminal so we don't
		// leave the user with something unusable.
		defer func() {
			buffer.cursor.Show()
			errnie.Handles(buffer.console.Reset())
		}()
	}

	<-make(chan struct{})
}

func (buffer *Buffer) Read(p []byte) (n int, err error) {
	var nb int64
	nb, err = io.Copy(bytes.NewBuffer(p), buffer.back)
	return int(nb), err
}

func (buffer *Buffer) Write(p []byte) (n int, err error) {
	var nb int64
	nb, err = io.Copy(buffer.back, bytes.NewBuffer(p))
	return int(nb), err
}
