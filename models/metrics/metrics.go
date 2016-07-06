/*
metrics package for basic visualization of numbers we have.
*/
package metrics

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/astaxie/beego"
	"time"
)

type Response struct {
	MembershipsByDay        map[string]float64
	MembershipsByMonth      map[string]float64
	MembershipCountsByMonth map[string]int
	ActivationsByDay        map[string]float64
	ActivationsByMonth      map[string]float64
	MinutesByDay            map[string]float64
	MinutesByMonth          map[string]float64
}

func NewResponse(data Data) (resp Response, err error) {
	resp.MembershipsByDay, err = data.sumMembershipsBy("2006-01-02")
	if err != nil {
		return
	}
	resp.MembershipsByMonth, err = data.sumMembershipsBy("2006-01")
	if err != nil {
		return
	}
	resp.MembershipCountsByMonth, err = data.sumMembershipCountsBy("2006-01")
	if err != nil {
		return
	}
	resp.ActivationsByDay, err = data.sumActivationsBy("2006-01-02")
	if err != nil {
		return
	}
	resp.ActivationsByMonth, err = data.sumActivationsBy("2006-01")
	if err != nil {
		return
	}
	resp.MinutesByDay, err = data.sumMinutesBy("2006-01-02")
	if err != nil {
		return
	}
	resp.MinutesByMonth, err = data.sumMinutesBy("2006-01")
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
	userMemberships []*user_memberships.UserMembership
	membershipsById map[int64]*memberships.Membership
}

func FetchData(locationId int64) (data Data, err error) {
	data.LocationId = locationId

	allInvoices, err := invutil.GetAllAt(locationId)
	if err != nil {
		return data, fmt.Errorf("Failed to get invoice summary: %v", err)
	}

	data.Invoices = filter(allInvoices)

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

	return
}

func filter(all []*invutil.Invoice) (filtered []*invutil.Invoice) {
	byUserIdYearMonth := make(map[int64]map[int]map[time.Month]*invutil.Invoice)
	filtered = make([]*invutil.Invoice, 0, len(all))

	for _, iv := range all {
		uid := iv.UserId
		m := time.Month(iv.Month)
		y := iv.Year

		if _, ok := byUserIdYearMonth[uid]; !ok {
			byUserIdYearMonth[uid] = make(map[int]map[time.Month]*invutil.Invoice)
		}
		if _, ok := byUserIdYearMonth[uid][y]; !ok {
			byUserIdYearMonth[uid][y] = make(map[time.Month]*invutil.Invoice)
		}
		if existing, ok := byUserIdYearMonth[uid][y][m]; !ok || iv.Id > existing.Id {
			byUserIdYearMonth[uid][y][m] = iv
		}
	}

	for _, byUid := range byUserIdYearMonth {
		for _, byYear := range byUid {
			for _, iv := range byYear {
				filtered = append(filtered, iv)
			}
		}
	}

	return
}

func (this Data) sumActivationsBy(timeFormat string) (sums map[string]float64, err error) {
	sums = make(map[string]float64)

	for _, inv := range this.Invoices {
		for _, purchase := range inv.Purchases {
			if purchase.Type == purchases.TYPE_ACTIVATION {
				priceTotalDisc, err := purchases.PriceTotalDisc(purchase)
				if err != nil {
					return nil, fmt.Errorf("PriceTotalDisc: %v", err)
				}
				var key string
				key = purchase.TimeStart.Format(timeFormat)
				sums[key] = sums[key] + priceTotalDisc
			}
		}
	}

	return
}

func (this Data) sumMembershipsBy(timeFormat string) (sums map[string]float64, err error) {
	sums = make(map[string]float64)

	for _, inv := range this.Invoices {
		for _, userMembership := range inv.UserMemberships.Data {
			membership, ok := this.membershipsById[userMembership.MembershipId]
			if !ok {
				return nil, fmt.Errorf("User Membership %v links to unknown Membership Id %v", userMembership.Id, userMembership.MembershipId)
			}

			t := time.Date(inv.Year, time.Month(inv.Month), 1, 12, 12, 12, 0, time.UTC)
			key := t.Format(timeFormat)
			sums[key] = sums[key] + float64(membership.MonthlyPrice)
		}
	}

	// Fab Lab Berlin (Location Id 1) was closed half of the December 2015.
	// At that time only half of the membership price was charged.
	if this.LocationId == 1 {
		d2015 := time.Date(2015, time.December, 1, 13, 0, 0, 0, time.UTC).
			Format(timeFormat)
		if _, ok := sums[d2015]; ok {
			beego.Info("Dividing 12-2015 memberships by 2 @ locId 1")
			sums[d2015] = sums[d2015] / 2
		}
	}

	return
}

func (this Data) sumMembershipCountsBy(timeFormat string) (sums map[string]int, err error) {
	sums = make(map[string]int)

	for _, inv := range this.Invoices {
		for _, userMembership := range inv.UserMemberships.Data {
			membership, ok := this.membershipsById[userMembership.MembershipId]
			if !ok {
				return nil, fmt.Errorf("User Membership %v links to unknown Membership Id %v", userMembership.Id, userMembership.MembershipId)
			}
			t := time.Date(inv.Year, time.Month(inv.Month), 1, 12, 12, 12, 0, time.UTC)
			key := t.Format(timeFormat)
			if membership.MonthlyPrice > 0 {
				sums[key] = sums[key] + 1
			}
		}
	}

	return
}

func (this Data) sumMinutesBy(timeFormat string) (sums map[string]float64, err error) {
	sums = make(map[string]float64)

	for _, inv := range this.Invoices {
		if inv.User.GetRole() != user_roles.STAFF && inv.User.GetRole() != user_roles.ADMIN {
			for _, purchase := range inv.Purchases {
				if purchase.Type == purchases.TYPE_ACTIVATION {
					key := purchase.TimeStart.Format(timeFormat)
					sums[key] = sums[key] + float64(purchase.Seconds())/60
				}
			}
		}
	}

	return
}
