package spd

import (
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {
	Convey("Given a Datagram instance", t, func() {
		datagram := New(APPXMP, TEST, UNIT)

		Convey("And it has a Layer in its Payload", func() {
			data := []byte("test read data")
			datagram.Write(data)

			Convey("When the Read method is called", func() {
				p := make([]byte, 0)
				n, err := datagram.Read(p)

				Convey("It should not report an error", func() {
					// io.EOF is in fact an error, however it is an
					// error that should be treated as nil.
					So(err, ShouldEqual, io.EOF)
				})

				Convey("It should read the correct number of bytes", func() {
					So(n, ShouldEqual, len(data))
				})

				Convey("It should read the correct value", func() {
					So(p, ShouldResemble, data)
				})
			})
		})
	})
}

func TestWrite(t *testing.T) {
	Convey("Given a Datagram instance", t, func() {
		datagram := New(APPXMP, TEST, UNIT)

		Convey("And we have some data to write", func() {
			data := []byte("test write data")

			Convey("When the Write method is called", func() {
				n, err := datagram.Write(data)

				Convey("It should not report an error", func() {
					// io.EOF is in fact an error, however it is an
					// error that should be treated as nil.
					So(err, ShouldBeNil)
				})

				Convey("It should write the correct number of bytes", func() {
					So(n, ShouldEqual, len(data))
				})

				Convey("It should write the correct value to a new Layer", func() {
					var p []byte
					datagram.Read(p)
					So(p, ShouldResemble, data)
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
