package memberstats

import (
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models/metrics/memberstats"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMemberstats(t *testing.T) {

	Convey("Testing Memberstats", t, func() {
		Convey("MapBin", func() {
			from := month.New(2015, 8)
			to := month.New(2017, 1)

			i, ok := memberstats.MapBin(from, to, from)
			So(i, ShouldEqual, 0)
			So(ok, ShouldBeTrue)

			i, ok = memberstats.MapBin(from, to, month.New(2015, 9))
			So(i, ShouldEqual, 1)
			So(ok, ShouldBeTrue)

			i, ok = memberstats.MapBin(from, to, month.New(2015, 11))
			So(i, ShouldEqual, 3)
			So(ok, ShouldBeTrue)

			i, ok = memberstats.MapBin(from, to, month.New(2015, 12))
			So(i, ShouldEqual, 4)
			So(ok, ShouldBeTrue)

			i, ok = memberstats.MapBin(from, to, month.New(2016, 12))
			So(i, ShouldEqual, 16)
			So(ok, ShouldBeTrue)

			i, ok = memberstats.MapBin(from, to, month.New(2017, 1))
			So(i, ShouldEqual, 17)
			So(ok, ShouldBeTrue)
		})
	})
}
