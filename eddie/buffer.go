package eddie

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/containerd/console"
	"github.com/mattn/go-localereader"
	"github.com/muesli/cancelreader"
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

func (buffer *Buffer) Focus() chan []interface{} {
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

	readDone := make(chan struct{})
	out := make(chan []interface{})
	reader, err := cancelreader.NewReader(buffer.input.(*os.File))
	errnie.Handles(err)

	go func() {
		defer close(readDone)

		for {
			var buf [256]byte
			numBytes, err := reader.Read(buf[:])
			errnie.Handles(err)

			b := buf[:numBytes]
			b, err = localereader.UTF8(b)
			errnie.Handles(err)

			var runeSets [][]rune
			var runes []rune

			for i, w := 0, 0; i < len(b); i += w {
				r, width := utf8.DecodeRune(b[i:])

				if r == utf8.RuneError {
					errnie.Handles(errors.New("could not decode rune"))
					return
				}

				if r == '\x1b' && len(runes) > 1 {
					runeSets = append(runeSets, runes)
					runes = []rune{}
				}

				runes = append(runes, r)
				w = width
			}

			runeSets = append(runeSets, runes)

			if len(runeSets) == 0 {
				errnie.Handles(errors.New("received 0 runes from input"))
				return
			}

			var msgs []interface{}
			for _, runes := range runeSets {
				// Is it a sequence, like an arrow key?
				if k, ok := sequences[string(runes)]; ok {
					msgs = append(msgs, KeyMsg(k))
					continue
				}

				// Some of these need special handling.
				hex := fmt.Sprintf("%x", runes)
				if k, ok := hexes[hex]; ok {
					msgs = append(msgs, KeyMsg(k))
					continue
				}

				// Is the alt key pressed? If so, the buffer will be prefixed with an
				// escape.
				alt := false
				if len(runes) > 1 && runes[0] == 0x1b {
					alt = true
					runes = runes[1:]
				}

				for _, v := range runes {
					// Is the first rune a control character?
					r := KeyType(v)
					if r <= keyUS || r == keyDEL {
						msgs = append(msgs, KeyMsg(Key{Type: r, Alt: alt}))
						continue
					}

					// If it's a space, override the type with KeySpace (but still include
					// the rune).
					if r == ' ' {
						msgs = append(msgs, KeyMsg(Key{Type: KeySpace, Runes: []rune{v}, Alt: alt}))
						continue
					}

					// Welp, just regular, ol' runes.
					msgs = append(msgs, KeyMsg(Key{Type: KeyRunes, Runes: []rune{v}, Alt: alt}))
				}
			}

			out <- msgs
		}
	}()

	return out
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
