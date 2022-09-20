package eddie

import "github.com/muesli/termenv"

type Buffer struct {
	altscreen bool
}

func NewBuffer() *Buffer {
	return &Buffer{
		altscreen: true,
	}
}

func (buffer *Buffer) Init() *Buffer {
	if buffer.altscreen {
		termenv.AltScreen()
		termenv.ClearScreen()
		termenv.MoveCursor(0, 0)
	}

	return buffer
}
