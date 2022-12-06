package ford

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/drknow"
)

type WorkspaceTestContext struct {
	instance *Workspace
}

func NewWorkspaceTestContext() *WorkspaceTestContext {
	return &WorkspaceTestContext{
		instance: NewWorkspace(
			NewWorkload([]*Assembly{
				NewAssembly(drknow.NewAbstract()),
			}...),
		),
	}
}

func TestNew(t *testing.T) {
	Convey("Given a new instance", t, func() {
	})
}

func BenchmarkNew(b *testing.B) {
	NewWorkspaceTestContext()

	for i := 0; i < b.N; i++ {
		_ = NewWorkspace()
	}
}
