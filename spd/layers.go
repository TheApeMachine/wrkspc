package spd

import (
	"io"

	"capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

/* Next ... */
func (dg *Datagram) Next() (p []byte) {
	if _, err := dg.Read(p); errnie.Handles(err) != nil {
		errnie.Inspects(p)
		return
	}

	return
}

/*
ReadAt read the entire payload of a Layer at a specified offset.

This was designed to work as a circular buffer, so if we assume a
DataList of [0, 1, 2], off = 2 returns 2 and off = 3 returns 0.
*/
func (dg *Datagram) ReadAt(b []byte, off int64) (n int, err error) {
	if !dg.HasLayers() {
		// Signals to the caller that a read was attempted before
		// any data was ever written.
		return 0, io.ErrUnexpectedEOF
	}

	// Set the Layer Pointer, making it circular if it is set outside
	// of its length in either direction.
	dg.SetPtr(int32(off) % int32(dg.layers().Len()))
	var p []byte

	// Get the layer we are interested in.
	if p, err = dg.layers().At(int(dg.Ptr())); err != nil {
		return 0, err
	}

	// Grow the buffer if needed.
	if lb, lp := len(b), len(p); lb < lp {
		b = b[:lp]
	}

	// Move the layer into the read buffer, and return io.EOF to
	// signal a graceful return.
	n = copy(b, p)
	return n, io.EOF
}

func (dg *Datagram) newLayer() capnp.DataList {
	var (
		layers capnp.DataList
		err    error
	)

	if !dg.HasLayers() {
		if _, err = dg.NewLayers(1); err != nil {
			errnie.Handles(err)
		}
	}

	if layers, err = dg.Layers(); err != nil {
		errnie.Handles(err)
		return layers
	}

	dg.SetPtr(int32(layers.Len() - 1))
	return layers
}

func (dg *Datagram) layers() capnp.DataList {
	var (
		layers capnp.DataList
		err    error
	)

	if layers, err = dg.Layers(); errnie.Handles(err) != nil {
		return layers
	}

	return layers
}
