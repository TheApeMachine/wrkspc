package datura

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

func BenchmarkWrite(b *testing.B) {
	errnie.Tracing(false)
	store := NewS3()

	for i := 0; i < b.N; i++ {
		store.Write(spd.NewCached(
			"datapoint", "test", "test.wrkspc.org", "test",
		))
	}
}

func TestS3WriteRead(t *testing.T) {
	Convey("Given an S3 connection", t, func() {
		tree := NewS3()

		Convey("And a value is written", func() {
			dg := spd.NewCached(
				"datapoint", "test", "test.wrkspc.org", "test",
			)

			tree.Write(dg)

			Convey("It should be able to retrieve the value", func() {
				for _, key := range []string{
					"v4.0.0/datapoint/test/test.wrkspc.org",
					"v4.0.0/datapoint/test/test.wr",
					"v4.0.0/datapoint/",
				} {
					q := spd.NewCached(
						"datapoint", "test", "test.wrkspc.org",
						key,
					)

					tree.Read(q)

					So(
						string(spd.Unmarshal(q).Payload()),
						ShouldEqual,
						string(spd.Unmarshal(dg).Payload()),
					)
				}
			})
		})
	})
}
