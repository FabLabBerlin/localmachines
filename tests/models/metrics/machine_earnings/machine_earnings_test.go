package machine_earnings

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/metrics/machine_earnings"
	"github.com/FabLabBerlin/localmachines/tests/models/mocks"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestMachineEarnings(t *testing.T) {

	Convey("Testing MachineEarnings", t, func() {
		Reset(setup.ResetDB)

		Convey("PayAsYouGo", func() {
			Convey("Staff pays nothing", func() {
				inv := mocks.LoadInvoice(4165)

				m := &machine.Machine{
					Id: 14,
				}

				me := machine_earnings.New(
					m,
					month.New(2016, 1),
					month.New(2016, 12),
					[]*invutil.Invoice{
						inv,
					},
				)

				So(float64(me.PayAsYouGo()), ShouldEqual, 0)
			})

			Convey("User without membership pays it all", func() {
				inv := mocks.LoadInvoice(4928)

				m := &machine.Machine{
					Id: 17,
				}

				me := machine_earnings.New(
					m,
					month.New(1, 2016),
					month.New(12, 2016),
					[]*invutil.Invoice{
						inv,
					},
				)
				fmt.Printf("me.PayAsYouGo()=%v\n", me.PayAsYouGo())
				So(math.Abs(float64(me.PayAsYouGo())-164.15) < 0.01, ShouldBeTrue)
			})
		})

		Convey("Memberships", func() {
			Convey("Using only one Machine should lead to 100 percent", func() {
				inv := mocks.LoadInvoice(5936)

				m := &machine.Machine{
					Id: 3,
				}

				me := machine_earnings.New(
					m,
					month.New(1, 2016),
					month.New(12, 2016),
					[]*invutil.Invoice{
						inv,
					},
				)

				So(float64(me.Memberships()), ShouldEqual, 150)
			})

			Convey("Using only two Machines 20%% / 80%% should be split accordingly", func() {
				inv := mocks.LoadInvoice(4402)

				// 66.4 cents ~ 21%
				m83Seconds80CentsPerMinute := &machine.Machine{
					Id: 3,
				}

				// 256 cents ~ 79%
				m160Seconds160CentsPerMinute := &machine.Machine{
					Id: 17,
				}

				me83Seconds80CentsPerMinute := machine_earnings.New(
					m83Seconds80CentsPerMinute,
					month.New(1, 2016),
					month.New(12, 2016),
					[]*invutil.Invoice{
						inv,
					},
				)

				me160Seconds160CentsPerMinute := machine_earnings.New(
					m160Seconds160CentsPerMinute,
					month.New(1, 2016),
					month.New(12, 2016),
					[]*invutil.Invoice{
						inv,
					},
				)
				So(math.Abs(float64(me83Seconds80CentsPerMinute.Memberships())-150*0.21) < 0.4, ShouldBeTrue)
				So(math.Abs(float64(me160Seconds160CentsPerMinute.Memberships())-150*0.79) < 0.4, ShouldBeTrue)
			})
		})
	})

}
