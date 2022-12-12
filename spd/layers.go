package spd

import (
	"io"

	"capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
ReadAt ...
*/
func (dg *Datagram) ReadAt(b []byte, off int64) (n int, err error) {
	if !dg.HasLayers() {
		return 0, io.EOF
	}

	dg.SetPtr(int32(off))
	var p []byte

	if p, err = dg.layers().At(int(dg.Ptr())); err != nil {
		return 0, err
	}

	n = copy(b, p)
	return n, io.EOF
}

func (dg *Datagram) newLayer() capnp.DataList {
	var (
		layers capnp.DataList
		err    error
	)

	if !dg.HasLayers() {
		// If we have never written any data to a Layer, we need to make
		// sure to instantiate a new DataList first.
		if layers, err = capnp.NewDataList(dg.Segment(), 1); err != nil {
			errnie.Handles(err)
			return layers
		}

		// Load the DataList into the Datagram Layers property.
		if err = dg.SetLayers(layers); err != nil {
			errnie.Handles(err)
			return layers
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
