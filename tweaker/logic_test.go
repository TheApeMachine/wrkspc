package tweaker

import (
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/test"
)

type LogicContext struct {
	expected map[string]int
}

func NewLogicContext() *LogicContext {
	test.NewConfigContext()

	return &LogicContext{
		expected: map[string]int{
			"twoface.pool.workers": runtime.NumCPU() * 2,
		},
	}
}

func TestGetLogic(t *testing.T) {
	ctx := NewLogicContext()

	Convey("Given a valid key", t, func() {
		for key, value := range ctx.expected {
			Convey(key+" should return the correct value", func() {
				So(GetLogic(key), ShouldEqual, value)
			})
		}
	})
}

func BenchmarkGetLogic(b *testing.B) {
	test.NewConfigContext()

	for i := 0; i < b.N; i++ {
		_ = GetLogic("twoface.pool.workers")
	}
}
