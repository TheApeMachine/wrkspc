package hefner

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/spd"
)

func TestNewPipe(t *testing.T) {
	Convey("It should work", t, func() {})
}

func BenchmarkNew(b *testing.B) {
	datagram := &spd.Empty

	for i := 0; i < b.N; i++ {
		_ = NewPipe(datagram)
	}
}
