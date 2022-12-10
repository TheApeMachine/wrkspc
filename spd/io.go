package spd

import (
	capnp "capnproto.org/go/capnp/v3"
	"github.com/theapemachine/wrkspc/errnie"
)

func (dg Datagram) Read(p []byte) (n int, err error) {
	var layers capnp.DataList

	if layers, err = dg.Layers(); err != nil {
		return len(p), errnie.Handles(err)
	}

	p, err = layers.At(0)
	return len(p), errnie.Handles(err)
}

func (dg Datagram) Write(p []byte) (n int, err error) {
	layers, err := dg.Layers()
	layers.Set(layers.Len(), p)
	return len(p), err
}

func (dg Datagram) Close() error {
	return nil
}
