package ford

import (
	"io"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/drknow"
	"github.com/theapemachine/wrkspc/spd"
)

func TestNew(t *testing.T) {
	Convey("Given a new instance", t, func() {
	})
}

func BenchmarkNew(b *testing.B) {
	errpipe := drknow.NewAbstract(spd.New(spd.APPBIN, spd.PIPE, spd.WAN))

	for i := 0; i < b.N; i++ {
		_ = NewWorkspace(NewWorkload(NewAssembly(errpipe)))
	}
}

func BenchmarkRead(b *testing.B) {
	errpipe := drknow.NewAbstract(spd.New(spd.APPBIN, spd.PIPE, spd.WAN))
	benchSPC := NewWorkspace(NewWorkload(NewAssembly(errpipe)))

	for i := 0; i < b.N; i++ {
		io.Copy(os.Stdout, benchSPC)
	}
}

func BenchmarkWrite(b *testing.B) {
	errpipe := drknow.NewAbstract(spd.New(spd.APPBIN, spd.PIPE, spd.WAN))
	benchSPC := NewWorkspace(NewWorkload(NewAssembly(errpipe)))
	datagram := spd.New(spd.APPXMP, spd.TEST, spd.BENCHMARK)

	for i := 0; i < b.N; i++ {
		io.Copy(benchSPC, datagram)
	}
}
