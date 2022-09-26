package passepartout

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/spd"
)

func TestManagerReadWrite(t *testing.T) {
	Convey("Given a manager", t, func() {
		manager := NewManager()

		Convey("And a value is written", func() {
			dg := spd.NewCached(
				"datapoint", "test", "test.wrkspc.org", "test",
			)

			manager.Write(dg)

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

					manager.Read(q)

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
