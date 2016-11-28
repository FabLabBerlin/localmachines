package retention

import (
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/metrics/retention"
	"github.com/FabLabBerlin/localmachines/tests/models/mocks"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestRetention(t *testing.T) {

	Convey("Testing Retention", t, func() {
		Reset(setup.ResetDB)

		Convey("Calculate", func() {
			Convey("Empty for no data", func() {
				r := retention.New(
					1,
					7,
					day.New(2016, 1, 1),
					day.New(2016, 12, 31),
					[]*invutil.Invoice{},
				)
				triangle := r.Calculate()
				So(len(triangle), ShouldEqual, 53)
				for i, row := range triangle {
					So(len(row.Returns), ShouldEqual, 52-i)
					So(row.StepDays, ShouldEqual, 7)
					So(row.Users, ShouldEqual, 0)
					for _, retrn := range row.Returns {
						So(retrn, ShouldEqual, 0.0)
					}
				}
			})

			Convey("Single staff user in October", func() {
				inv := mocks.LoadInvoice(4165)

				r := retention.New(
					1,
					1,
					day.New(2016, 10, 1),
					day.New(2016, 10, 31),
					[]*invutil.Invoice{
						inv,
					},
				)
				triangle := r.Calculate()
				So(len(triangle), ShouldEqual, 31)
				for i, row := range triangle {
					So(len(row.Returns), ShouldEqual, 30-i)
					So(row.StepDays, ShouldEqual, 1)
					//So(row.Users, ShouldEqual, 1)
					//So(row.NewUsers(), ShouldEqual, []uint64{19})
					for _, retrn := range row.Returns {
						So(retrn, ShouldEqual, 0.0)
					}
				}
			})
		})
	})
}
