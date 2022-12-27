package drknow

import (
	"bytes"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {
	Convey("Given an Abstract instance", t, func() {
		abstract := NewAbstract("test")

		Convey("And there is no Data stored in the Buffer", func() {
			data := ""
			abstract.buffer.WriteString(data)

			Convey("When calling its Read method", func() {
				buf := bytes.NewBuffer([]byte{})
				io.Copy(buf, abstract)

				Convey("It should not read any Data", func() {
					So(buf.Len(), ShouldEqual, 0)
				})
			})
		})

		Convey("And there is Data stored in the Buffer", func() {
			data := "test read data"
			abstract.buffer.WriteString(data)

			Convey("When calling its Read method", func() {
				buf := bytes.NewBuffer([]byte{})
				io.Copy(buf, abstract)

				Convey("It should read the correct Data", func() {
					So(buf.String(), ShouldEqual, data)
				})
			})
		})
	})
}

func TestWrite(t *testing.T) {
	Convey("Given an Abstract instance", t, func() {
		abstract := NewAbstract("test")

		Convey("And there is no Data stored in the Buffer", func() {
			abstract.buffer.Reset()

			Convey("When calling its Write method", func() {
				data := []byte("test write data")
				buf := bytes.NewBuffer(data)
				io.Copy(abstract, buf)

				Convey("It should contain the new Data", func() {
					So(
						abstract.buffer.Bytes(),
						ShouldResemble,
						data,
					)
				})
			})
		})

		Convey("And there is some Data stored in the Buffer", func() {
			abstract.buffer.Reset()
			some := []byte("some")
			abstract.buffer = bytes.NewBuffer(some)

			Convey("When calling its Write method", func() {
				data := []byte("test write data")
				buf := bytes.NewBuffer(data)
				io.Copy(abstract, buf)

				Convey("It should append the new Data", func() {
					res := append(some, data...)

					So(
						abstract.buffer.Bytes(),
						ShouldResemble,
						res,
					)
				})
			})
		})
	})
}

func BenchmarkRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func BenchmarkWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
