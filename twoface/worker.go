package twoface

import "io"

/*
Worker wraps a concurrent process that is able to process Job types
scheduled onto a Pool.
*/
type Worker interface {
	io.ReadWriteCloser
}
