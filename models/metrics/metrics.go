/*
metrics package for basic visualization of numbers we have.
*/
package metrics

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/metrics/bin"
	"github.com/FabLabBerlin/localmachines/models/metrics/filter"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/astaxie/beego"
	"time"
)

type Response struct {
	MembershipsByDay           map[string]float64
	MembershipsByMonth         map[string]float64
	MembershipsByDayRnD        map[string]float64
	MembershipsByMonthRnD      map[string]float64
	MembershipCountsByMonth    map[string]int
	MembershipCountsByMonthRnD map[string]int
	ActivationsByDay           map[string]float64
	ActivationsByMonth         map[string]float64
	MinutesByDay               map[string]float64
	MinutesByMonth             map[string]float64
}

func NewResponse(data Data) (resp Response, err error) {
	resp.MembershipsByDay, err = data.sumMembershipsBy(bin.NewWidth(bin.DAY), false)
	if err != nil {
		return
	}
	resp.MembershipsByMonth, err = data.sumMembershipsBy(bin.NewWidth(bin.MONTH), false)
	if err != nil {
		return
	}
	resp.MembershipsByDayRnD, err = data.sumMembershipsBy(bin.NewWidth(bin.DAY), true)
	if err != nil {
		return
	}
	resp.MembershipsByMonthRnD, err = data.sumMembershipsBy(bin.NewWidth(bin.MONTH), true)
	if err != nil {
		return
	}
	resp.MembershipCountsByMonth, err = data.sumMembershipCountsBy(bin.NewWidth(bin.MONTH), false)
	if err != nil {
		return
	}
	resp.MembershipCountsByMonthRnD, err = data.sumMembershipCountsBy(bin.NewWidth(bin.MONTH), true)
	if err != nil {
		return
	}
	resp.ActivationsByDay, err = data.sumActivationsBy(bin.NewWidth(bin.DAY))
	if err != nil {
		return
	}
	resp.ActivationsByMonth, err = data.sumActivationsBy(bin.NewWidth(bin.MONTH))
	if err != nil {
		return
	}
	resp.MinutesByDay, err = data.sumMinutesBy(bin.NewWidth(bin.DAY))
	if err != nil {
		return
	}
	resp.MinutesByMonth, err = data.sumMinutesBy(bin.NewWidth(bin.MONTH))
	if err != nil {
		return
	}

	return
}

type Data struct {
	LocationId      int64
	startTime       time.Time
	endTime         time.Time
	Invoices        []*invutil.Invoice
	userLocations   user_locations.UserLocations
	userMemberships []*user_memberships.UserMembership
	membershipsById map[int64]*memberships.Membership
}

func FetchData(locationId int64, interval lib.Interval) (data Data, err error) {
	data.LocationId = locationId

	allInvoices, err := invutil.GetAllAt(locationId)
	if err != nil {
		return data, fmt.Errorf("Failed to get invoice summary: %v", err)
	}

	data.Invoices = filter.Invoices(allInvoices, interval.DayFrom(), interval.DayTo())

	ms, err := memberships.GetAllAt(locationId)
	if err != nil {
		return data, fmt.Errorf("Failed to get memberships: %v", err)
	}

	data.membershipsById = make(map[int64]*memberships.Membership)
	for _, m := range ms {
		data.membershipsById[m.Id] = m
	}

	data.userMemberships, err = user_memberships.GetAllAt(locationId)
	if err != nil {
		return data, fmt.Errorf("Failed to get user memberships: %v", err)
	}

	data.userLocations, err = user_locations.GetAllForLocation(locationId)
	if err != nil {
		return data, fmt.Errorf("Failed to get user locations: %v", err)
	}

	return
}

func (this Data) sumActivationsBy(w bin.Width) (sums map[string]float64, err error) {
	sums = make(map[string]float64)

	for _, inv := range this.Invoices {
		for _, purchase := range inv.Purchases {
			if purchase.Type == purchases.TYPE_ACTIVATION {
				priceTotalDisc, err := purchases.PriceTotalDisc(purchase)
				if err != nil {
					return nil, fmt.Errorf("PriceTotalDisc: %v", err)
				}
				var key string
				key = purchase.TimeStart.Format(w.TimeFormat())
				sums[key] = sums[key] + priceTotalDisc
			}
		}
	}

	return
}

func (this Data) sumMembershipsBy(w bin.Width, rndOnly bool) (sums map[string]float64, err error) {
	sums = make(map[string]float64)

	for _, inv := range this.Invoices {
		for _, userMembership := range inv.InvUserMemberships {
			membership, ok := this.membershipsById[userMembership.MembershipId]
			if !ok {
				return nil, fmt.Errorf("User Membership %v links to unknown Membership Id %v", userMembership.Id, userMembership.MembershipId)
			}

			if !rndOnly || membership.IsRndCentre() {
				t := time.Date(inv.Year, time.Month(inv.Month), 1, 12, 12, 12, 0, time.UTC)
				key := t.Format(w.TimeFormat())
				sums[key] = sums[key] + float64(membership.MonthlyPrice)
			}
		}
	}

	// Fab Lab Berlin (Location Id 1) was closed half of the December 2015.
	// At that time only half of the membership price was charged.
	if this.LocationId == 1 {
		d2015 := time.Date(2015, time.December, 1, 13, 0, 0, 0, time.UTC).
			Format(w.TimeFormat())
		if _, ok := sums[d2015]; ok {
			beego.Info("Dividing 12-2015 memberships by 2 @ locId 1")
			sums[d2015] = sums[d2015] / 2
		}
	}

	return
}

func (this Data) sumMembershipCountsBy(w bin.Width, rndOnly bool) (sums map[string]int, err error) {
	sums = make(map[string]int)

	for _, inv := range this.Invoices {
		for _, userMembership := range inv.InvUserMemberships {
			membership, ok := this.membershipsById[userMembership.MembershipId]
			if !ok {
				return nil, fmt.Errorf("User Membership %v links to unknown Membership Id %v", userMembership.Id, userMembership.MembershipId)
			}
			if !rndOnly || membership.IsRndCentre() {
				t := time.Date(inv.Year, time.Month(inv.Month), 1, 12, 12, 12, 0, time.UTC)
				key := t.Format(w.TimeFormat())
				if membership.MonthlyPrice > 0 {
					sums[key] = sums[key] + 1
				}
			}
		}
	}

	return
}

func (this Data) sumMinutesBy(w bin.Width) (sums map[string]float64, err error) {
	sums = make(map[string]float64)

	for _, inv := range this.Invoices {
		r, _ := this.userLocations.UserRoleOf(this.LocationId, inv.User.Id)
		if r != user_roles.STAFF && r != user_roles.ADMIN && !inv.User.SuperAdmin {
			for _, purchase := range inv.Purchases {
				if purchase.Type == purchases.TYPE_ACTIVATION {
					key := purchase.TimeStart.Format(w.TimeFormat())
					sums[key] = sums[key] + float64(purchase.Seconds())/60
				}
			}
		}
	}

	return
}
