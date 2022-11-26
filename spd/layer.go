package spd

import (
	"bytes"

	"github.com/theapemachine/wrkspc/errnie"
)

func (dg Datagram) Write() []*bytes.Buffer {
	list, err := dg.Layers()
	errnie.Handles(err)

	count := list.Len()

	buf := make([]*bytes.Buffer, count)

	for idx := 0; idx < count; idx++ {
		data, err := list.At(idx)
		errnie.Handles(err)
		buf[idx] = bytes.NewBuffer(data)

	}

	return buf
}
