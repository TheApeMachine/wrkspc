package spd

import (
	"io"

	"capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Next returns the Layer the Ptr is currently pointing to, then advances
the Ptr one step.
*/
func (dg *Datagram) Next() []byte {
	errnie.Trace()

	var (
		p   []byte
		err error
	)

	if p, err = dg.ReadAt(int64(dg.Ptr())); errnie.Handles(err) != nil {
		return p
	}

	dg.SetPtr(dg.Ptr() + 1)

	return p
}

/*
ReadAt read the entire payload of a Layer at a specified offset.

This was designed to work as a circular buffer, so if we assume a
DataList of [0, 1, 2], off = 2 returns 2 and off = 3 returns 0.
*/
func (dg *Datagram) ReadAt(off int64) (p []byte, err error) {
	errnie.Trace()

	if !dg.HasLayers() {
		// Signals to the caller that a read was attempted before
		// any data was ever written.
		return []byte{}, io.ErrUnexpectedEOF
	}

	// Set the Layer Pointer, making it circular if it is set outside
	// of its length in either direction.
	dg.SetPtr(int32(off) % int32(dg.layers().Len()))
	errnie.Debugs("spd.Datagram.ReadAt Ptr ->", dg.Ptr())

	// Get the layer we are interested in.
	return dg.layers().At(int(dg.Ptr()))
}

func (dg *Datagram) newLayer() capnp.DataList {
	errnie.Trace()

	var (
		layers capnp.DataList
		err    error
	)

	if !dg.HasLayers() {
		if layers, err = dg.NewLayers(1); err != nil {
			errnie.Handles(err)
		}
	}

	dg.SetPtr(int32(dg.layers().Len() - 1))
	errnie.Debugs("spd.Datagram.newLayer Ptr ->", dg.Ptr())
	return layers
}

func (dg *Datagram) layers() capnp.DataList {
	errnie.Trace()

	var (
		layers capnp.DataList
		err    error
	)

	if layers, err = dg.Layers(); errnie.Handles(err) != nil {
		return layers
	}

	return layers
}
