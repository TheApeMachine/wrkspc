package infra

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type MockClient struct{}

func (client MockClient) Apply(vendor, name, tag string) {}

func TestNewClient(t *testing.T) {
	Convey("Given a new struct instance", t, func() {
		inst := MockClient{}

		Convey("It should be a Client", func() {
			So(
				inst,
				ShouldHaveSameTypeAs,
				NewClient(inst),
			)
		})
	})
}
