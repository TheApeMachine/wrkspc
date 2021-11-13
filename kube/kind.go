package kube

import (
	"time"

	"github.com/theapemachine/wrkspc/errnie"
)

type MigratableKind interface {
	Up() error
	Check() bool
	Down() error
	Delete() error
	Name() string
}

type Base struct {
	kind MigratableKind
}

func NewBase(kind MigratableKind) Base {
	return Base{
		kind: kind,
	}
}

func (base Base) waiter(direction bool) {
	count := 0

	// We needed some form of backoff and give up thingy, otherwise when there is
	// a problem we're essentially stuck, and I don't like being stuck.
	// TODO: Needs some way to determine if the deployment is significantly compromised
	//       or reporting on a minor error is enough and it can be fixed later.
	for {
		//count++
		count = 1 // temp override for testing.
		time.Sleep(time.Duration(count) * time.Second)

		if base.kind.Check() == direction || count == 10 {
			break
		}
	}
}

func (base Base) teardown() error {
	err := base.kind.Delete()
	errnie.Handles(err).With(errnie.KILL)
	base.waiter(false)
	return err
}
