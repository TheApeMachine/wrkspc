package tweaker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/test"
)

type BoolContext struct {
	expected map[string]bool
}

func NewBoolContext() *BoolContext {
	test.NewConfigContext()

	return &BoolContext{
		expected: map[string]bool{
			"errnie.local": false,
			"errnie.debug": true,
			"errnie.trace": true,
		},
	}
}

func TestGetBool(t *testing.T) {
	ctx := NewBoolContext()

	Convey("Given a valid key", t, func() {
		for key, value := range ctx.expected {
			Convey(key+" should return the correct value", func() {
				So(GetBool(key), ShouldEqual, value)
			})
		}
	})
}

func BenchmarkGetBool(b *testing.B) {
	test.NewConfigContext()

	for i := 0; i < b.N; i++ {
		_ = GetBool("errnie.local")
		_ = GetBool("errnie.debug")
		_ = GetBool("errnie.trace")
	}
}
