package countries

import (
	"github.com/FabLabBerlin/localmachines/models/countries"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMachineEarnings(t *testing.T) {

	Convey("Testing countries", t, func() {

		Convey("GetByCode", func() {
			c, ok := countries.GetByCode("DE")
			So(ok, ShouldBeTrue)
			So(c.Code, ShouldEqual, "DE")
			So(c.Name, ShouldEqual, "Germany")

			_, ok = countries.GetByCode("LOL123")
			So(ok, ShouldBeFalse)
		})

		Convey("GetAll", func() {
			cs := countries.GetAll()
			So(len(cs), ShouldEqual, 249)
		})

	})
}
