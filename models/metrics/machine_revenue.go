package metrics

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/invoices/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"time"
)

const TROTEC_LOCATION_ID = 1

const TROTEC_ID = 17

const (
	PLUS_M         = 6
	PLUS_M_B2B     = 7
	PLUS_M_STUDENT = 9
	PLUS_M_PAO     = 12
)

type Trotec struct {
	TrotecMinutesByMonth               map[string]float64
	TrotecPaygMinutesByMonth           map[string]float64
	TrotecPlusMembershipMinutesByMonth map[string]float64
	AllPlusMembershipMinutesByMonth    map[string]float64

	TrotecReservationsEurByMonth map[string]float64

	TrotecPaygEurByMonth                       map[string]float64
	TrotecPlusMembershipUndiscountedEurByMonth map[string]float64
	AllPlusMembershipUndiscountedEurByMonth    map[string]float64

	AllPlusMembershipsEurByMonth map[string]float64
}

func NewTrotecStats() (t *Trotec, err error) {
	t = &Trotec{
		TrotecMinutesByMonth:               make(map[string]float64),
		TrotecPaygMinutesByMonth:           make(map[string]float64),
		TrotecPlusMembershipMinutesByMonth: make(map[string]float64),
		AllPlusMembershipMinutesByMonth:    make(map[string]float64),

		TrotecReservationsEurByMonth: make(map[string]float64),

		TrotecPaygEurByMonth:                       make(map[string]float64),
		TrotecPlusMembershipUndiscountedEurByMonth: make(map[string]float64),
		AllPlusMembershipUndiscountedEurByMonth:    make(map[string]float64),

		AllPlusMembershipsEurByMonth: make(map[string]float64),
	}
	endTime := time.Now()
	interval := lib.Interval{
		MonthFrom: int(time.November),
		YearFrom:  2015,
		MonthTo:   int(endTime.Month()),
		YearTo:    endTime.Year(),
	}
	monthlyEarning, err := monthly_earning.New(TROTEC_LOCATION_ID, interval)
	if err != nil {
		return nil, fmt.Errorf("Failed to get invoice summary: %v", err)
	}
	for _, inv := range monthlyEarning.Invoices {
		for _, purchase := range inv.Purchases {
			if hasFreeMembership(purchase.Memberships) {
				continue
			}
			isStaff := false
			switch purchase.User.GetRole().String() {
			case user_roles.STAFF.String(), user_roles.ADMIN.String(), user_roles.SUPER_ADMIN.String():
				isStaff = true
			}
			uls, err := user_locations.GetAllForUser(purchase.UserId)
			for _, ul := range uls {
				if ul.LocationId == TROTEC_LOCATION_ID {
					switch ul.GetRole().String() {
					case user_roles.STAFF.String(), user_roles.ADMIN.String(), user_roles.SUPER_ADMIN.String():
						isStaff = true
					}
				}
			}
			if isStaff {
				continue
			}
			priceTotalDisc, err := purchases.PriceTotalDisc(purchase)
			if err != nil {
				return nil, fmt.Errorf("PriceTotalDisc: %v", err)
			}
			priceUndiscounted := purchases.PriceTotalExclDisc(purchase)
			month := purchase.TimeStart.Month().String()
			if purchase.Type == purchases.TYPE_ACTIVATION {
				if purchase.MachineId == TROTEC_ID {
					t.TrotecMinutesByMonth[month] = t.TrotecMinutesByMonth[month] + purchase.Quantity
					if hasTrotecRebate(purchase.Memberships) {
						t.TrotecPlusMembershipMinutesByMonth[month] = t.TrotecPlusMembershipMinutesByMonth[month] + purchase.Quantity
						t.TrotecPlusMembershipUndiscountedEurByMonth[month] = t.TrotecPlusMembershipUndiscountedEurByMonth[month] + priceUndiscounted
					} else {
						t.TrotecPaygMinutesByMonth[month] = t.TrotecPaygMinutesByMonth[month] + purchase.Quantity
						t.TrotecPaygEurByMonth[month] = t.TrotecPaygEurByMonth[month] + priceTotalDisc
					}
				}
				if hasTrotecRebate(purchase.Memberships) {
					t.AllPlusMembershipMinutesByMonth[month] = t.AllPlusMembershipMinutesByMonth[month] + purchase.Quantity
					t.AllPlusMembershipUndiscountedEurByMonth[month] = t.AllPlusMembershipUndiscountedEurByMonth[month] + priceUndiscounted
				}
			} else if purchase.Type == purchases.TYPE_RESERVATION {
				if purchase.MachineId == TROTEC_ID {
					t.TrotecReservationsEurByMonth[month] = t.TrotecReservationsEurByMonth[month] + priceUndiscounted
				}
			}
		}
	}

	ms, err := memberships.GetAllMembershipsAt(TROTEC_LOCATION_ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get memberships: %v", err)
	}
	membershipsById := make(map[int64]*memberships.Membership)
	for _, m := range ms {
		membershipsById[m.Id] = m
	}

	userMemberships, err := user_memberships.GetAllUserMembershipsAt(TROTEC_LOCATION_ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user memberships: %v", err)
	}

	midMonth := time.Date(2015, time.November, 15, 5, 5, 5, 5, time.UTC)
	for ; midMonth.Month() != time.April; midMonth = midMonth.AddDate(0, 1, 0) {
		month := midMonth.Month()
		for _, um := range userMemberships {
			if um.Interval().Contains(midMonth) && membershipIdHasTrotecRebate(um.MembershipId) {
				m := membershipsById[um.MembershipId]
				t.AllPlusMembershipsEurByMonth[month.String()] = t.AllPlusMembershipsEurByMonth[month.String()] + m.MonthlyPrice
			}
		}
	}

	return
}

func hasTrotecRebate(ms []*memberships.Membership) bool {
	for _, m := range ms {
		if membershipIdHasTrotecRebate(m.Id) {
			return true
		}
	}
	return false
}

func membershipIdHasTrotecRebate(id int64) bool {
	switch id {
	case PLUS_M, PLUS_M_B2B, PLUS_M_STUDENT, PLUS_M_PAO:
		return true
	}
	return false
}

func hasFreeMembership(ms []*memberships.Membership) bool {
	for _, m := range ms {
		if m.MonthlyPrice < 1 {
			return true
		}
	}
	return false
}
