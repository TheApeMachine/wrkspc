package spd

import (
	"github.com/wrk-grp/errnie"
)

/*
Read implements the io.Reader interface, and every call should return
the next layer of the payload.
*/
func (dg Datagram) Read(p []byte) (n int, err error) {
	errnie.Trace()

	var b []byte
	if b, err = dg.ReadAt(int64(dg.Ptr())); errnie.IOError(err) {
		errnie.Handles(err)
		return
	}

	return copy(p, b), err
}

/*
Write implements the io.Writer interface, and every call should append
the incoming data to a new layer of the payload.
*/
func (dg *Datagram) Write(p []byte) (n int, err error) {
	errnie.Trace()

	layers := dg.newLayer()

	if err = layers.Set(int(dg.Ptr()), p); errnie.IOError(err) {
		return 0, err
	}

	return len(p), nil
}

/*
Close implements the io.Closer interface, and should be called when the
payload request is done.
*/
func (dg *Datagram) Close() error {
	errnie.Trace()
	return nil
}
