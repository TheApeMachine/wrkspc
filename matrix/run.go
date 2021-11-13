package matrix

import (
	"context"
	"path/filepath"

	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Run is a wrapper that takes a container and defines a way to run it.
*/
type Run struct {
	ctx       context.Context
	root      string
	name      string
	container Container
}

/*
NewRun constructs an instance of Run and returns it.
TODO: Context should come from here.
*/
func NewRun(name string) Run {
	errnie.Traces()

	return Run{
		name:      name,
		ctx:       context.Background(),
		root:      filepath.FromSlash(brazil.HomePath() + ".wrkspc/"),
		container: NewContainer(name),
	}
}

/*
Cycle executes a Run.
*/
func (run Run) Cycle() error {
	errnie.Traces()
	return run.container.Run()
}
