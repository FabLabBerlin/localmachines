// filter data
package filter

import (
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
)

func Invoices(all []*invutil.Invoice, from, to month.Month) (ivs []*invutil.Invoice) {
	ivs = make([]*invutil.Invoice, 0, len(all))

	for _, iv := range all {
		if iv.GetMonth().Before(from) {
			continue
		}

		if iv.GetMonth().After(to) {
			continue
		}

		if !iv.Canceled {
			ivs = append(ivs, iv)
		}
	}

	return
}

func InvoicesByUsers(
	locId int64,
	all []*invutil.Invoice,
	userLocations []*user_locations.UserLocation,
	excludeStaff bool,
	excludeNeverActive bool,
) (invs []*invutil.Invoice) {

	invs = make([]*invutil.Invoice, 0, len(all))

	rolesByUid := make(map[int64]user_roles.Role)
	for _, ul := range userLocations {
		if ul.LocationId == locId {
			rolesByUid[ul.UserId] = ul.GetRole()
		}
	}

	everActiveUid := make(map[int64]struct{})

	everStaffUid := make(map[int64]struct{})

	for _, inv := range all {
		for _, p := range inv.Purchases {
			everActiveUid[p.UserId] = struct{}{}
		}
		for _, ium := range inv.InvUserMemberships {
			everActiveUid[ium.UserId] = struct{}{}
			if ium.MembershipId == 3 {
				everStaffUid[ium.UserId] = struct{}{}
			}
		}
	}

	for _, inv := range all {
		u := inv.User

		if excludeNeverActive {
			if _, everActive := everActiveUid[u.Id]; !everActive {
				continue
			}
		}
		if excludeStaff {
			if u.SuperAdmin {
				continue
			}

			if _, everStaff := everStaffUid[u.Id]; everStaff {
				continue
			}
		}
		if r, ok := rolesByUid[u.Id]; (!ok || r == user_roles.MEMBER) || !excludeStaff {
			invs = append(invs, inv)
		}
	}

	return
}
