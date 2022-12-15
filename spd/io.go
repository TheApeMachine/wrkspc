package spd

import (
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Read data layer by layer into p sequentially..

Each time the Read method is called it will return the
next layer down the stack.
*/
func (dg *Datagram) Read(p []byte) (n int, err error) {
	errnie.Trace()

	errnie.Debugs("spd.Datagram.Read Ptr ->", dg.Ptr())

	var b []byte
	if b, err = dg.ReadAt(int64(dg.Ptr())); errnie.IOError(err) {
		errnie.Handles(err)
		return
	}

	p = append(p, b...)
	errnie.Debugs("spd.Datagram.Read ->", string(p))
	return n, err
}

/*
Write the contents of p into a new Layer on the Datagram instance.
*/
func (dg *Datagram) Write(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("spd.Datagram.Write <-", string(p))

	layers := dg.newLayer()
	errnie.Debugs("spd.Datagram.Write Ptr ->", dg.Ptr())

	if err = layers.Set(int(dg.Ptr()), p); errnie.IOError(err) {
		return 0, err
	}

	errnie.Debugs("spd.Datagram.Write ->", layers)
	return len(p), err
}

/*
Close ...
*/
func (dg *Datagram) Close() error {
	return nil
}
