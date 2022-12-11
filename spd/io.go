package spd

import (
	"fmt"
	"io"

	capnp "capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Read the contents of the first layer into p.

Make sure p has a pre-allocated length that corresponds to the length
of the data you want to read into it.
*/
func (dg *Datagram) Read(p []byte) (n int, err error) {
	var (
		layers capnp.DataList
		b      = make([]byte, len(p))
	)

	if layers, err = dg.Layers(); err != nil {
		return len(p), err
	}

	if dg.Ptr() < 0 {
		// Make multiple Read calls return layer data in a
		// circular fashion.
		dg.SetPtr(int32(layers.Len() - 1))
	}

	errnie.Debugs(fmt.Sprintf("READ Layer: %d", dg.Ptr()))

	if b, err = layers.At(int(dg.Ptr())); err != nil {
		return len(p), err
	}

	n = copy(p, b)

	// Reduce the Layer Pointer by one layer.
	dg.SetPtr(dg.Ptr() - 1)

	return n, io.EOF
}

/*
Write the contents of p into a new Layer on the Datagram instance.
*/
func (dg *Datagram) Write(p []byte) (n int, err error) {
	var layers capnp.DataList

	if !dg.HasLayers() {
		// If we have never written any data to a Layer, we need to make
		// sure to instantiate a new DataList first.
		if layers, err = capnp.NewDataList(dg.Segment(), 1); err != nil {
			return 0, err
		}

		// Load the DataList into the Datagram Layers property.
		dg.SetLayers(layers)
		dg.SetPtr(-1)
	}

	if layers, err = dg.Layers(); err != nil {
		return 0, err
	}

	// Write to a new Layer.
	dg.SetPtr(int32(layers.Len() - 1))
	errnie.Debugs(fmt.Sprintf("WRITE Layer: %d", dg.Ptr()))

	if err = layers.Set(int(dg.Ptr()), p); err != nil {
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
