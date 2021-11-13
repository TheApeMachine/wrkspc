package kube

import (
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
MigratableKind I do not remember why I called it this.
*/
type MigratableKind interface {
	Up()
	Check() bool
	Delete() error
}

/*
Base holds some shared functionality for Kinds.
*/
type Base struct {
	file   []byte
	kind   MigratableKind
	client RestClient
	err    error
}

/*
NewBase constructs a new Base for a Kind.
*/
func NewBase(file []byte, kind MigratableKind, client RestClient) *Base {
	return &Base{
		file:   file,
		kind:   kind,
		client: client,
	}
}

/*
waiter hangs around until it either finds the kind fully deployed and working, or it times out.
*/
func (base *Base) waiter() {
	// Run the liveness check as a retrying process.
	twoface.NewRepeater(10, twoface.NewRetryStrategy(
		twoface.Fibonacci{MaxTries: 10},
	)).Attempt(1, base.kind.Check)
}

func (base *Base) teardown() error {
	err := base.kind.Delete()
	errnie.Handles(err).With(errnie.KILL)
	base.waiter()
	return err
}
