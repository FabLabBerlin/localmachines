package machine_earnings

import (
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/metrics/machine_earnings"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships/inv_user_memberships"
	"github.com/FabLabBerlin/localmachines/tests/models/mocks"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestMachineEarnings(t *testing.T) {

	Convey("Testing MachineEarnings", t, func() {
		Reset(setup.ResetDB)

		Convey("PayAsYouGo", func() {
			inv := mocks.LoadInvoice(4165)

			data := &invutil.PrefetchedData{
				LocationId: 1,
				UmbsByUid:  make(map[int64][]*user_memberships.UserMembership),
				IumbsByUid: make(map[int64][]*inv_user_memberships.InvoiceUserMembership),
			}

			data.UmbsByUid[19] = []*user_memberships.UserMembership{
				mocks.LoadUserMembership(14),
				mocks.LoadUserMembership(15),
			}

			for _, ium := range inv.InvUserMemberships {
				data.IumbsByUid[19] = append(data.IumbsByUid[19], ium)
			}

			m := &machine.Machine{
				Id: 14,
			}

			me := machine_earnings.New(
				m,
				day.New(2016, 1, 1),
				day.New(2016, 12, 31),
				[]*invutil.Invoice{
					inv,
				},
				data,
			)

			So(float64(me.PayAsYouGo()), ShouldEqual, 0)
		})
	})

}
