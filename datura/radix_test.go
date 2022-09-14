package datura

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/spd"
)

func TestRadixRead(t *testing.T) {
	Convey("Given a radix tree", t, func() {
		tree := NewRadix()

		Convey("And a value is written", func() {
			dg := spd.NewCached(
				"datapoint", "test", "test.wrkspc.org", "test",
			)

			tree.Write(dg)

			Convey("It should be able to retrieve the value", func() {
				q := spd.NewCached(
					"datapoint", "test", "test.wrkspc.org",
					"v4.0.0/datapoint/test/test.wrkspc.org",
				)

				tree.Read(q)

				So(
					string(spd.Unmarshal(q).Payload()),
					ShouldEqual,
					string(spd.Unmarshal(dg).Payload()),
				)
			})
		})
	})
}
