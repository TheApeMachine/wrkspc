package spd

import (
	"io"
	"testing"

	"capnproto.org/go/capnp/v3"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/errnie"
)

func readLayer(datagram *Datagram, t *testing.T) []byte {
	layers, err := datagram.Layers()
	if errnie.Handles(err) != nil {
		t.FailNow()
	}

	var data []byte
	if data, err = layers.At(0); err != nil {
		t.FailNow()
	}

	return data
}

func writeLayer(datagram *Datagram, data []byte, t *testing.T) {
	dl, err := capnp.NewDataList(datagram.Segment(), 1)
	if errnie.Handles(err) != nil {
		t.FailNow()
	}

	datagram.SetLayers(dl)

	layers, err := datagram.Layers()
	if errnie.Handles(err) != nil {
		t.FailNow()
	}

	if err := layers.Set(0, data); errnie.Handles(err) != nil {
		t.FailNow()
	}

	datagram.SetLayers(layers)
}

func TestRead(t *testing.T) {
	Convey("Given a Datagram instance", t, func() {
		datagram := New(APPXMP, TEST, UNIT)

		Convey("And it has a Layer in its Payload", func() {
			data := []byte("test read data")
			writeLayer(datagram, data, t)

			Convey("When the Read method is called", func() {
				p := make([]byte, len(data))
				n, err := datagram.Read(p)

				Convey("It should not report an error", func() {
					// io.EOF is in fact an error, however it is an
					// error that should be treated as nil.
					So(err, ShouldEqual, io.EOF)
				})

				Convey("It should read the correct number of bytes", func() {
					So(n, ShouldEqual, len(p))
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
					So(err, ShouldEqual, io.EOF)
				})

				Convey("It should write the correct number of bytes", func() {
					So(n, ShouldEqual, len(data))
				})

				Convey("It should write the correct value to a new Layer", func() {
					p := readLayer(datagram, t)
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
