package ford

import (
	"bytes"
	"io"

	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/drknow"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/hefner"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/tui"
	"github.com/theapemachine/wrkspc/twoface"
)

type Workspace struct {
	ID        uuid.UUID
	ctx       *twoface.Context
	tree      *drknow.Tree
	pipes     map[string]*hefner.Pipe
	links     map[string][]io.ReadWriteCloser
	workloads []*Workload
	pool      *twoface.Pool
	err       chan error
}

func NewWorkspace(workloads ...*Workload) *Workspace {
	errnie.Trace()
	ctx := twoface.NewContext()

	return &Workspace{
		uuid.New(),
		ctx,
		drknow.NewTree(),
		make(map[string]*hefner.Pipe),
		make(map[string][]io.ReadWriteCloser),
		workloads,
		twoface.NewPool(ctx).Run(),
		make(chan error),
	}
}

func (wrkspc *Workspace) AddWork(work ...*Workload) *Workspace {
	wrkspc.workloads = append(wrkspc.workloads, work...)
	return wrkspc
}

func (wrkspc *Workspace) Read(p []byte) (n int, err error) {
	errnie.Trace()

	var b []byte

	for _, wrkld := range wrkspc.workloads {
		if n, err = wrkld.Read(b); err != nil {
			return n, errnie.Handles(err)
		}

		dg := &spd.Empty
		if err = dg.Decode(b); errnie.Handles(err) != nil {
			return
		}

		if err = wrkspc.interpret(dg); errnie.Handles(err) != nil {
			return
		}

		var scope spd.ScopeType
		if scope, err = dg.Scope(); errnie.Handles(err) != nil {
			return
		}

		if bytes.Equal(scope, spd.UI) {
			// This is a ui related message, so we should send it
			// to the buffer, which is copied to stdout.
			n = copy(p, b)
		}
	}

	for key, pipe := range wrkspc.pipes {
		for _, link := range wrkspc.links[key] {
			io.Copy(link, pipe)
		}
	}

	return
}

func (wrkspc *Workspace) Write(p []byte) (n int, err error) {
	errnie.Trace()

	for _, workload := range wrkspc.workloads {
		if n, err = workload.Write(p); err != nil {
			wrkspc.Close()
		}
	}

	return
}

func (wrkspc *Workspace) Close() error {
	errnie.Trace()
	errnie.Informs("closing workspace", wrkspc.ID)

	for _, workload := range wrkspc.workloads {
		workload.Close()
	}

	wrkspc.err <- errnie.NewError(nil)
	return nil
}

func (wrkspc *Workspace) interpret(dg *spd.Datagram) error {
	errnie.Trace()

	var (
		role  spd.RoleType
		scope spd.ScopeType
		err   error
	)

	if role, err = dg.Role(); errnie.Handles(err) != nil {
		return err
	}

	if scope, err = dg.Scope(); errnie.Handles(err) != nil {
		return err
	}

	if bytes.Equal(role, spd.PIPE) {
		if bytes.Equal(scope, spd.UI) {
			ui := tui.NewUI(dg)
			wrkspc.pipes[string(scope)] = hefner.NewPipe(ui)

			wrkspc.AddWork(NewWorkload(NewAssembly(
				drknow.NewAbstract(ui),
			)))
		}
	}

	if bytes.Equal(role, spd.LINK) {
		wrkspc.links[string(scope)] = append(
			wrkspc.links[string(scope)], dg,
		)
	}

	return errnie.Handles(nil)
}
