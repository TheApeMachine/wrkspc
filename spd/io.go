package spd

import (
	"io"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Read the contents of the first layer into p.

Make sure p has a pre-allocated length that corresponds to the length
of the data you want to read into it.
*/
func (dg *Datagram) Read(p []byte) (n int, err error) {
	if n, err = dg.ReadAt(p, int64(dg.Ptr())); err != nil {
		errnie.Handles(err)
		return n, err
	}

	return n, err
}

/*
Write the contents of p into a new Layer on the Datagram instance.
*/
func (dg *Datagram) Write(p []byte) (n int, err error) {
	if err = dg.newLayer().Set(int(dg.Ptr()), p); err != nil {
		return 0, err
	}

	return len(p), io.EOF
}

/*
Close ...
*/
func (dg *Datagram) Close() error {
	return nil
}
