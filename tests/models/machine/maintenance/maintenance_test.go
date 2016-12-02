package maintenance

import (
	"github.com/FabLabBerlin/localmachines/models/machine/maintenance"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func init() {
	setup.ConfigDB()
}

func TestMachine(t *testing.T) {
	Convey("Testing Maintenance model", t, func() {
		Reset(setup.ResetDB)
		Convey("On", func() {
			Convey("Works", func() {
				mt, err := maintenance.On(123)
				if err != nil {
					panic(err.Error())
				}
				So(mt.Id, ShouldNotEqual, 0)

				So(mt.Start.Year(), ShouldEqual, time.Now().Year())
				So(mt.Start.Month(), ShouldEqual, time.Now().Month())
				So(mt.Start.Day(), ShouldEqual, time.Now().Day())
			})
		})
	})
}
