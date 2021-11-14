package sonic

import (
	"container/ring"
)

type Ring struct {
	Ptr    int
	Buffer *ring.Ring
}

func NewRing(size int) *Ring {
	return &Ring{Buffer: ring.New(size)}
}

func (ring *Ring) MoveTo(pos int) {
	ring.Buffer.Move(pos - ring.Ptr)
	ring.Ptr = pos
}

func (ring *Ring) Peek(value interface{}) {
	value = &ring.Buffer.Prev().Value
	ring.Buffer.Next()
}

func (ring *Ring) Poke(value interface{}) {
	ring.Buffer.Value = value
	ring.Buffer.Next()
	ring.Ptr++
}
