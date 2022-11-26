package spd

import "github.com/theapemachine/wrkspc/errnie"

func (dg Datagram) Read(p []byte) (n int, err error) {
	buf, err := dg.Message().MarshalPacked()
	if errnie.Handles(err) != nil {
		return
	}

	copy(p, buf)
	return len(p), err
}

func (dg Datagram) Write(p []byte) (n int, err error) {
	dg.Unmarshal(p)

	if dg.HasUuid() {
		return dg.append(p)
	}

	return len(p), err
}

func (dg Datagram) append(p []byte) (n int, err error) {
	layers, err := dg.Layers()
	layers.Set(layers.Len(), p)
	return len(p), err
}
