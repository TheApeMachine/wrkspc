package eddie

import (
	"fmt"
	"io"

	"github.com/muesli/termenv"
)

type Cursor struct {
	w io.Writer
}

func NewCursor(w io.Writer) *Cursor {
	return &Cursor{
		w: w,
	}
}

func (cursor *Cursor) Show() *Cursor {
	fmt.Fprintf(cursor.w, termenv.CSI+termenv.ShowCursorSeq)
	return cursor
}

func (cursor *Cursor) Hide() *Cursor {
	fmt.Fprintf(cursor.w, termenv.CSI+termenv.HideCursorSeq)
	return cursor
}
