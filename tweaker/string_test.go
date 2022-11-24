package tweaker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type StringContext struct {
	expected map[string]string
}

func NewStringContext() *StringContext {
	NewConfigContext()

	return &StringContext{
		expected: map[string]string{
			"server.port": "1984",
		},
	}
}

func TestGetString(t *testing.T) {
	ctx := NewStringContext()

	Convey("Given a valid key", t, func() {
		for key, value := range ctx.expected {
			Convey(key+" should return the correct value", func() {
				So(GetString(key), ShouldEqual, value)
			})
		}
	})
}

func BenchmarkGetString(b *testing.B) {
	NewConfigContext()

	for i := 0; i < b.N; i++ {
		_ = GetString("server.port")
	}
}
