/*
metrics package for basic visualization of numbers we have.
*/
package metrics

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/metrics/bin"
	"github.com/FabLabBerlin/localmachines/models/metrics/filter"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/astaxie/beego"
	"time"
)

type Response struct {
	Memberships         map[string]float64
	MembershipsRnD      map[string]float64
	MembershipCounts    map[string]int
	MembershipCountsRnD map[string]int
	Activations         map[string]float64
	Minutes             map[string]float64
}

func NewResponse(data Data) (resp Response, err error) {
	if data.BinWidth.IsMonth() {
		resp.Memberships, err = data.sumMembershipsBy(data.BinWidth, false)
		if err != nil {
			return
		}
		resp.MembershipsRnD, err = data.sumMembershipsBy(data.BinWidth, true)
		if err != nil {
			return
		}
		resp.MembershipCounts, err = data.sumMembershipCountsBy(data.BinWidth, false)
		if err != nil {
			return
		}
		resp.MembershipCountsRnD, err = data.sumMembershipCountsBy(data.BinWidth, true)
		if err != nil {
			return
		}
	} else {
		resp.Memberships = make(map[string]float64)
		resp.MembershipsRnD = make(map[string]float64)
		resp.MembershipCounts = make(map[string]int)
		resp.MembershipCountsRnD = make(map[string]int)
	}
	resp.Activations, err = data.sumActivationsBy(data.BinWidth)
	if err != nil {
		return
	}
	resp.Minutes, err = data.sumMinutesBy(data.BinWidth)
	if err != nil {
		return
	}

	return
}

type Data struct {
	LocationId int64
	Invoices   []*invutil.Invoice
	BinWidth   bin.Width
}

func FetchData(locationId int64, interval lib.Interval, binWidth bin.Width) (data Data, err error) {
	data.LocationId = locationId

	allInvoices, err := invutil.GetAllAt(locationId)
	if err != nil {
		return data, fmt.Errorf("Failed to get invoice summary: %v", err)
	}

	data.Invoices = filter.Invoices(
		allInvoices,
		interval.DayFrom().Month(),
		interval.DayTo().Month(),
	)

	userLocations, err := user_locations.GetAllForLocation(locationId)
	if err != nil {
		return data, fmt.Errorf("Failed to get user locations: %v", err)
	}

	data.Invoices = filter.InvoicesByUsers(
		data.LocationId,
		data.Invoices,
		userLocations,
		true,
		false,
	)

	data.BinWidth = binWidth

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
				key = w.TimeIndex(purchase.TimeStart)
				sums[key] = sums[key] + priceTotalDisc
			}
		}
	}

	return
}

func (this Data) sumMembershipsBy(w bin.Width, rndOnly bool) (sums map[string]float64, err error) {
	sums = make(map[string]float64)

	for _, inv := range this.Invoices {
		for _, ium := range inv.InvUserMemberships {
			if !rndOnly || ium.Membership().IsRndCentre() {
				t := time.Date(inv.Year, time.Month(inv.Month), 1, 12, 12, 12, 0, time.UTC)
				key := w.TimeIndex(t)
				sums[key] = sums[key] + float64(ium.Membership().MonthlyPrice)
			}
		}
	}

	// Fab Lab Berlin (Location Id 1) was closed half of the December 2015.
	// At that time only half of the membership price was charged.
	if this.LocationId == 1 {
		d2015 := w.TimeIndex(time.Date(2015, time.December, 1, 13, 0, 0, 0, time.UTC))
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
		for _, ium := range inv.InvUserMemberships {
			if !rndOnly || ium.Membership().IsRndCentre() {
				t := time.Date(inv.Year, time.Month(inv.Month), 1, 12, 12, 12, 0, time.UTC)
				key := w.TimeIndex(t)
				if ium.Membership().MonthlyPrice > 0 {
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
		for _, purchase := range inv.Purchases {
			if purchase.Type == purchases.TYPE_ACTIVATION {
				key := w.TimeIndex(purchase.TimeStart)
				sums[key] = sums[key] + float64(purchase.Seconds())/60
			}
		}
	}

	return
}
