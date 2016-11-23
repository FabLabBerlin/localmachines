package machine_earnings

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"math"
	"time"
)

type Money float64

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

func (me MachineEarning) PayAsYouGoCached() (sum Money) {
	key := fmt.Sprintf("PayAsYouGo(%v)", me.m.Id)

	redis.Cached(key, 3600, &sum, func() interface{} {
		return me.PayAsYouGo()
	})

	return
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

func (me MachineEarning) MembershipsCached() (sum Money) {
	key := fmt.Sprintf("Memberships(%v)", me.m.Id)

	redis.Cached(key, 3600, &sum, func() interface{} {
		return me.Memberships()
	})

	return
}

// Memberships money channel. The proportion is approximated by the
// undiscounted PAYG price.
//
// Why not timebase:
// Imagine 100h 3D printing on Replicator Mini =>  600€ Payg
//          20h Lasercutting                   => 1900€ Payg
//
//
//
//  Memberships(Machine) = sum(Membership(MembershipId, Machine))
//
func (me MachineEarning) Memberships() (sum Money) {
	membershipIds := make(map[int64]bool)

	for _, inv := range me.invs {
		for _, p := range inv.Purchases {
			if p.MachineId != me.m.Id {
				continue
			}
			if !me.ContainsTime(p.TimeStart) {
				continue
			}

			for _, m := range p.Memberships {
				membershipIds[m.Id] = true
			}
		}
	}

	for membershipId := range membershipIds {
		sum += me.Membership(membershipId)
	}

	return
}

//
//                                Undiscounted cost for Machine
//  Membership(Machine) = -------------------------------------------------- * sum(Membership.MonthlyPrice)
//                         Undiscounted cost for all Machines in Membership
//
func (me MachineEarning) Membership(membershipId int64) (sum Money) {
	undiscountedForMachine := Money(0.0)
	undiscountedForMembership := Money(0.0)
	sumMonthlyPrice := Money(0.0)

	// Step 1) Calculate undiscountedForMachine and undiscountedForMembership
	for _, inv := range me.invs {
		for _, p := range inv.Purchases {
			if !me.ContainsTime(p.TimeStart) {
				continue
			}

			for _, m := range p.Memberships {
				if m.Id == membershipId {
					if p.MachineId == me.m.Id {
						undiscountedForMachine += Money(p.PricePerUnit * p.Quantity)
					}
					undiscountedForMembership += Money(p.PricePerUnit * p.Quantity)
				}
			}
		}
	}

	// Step 2) Calculate sumMonthlyPrice
	for _, inv := range me.invs {
		for _, ium := range inv.InvUserMemberships {
			if ium.MembershipId == membershipId {
				sumMonthlyPrice += Money(ium.UserMembership.Membership.MonthlyPrice)
			}
		}
	}

	if math.Abs(float64(undiscountedForMachine)) < 0.01 {
		return 0
	}

	lhs := undiscountedForMachine / undiscountedForMembership
	rhs := sumMonthlyPrice

	return lhs * rhs
}
