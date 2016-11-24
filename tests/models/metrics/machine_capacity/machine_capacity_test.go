package machine_capacity

import (
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/metrics/machine_capacity"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestMachineCapacity(t *testing.T) {

	Convey("Testing MachineCapacity", t, func() {
		Convey("Total", func() {
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
