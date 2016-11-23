package machine_earnings

import (
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/astaxie/beego"
	"math"
	"time"
)

type Money float64

type T struct {
	locationId int64
	invoices   []*invutil.Invoice
}

type MachineEarning struct {
	m    *machine.Machine
	from month.Month
	to   month.Month
	invs []*invutil.Invoice
}

func New(
	m *machine.Machine,
	from month.Month,
	to month.Month,
	invs []*invutil.Invoice,
) *MachineEarning {

	return &MachineEarning{
		m:    m,
		from: from,
		to:   to,
		invs: invs,
	}
}

func (me MachineEarning) ContainsInvoice(inv *invutil.Invoice) bool {
	m := month.New(time.Month(inv.Month), inv.Year)
	return me.from.BeforeOrEqual(m) && me.to.AfterOrEqual(m)
}

func (me MachineEarning) ContainsTime(t time.Time) bool {
	return !me.from.AfterTime(t) && !me.to.BeforeTime(t)
}

func (me MachineEarning) PayAsYouGo() (sum Money) {
	for _, inv := range me.invs {
		for _, p := range inv.Purchases {
			if !me.ContainsTime(p.TimeStart) {
				continue
			}

			if p.MachineId == me.m.Id {
				sum += Money(p.DiscountedTotal)
			}
		}
	}

	return
}

type UserMembershipId int64

// Memberships money channel. The proportion is approximated by the
// undiscounted PAYG price.
//
// Why not timebase:
// Imagine 100h 3D printing on Replicator Mini =>  600€ Payg
//          20h Lasercutting                   => 1900€ Payg
//
//
//                   w[1]*sum(Membership[1]) + ... + w[n]*sum(Membership[n])
//  Memberships() = ---------------------------------------------------------
//                                 w[1] + ... + w[n]
//
//   where
//
//  w[i] = sum(Purchase.Undiscounted, Purchase affected by Membership[i]),
//  sum(Membership[i]) = sum(InvUserMembership.MonthlyPrice, Invoice in [from, to] and
//                                                           MembershipId fits)
//
func (me MachineEarning) Memberships() (sum Money) {
	ms := make(map[int64]*memberships.Membership)
	ws := make(map[int64]Money)

	// Step 1: obtain w[i] and m[i]
	for _, inv := range me.invs {
		for _, p := range inv.Purchases {
			if p.MachineId != me.m.Id {
				continue
			}
			if !me.ContainsTime(p.TimeStart) {
				continue
			}

			if len(p.Memberships) > 1 {
				beego.Warn("len(p.Memberships) > 1")
			}

			for _, m := range p.Memberships {
				ws[m.Id] += Money(p.PricePerUnit * p.Quantity)
				ms[m.Id] = m
			}
		}
	}

	// Step 2: plug in w[i] and get membership earnings on-the-fly
	enumerator := 0.0

	for membershipId := range ms {
		innerSum := Money(0)
		for _, inv := range me.invs {
			for _, ium := range inv.InvUserMemberships {
				if ium.MembershipId == membershipId {
					innerSum += Money(ium.Membership().MonthlyPrice)
				}
			}
		}

		enumerator += float64(innerSum) * float64(ws[membershipId])
	}

	denominator := 0.0
	for _, w := range ws {
		denominator += float64(w)
	}

	if math.Abs(enumerator) < 0.01 {
		return 0
	}

	return Money(enumerator / denominator)
}
