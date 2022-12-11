package hefner

import (
	"bytes"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type PipeTestContext struct {
	pipe *Pipe
}

func NewPipeTestContext() *PipeTestContext {
	return &PipeTestContext{
		pipe: NewPipe(IPC),
	}
}

func TestNewPipe(t *testing.T) {
	ctx := NewPipeTestContext()

	Convey("Given a pointer to a new instance of Pipe", t, func() {
		pipe := ctx.pipe

		Convey("It should have a read output", func() {
			So(pipe.r, ShouldNotBeNil)
		})

		Convey("It should have a write input", func() {
			So(pipe.w, ShouldNotBeNil)
		})

		Convey("It should read from the write input", func() {
			data := []byte("test data")
			recv := make([]byte, len(data))

			pipe.Write(data)
			pipe.Read(recv)

			So(recv, ShouldResemble, data)
		})

		Convey("It should copy write to read", func() {
			data := bytes.NewBuffer([]byte("test data"))
			recv := bytes.NewBuffer([]byte{})

			io.Copy(pipe, data)

			So(recv, ShouldResemble, data)
		})
	})
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPipe(IPC)
	}
}
