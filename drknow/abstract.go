package drknow

import "io"

/*
Abstract is an item that holds both data and behavior,
which we can group as `knowledge` and `skill`.
*/
type Abstract interface {
	io.ReadWriteCloser
}
