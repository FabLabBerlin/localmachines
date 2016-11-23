package machine_earnings

import (
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
)

type Money float64

type T struct {
	locationId int64
	invoices   []*invutil.Invoice
}

type MachineEarning struct {
	m    *machine.Machine
	from day.Day
	to   day.Day
	invs []*invutil.Invoice
	data *invutil.PrefetchedData
}

func New(
	m *machine.Machine,
	from day.Day,
	to day.Day,
	invs []*invutil.Invoice,
	data *invutil.PrefetchedData,
) *MachineEarning {

	return &MachineEarning{
		m:    m,
		from: from,
		to:   to,
		invs: invs,
		data: data,
	}
}

func (me MachineEarning) PayAsYouGo() (sum Money) {
	for _, inv := range me.invs {
		for _, p := range inv.Purchases {
			if p.MachineId == me.m.Id {
				sum += Money(p.DiscountedTotal)
			}
		}
	}

	return
}

// Memberships money channel. The proportion is approximated by the
// undiscounted PAYG price.
//
// Why not timebase:
// Imagine 100h 3D printing on Replicator Mini =>  600€ Payg
//          20h Lasercutting                   => 1900€ Payg
func (me MachineEarning) Memberships() (sum Money) {
	return
}
