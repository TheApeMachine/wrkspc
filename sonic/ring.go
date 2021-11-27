package sonic

import (
	"container/ring"
)

/*
Ring is a wrapper around ring buffers.
*/
type Ring struct {
	Ptr    int
	Buffer *ring.Ring
}

/*
NewRing returns a reference pointer to a new ring buffer.
*/
func NewRing(size int) *Ring {
	return &Ring{Buffer: ring.New(size)}
}

/*
MoveTo an exact position in the Ring.
*/
func (ring *Ring) MoveTo(pos int) {
	ring.Buffer.Move(pos - ring.Ptr)
	ring.Ptr = pos
}

/*
Peek the current value.
*/
func (ring *Ring) Peek() interface{} {
	return &ring.Buffer.Value
}

/*
Poke a value in the next position.
*/
func (ring *Ring) Poke(value interface{}) {
	ring.Buffer.Next()
	ring.Ptr++
	ring.Buffer.Value = value
}
