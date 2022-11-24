package tweaker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type IntContext struct {
	expected map[string]int
}

func NewIntContext() *IntContext {
	NewConfigContext()

	return &IntContext{
		expected: map[string]int{
			"twoface.pool.job.buffer": 256,
		},
	}
}

func TestGetInt(t *testing.T) {
	ctx := NewIntContext()

	Convey("Given a valid key", t, func() {
		for key, value := range ctx.expected {
			Convey(key+" should return the correct value", func() {
				So(GetInt(key), ShouldEqual, value)
			})
		}
	})
}

func BenchmarkGetInt(b *testing.B) {
	NewConfigContext()

	for i := 0; i < b.N; i++ {
		_ = GetBool("twoface.pool.job.buffer")
	}
}
