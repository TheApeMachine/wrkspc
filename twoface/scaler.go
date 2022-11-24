package twoface

import "io"

/*
Scaler is a process that evaluates the resource load of a Pool and
adds or removes Workers according to its opinion about how to divide
the machine resources available.
*/
type Scaler interface {
	io.ReadWriteCloser
}
