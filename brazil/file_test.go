package brazil

import (
	"bytes"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/test"
)

type FileTestContext struct{}

func NewFileTestContext() *FileTestContext {
	test.NewConfigContext()

	return &FileTestContext{}
}

func TestNewFile(t *testing.T) {
	NewFileTestContext()

	Convey("Given a path that does not exist", t, func() {
		path := "/tmp/noexist"

		Convey("And a file name that does not exist", func() {
			name := "noexist.tmp"

			Convey("And some data to write to a file", func() {
				data := bytes.NewBufferString("tmp data")

				Convey("When I try to instantiate the file", func() {
					fh := NewFile(path, name, data)

					Convey("It should create the path", func() {
						_, err := os.Stat(path)
						So(err, ShouldEqual, os.IsExist(err))
					})

					Convey("It should create the file", func() {
						_, err := os.Stat(path + "/" + name)
						So(err, ShouldEqual, os.IsExist(err))
					})

					Convey("It should contain the data", func() {
						So(fh.Data, ShouldEqual, data)
					})
				})
			})
		})
	})
}

func BenchmarkGetBool(b *testing.B) {
	NewFileTestContext()

	for i := 0; i < b.N; i++ {
		NewFile("/tmp", "tmp.tmp", bytes.NewBufferString("tmp"))
	}
}
