package spd

import (
	"io"

	"capnproto.org/go/capnp/v3"
)

/*
ReadAt ...
*/
func (dg *Datagram) ReadAt(b []byte, off int64) (n int, err error) {
	if !dg.HasLayers() {
		return 0, io.EOF
	}

	var p []byte
	if p, err = dg.layer(off); err != nil {
		return n, err
	}

	n = copy(b, p)
	return n, io.EOF
}

func (dg *Datagram) layer(ptr int64) ([]byte, error) {
	var (
		layers capnp.DataList
		err    error
	)

	if layers, err = dg.Layers(); err != nil {
		return []byte{}, err
	}

	dg.SetPtr(int32(ptr))
	return layers.At(int(dg.Ptr()))
}
