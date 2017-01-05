package bin

import (
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models/metrics/bin"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBin(t *testing.T) {

	Convey("Testing bin", t, func() {
		Convey("Map", func() {
			from := month.New(2015, 8)
			to := month.New(2017, 1)

			i, ok := bin.Map(from, to, from)
			So(i, ShouldEqual, 0)
			So(ok, ShouldBeTrue)

			i, ok = bin.Map(from, to, month.New(2015, 9))
			So(i, ShouldEqual, 1)
			So(ok, ShouldBeTrue)

			i, ok = bin.Map(from, to, month.New(2015, 11))
			So(i, ShouldEqual, 3)
			So(ok, ShouldBeTrue)

			i, ok = bin.Map(from, to, month.New(2015, 12))
			So(i, ShouldEqual, 4)
			So(ok, ShouldBeTrue)

			i, ok = bin.Map(from, to, month.New(2016, 12))
			So(i, ShouldEqual, 16)
			So(ok, ShouldBeTrue)

			i, ok = bin.Map(from, to, month.New(2017, 1))
			So(i, ShouldEqual, 17)
			So(ok, ShouldBeTrue)
		})
	})
}
