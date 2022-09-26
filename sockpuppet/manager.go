package sockpuppet

import "io"

type Manager interface {
	io.ReadWriter
	PoolSize() int
}
