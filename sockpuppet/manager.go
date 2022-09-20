package sockpuppet

import "github.com/theapemachine/wrkspc/twoface"

type Manager interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	PoolSize() int
	Manage(*twoface.Context)
}
