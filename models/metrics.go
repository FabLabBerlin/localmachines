package models

import (
	"fmt"
	"time"
)

type MetricsResponse struct {
	MembershipsByDay   map[string]float64
	MembershipsByMonth map[string]float64
	ActivationsByDay   map[string]float64
	ActivationsByMonth map[string]float64
}

func NewMetricsResponse(data MetricsData) (resp MetricsResponse, err error) {
	resp.MembershipsByDay, err = data.sumMembershipsBy("2006-01-02")
	if err != nil {
		return
	}
	resp.MembershipsByMonth, err = data.sumMembershipsBy("2006-01")
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

	return
}

type MetricsData struct {
	startTime       time.Time
	endTime         time.Time
	invoice         Invoice
	userMemberships []*UserMembership
	membershipsById map[int64]*Membership
}

func FetchMetricsData() (data MetricsData, err error) {
	endTime := time.Now()
	startTime := endTime.Add(-180 * 24 * time.Hour)

	data.invoice, err = CalculateInvoiceSummary(startTime, endTime)
	if err != nil {
		return data, fmt.Errorf("Failed to get invoice summary: %v", err)
	}

	memberships, err := GetAllMemberships()
	if err != nil {
		return data, fmt.Errorf("Failed to get memberships: %v", err)
	}
	data.membershipsById = make(map[int64]*Membership)
	for _, membership := range memberships {
		data.membershipsById[membership.Id] = membership
	}

	data.userMemberships, err = GetAllUserMemberships()
	if err != nil {
		return data, fmt.Errorf("Failed to get user memberships: %v", err)
	}

	return
}

func (this MetricsData) sumActivationsBy(timeFormat string) (sums map[string]float64, err error) {
	sums = make(map[string]float64)

	for _, userSummary := range this.invoice.UserSummaries {
		for _, purchase := range userSummary.Purchases.Data {
			priceTotalDisc, err := PriceTotalDisc(purchase)
			if err != nil {
				return nil, fmt.Errorf("PriceTotalDisc: %v", err)
			}
			var key string
			if purchase.Activation != nil {
				key = purchase.Activation.TimeStart.Format(timeFormat)
			} else {
				key = purchase.Reservation.TimeStart.Format(timeFormat)
			}
			sums[key] = sums[key] + priceTotalDisc
		}
	}

	return
}

func (this MetricsData) sumMembershipsBy(timeFormat string) (sums map[string]float64, err error) {
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
