package datura

import "io"

type Store interface {
	io.ReadWriter
}
