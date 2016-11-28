package retention

import (
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/metrics/retention"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/FabLabBerlin/localmachines/tests/models/mocks"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
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
					[]*users.User{},
				)
				triangle := r.Calculate()
				So(len(triangle), ShouldEqual, 52)
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

				us := []*users.User{
					{
						Id:      19,
						Created: time.Date(2016, time.October, 1, 12, 0, 0, 0, time.UTC),
					},
				}

				r := retention.New(
					1,
					1,
					day.New(2016, 10, 1),
					day.New(2016, 10, 31),
					[]*invutil.Invoice{
						inv,
					},
					us,
				)
				triangle := r.Calculate()
				So(len(triangle), ShouldEqual, 30)

				// User who signed up on Oct 1, did prints on
				// Oct 4, 6, 12, 14, 20, 21, 31
				for i, row := range triangle {
					So(len(row.Returns), ShouldEqual, 30-i)
					So(row.StepDays, ShouldEqual, 1)
					if i == 0 {
						So(row.Users, ShouldEqual, 1)
						So(row.NewUsers(), ShouldResemble, []int64{19})
						for j, percentage := range row.Returns {
							switch j {
							case 3, 5, 11, 13, 19, 20, 30:
								So(percentage, ShouldEqual, 1.0)
							default:
								So(percentage, ShouldEqual, 0)
							}
						}
					} else {
						So(row.Users, ShouldEqual, 0)
						So(row.NewUsers(), ShouldResemble, []int64{})
						for _, retrn := range row.Returns {
							So(retrn, ShouldEqual, 0.0)
						}
					}
				}
			})
		})
	})
}
