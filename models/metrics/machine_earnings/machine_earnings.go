package machine_earnings

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/metrics/filter"
	"github.com/astaxie/beego"
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
		invs: filter.Invoices(invs, from, to),
	}
}

func (me MachineEarning) ContainsTime(t time.Time) bool {
	return !me.from.AfterTime(t) && !me.to.BeforeTime(t)
}

func (me MachineEarning) PayAsYouGoCached() (sum Money, err error) {
	key := fmt.Sprintf("PayAsYouGo(%v)-%v-%v", me.m.Id, me.from, me.to)

	err = redis.Cached(key, 3600, &sum, func() (interface{}, error) {
		return me.PayAsYouGo(), nil
	})

	return
}

func (me MachineEarning) PayAsYouGo() (sum Money) {
	for _, inv := range me.invs {
		for _, p := range inv.Purchases {
			if p.Archived || p.Cancelled {
				continue
			}
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
	key := fmt.Sprintf("Memberships(%v)-%v-%v", me.m.Id, me.from, me.to)

	redis.Cached(key, 3600, &sum, func() (interface{}, error) {
		return me.Memberships(), nil
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
	memberships := make(map[int64]*memberships.Membership)

	for _, inv := range me.invs {
		for _, p := range inv.Purchases {
			if p.MachineId != me.m.Id {
				continue
			}
			if !me.ContainsTime(p.TimeStart) {
				continue
			}

			for _, m := range p.Memberships {
				memberships[m.Id] = m
			}
		}
	}

	for membershipId, m := range memberships {
		affected, err := m.IsMachineAffected(me.m)
		if err != nil {
			beego.Error(err.Error())
		}
		if affected {
			sum += me.Membership(membershipId)
		}
	}

	return
}

//
//                                Undiscounted cost for Machine
//  Membership(Machine) = -------------------------------------------------- * sum(Membership.MonthlyPrice)
//                         Undiscounted cost for all Machines in Membership
//
func (me MachineEarning) Membership(membershipId int64) (sum Money) {
	iumIds := make(map[int64]struct{})
	undiscountedForMachine := Money(0.0)
	undiscountedForMembership := Money(0.0)
	sumMonthlyPrice := Money(0.0)

	// Step 1) Calculate undiscountedForMachine and undiscountedForMembership
	for _, inv := range me.invs {
		for _, p := range inv.Purchases {
			if p.Archived || p.Cancelled {
				continue
			}

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
				if _, exists := iumIds[ium.Id]; exists {
					beego.Error("duplicate sum!!!")
				}
				iumIds[ium.Id] = struct{}{}
				sumMonthlyPrice += Money(ium.UserMembership.Membership.MonthlyPrice)
			}
		}
	}

	if math.Abs(float64(undiscountedForMachine)) < 0.01 {
		return 0
	}

	lhs := undiscountedForMachine / undiscountedForMembership
	if lhs > 1.1 {
		beego.Error("lhs > 1.1")
	}
	rhs := sumMonthlyPrice

	return lhs * rhs
}
