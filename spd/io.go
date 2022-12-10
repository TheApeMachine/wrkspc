package spd

import (
	"io"
	"log"

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
		log.Fatal(err)
		return len(p), errnie.Handles(err)
	}

	if b, err = layers.At(0); err != nil {
		log.Fatal(err)
		return len(p), errnie.Handles(err)
	}

	n = copy(p, b)
	return n, io.EOF
}

func (dg *Datagram) Write(p []byte) (n int, err error) {
	var layers capnp.DataList

	if !dg.HasLayers() {
		if layers, err = capnp.NewDataList(dg.Segment(), 1); err != nil {
			return 0, err
		}

		dg.SetLayers(layers)
	}

	if layers, err = dg.Layers(); err != nil {
		return 0, err
	}

	if err = layers.Set(layers.Len()-1, p); err != nil {
		return 0, err
	}

	return len(p), io.EOF
}

func (dg *Datagram) Close() error {
	return nil
}
