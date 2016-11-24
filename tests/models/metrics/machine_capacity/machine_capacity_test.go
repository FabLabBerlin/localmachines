package machine_capacity

import (
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/metrics/machine_capacity"
	"github.com/FabLabBerlin/localmachines/tests/models/mocks"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestMachineCapacity(t *testing.T) {

	Convey("Testing MachineCapacity", t, func() {
		Convey("Opening", func() {
			Reset(setup.ResetDB)

			inv := mocks.LoadInvoice(4165)

			m := &machine.Machine{
				Id: 2,
			}

			mc := machine_capacity.New(
				m,
				month.New(1, 2016),
				month.New(12, 2016),
				[]*invutil.Invoice{
					inv,
				},
			)

			So(mc.Opening().Day(), ShouldEqual, 21)
			So(mc.Opening().Month(), ShouldEqual, 10)
			So(mc.Opening().Year(), ShouldEqual, 2016)
		})

		Convey("Usage", func() {
			Reset(setup.ResetDB)

			inv := mocks.LoadInvoice(4165)

			m := &machine.Machine{
				Id: 2,
			}

			mc := machine_capacity.New(
				m,
				month.New(1, 2016),
				month.New(12, 2016),
				[]*invutil.Invoice{
					inv,
				},
			)

			So(math.Abs(float64(mc.Usage().Hours())-6) < 0.1, ShouldBeTrue)
		})

		Convey("Utilization", func() {
			Reset(setup.ResetDB)

			Convey("0 when DB empty", func() {
				m := &machine.Machine{
					Id: 3,
				}
				mc := machine_capacity.New(
					m,
					month.New(1, 2016),
					month.New(12, 2016),
					[]*invutil.Invoice{},
				)
				So(mc.Total(), ShouldEqual, 0)
			})
		})
	})

}
