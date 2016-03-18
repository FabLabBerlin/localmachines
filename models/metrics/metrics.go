/*
metrics package for basic visualization of numbers we have.
*/
package metrics

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
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
	startTime       time.Time
	endTime         time.Time
	invoice         invoices.Invoice
	userMemberships []*models.UserMembership
	membershipsById map[int64]*models.Membership
}

func FetchData(locationId int64) (data Data, err error) {
	endTime := time.Now()
	interval := lib.Interval{
		MonthFrom: int(time.August),
		YearFrom:  2015,
		MonthTo:   int(endTime.Month()),
		YearTo:    endTime.Year(),
	}

	data.invoice, err = invoices.New(locationId, interval)
	if err != nil {
		return data, fmt.Errorf("Failed to get invoice summary: %v", err)
	}

	memberships, err := models.GetAllMembershipsAt(locationId)
	if err != nil {
		return data, fmt.Errorf("Failed to get memberships: %v", err)
	}
	data.membershipsById = make(map[int64]*models.Membership)
	for _, membership := range memberships {
		data.membershipsById[membership.Id] = membership
	}

	data.userMemberships, err = models.GetAllUserMembershipsAt(locationId)
	if err != nil {
		return data, fmt.Errorf("Failed to get user memberships: %v", err)
	}

	return
}

func (this Data) sumActivationsBy(timeFormat string) (sums map[string]float64, err error) {
	sums = make(map[string]float64)

	for _, userSummary := range this.invoice.UserSummaries {
		for _, purchase := range userSummary.Purchases.Data {
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

	for _, userMembership := range this.userMemberships {
		membership, ok := this.membershipsById[userMembership.MembershipId]
		if !ok {
			return nil, fmt.Errorf("User Membership %v links to unknown Membership Id %v", userMembership.Id, userMembership.MembershipId)
		}
		for t := userMembership.StartDate; t.Before(userMembership.EndDate); t = t.AddDate(0, 1, 0) {
			key := t.Format(timeFormat)
			sums[key] = sums[key] + float64(membership.MonthlyPrice)
		}
	}

	return
}

func (this Data) sumMembershipCountsBy(timeFormat string) (sums map[string]int, err error) {
	sums = make(map[string]int)

	for _, userMembership := range this.userMemberships {
		membership, ok := this.membershipsById[userMembership.MembershipId]
		if !ok {
			return nil, fmt.Errorf("User Membership %v links to unknown Membership Id %v", userMembership.Id, userMembership.MembershipId)
		}
		for t := userMembership.StartDate; t.Before(userMembership.EndDate); t = t.AddDate(0, 1, 0) {
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

	for _, userSummary := range this.invoice.UserSummaries {
		if userSummary.User.GetRole() != user_roles.STAFF && userSummary.User.GetRole() != user_roles.ADMIN {
			for _, purchase := range userSummary.Purchases.Data {
				if purchase.Type == purchases.TYPE_ACTIVATION {
					key := purchase.TimeStart.Format(timeFormat)
					sums[key] = sums[key] + float64(purchase.Seconds())/60
				}
			}
		}
	}

	return
}
