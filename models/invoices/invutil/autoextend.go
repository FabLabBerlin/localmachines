package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"time"
)

func TaskAutoExtend() (err error) {
	ls, err := locations.GetAll()
	if err != nil {
		return fmt.Errorf("get all locations: %v", err)
	}

	for _, l := range ls {
		if err := taskAutoExtend(l.Id); err != nil {
			return fmt.Errorf("auto extend @ location %v: %v", l.Id, err)
		}
	}

	return
}

func taskAutoExtend(locId int64) (err error) {
	umsByUid, invsByUid, us, err := taskData(locId)
	if err != nil {
		return fmt.Errorf("load task data: %v", err)
	}

	data := NewPrefetchedData(locId)

	if err := data.Prefetch(); err != nil {
		return fmt.Errorf("prefetch data: %v", err)
	}

	for _, u := range us {
		iv, ok := invsByUid[u.Id]
		if !ok {
			iv, err = invoices.GetDraft(locId, u.Id, time.Now())
			if err != nil {
				return fmt.Errorf("get draft: %v", err)
			}
		}

		inv, err := Get(iv.Id)
		if err != nil {
			return fmt.Errorf("invutil.Get(%v): %v", iv.Id, err)
		}

		if err := inv.InvoiceUserMemberships(data); err != nil {
			return fmt.Errorf("invoice user memberships for %v: %v", u.Id, err)
		}
	}

	return
}

func taskData(locId int64) (
	umsByUid map[int64]*user_memberships.UserMembership,
	invsByUid map[int64][]*invoices.Invoice,
	us []*users.User,
	err error,
) {
	ums, err := user_memberships.GetAllAt(locId)
	if err != nil {
		return fmt.Errorf("get all user memberships: %v", err)
	}

	umsByUid = make(map[int64]*user_memberships.UserMembership)
	for _, um := range ums {
		umsByUid[um.UserId] = um
	}

	invs, err := invoices.GetAllAt(locId)
	if err != nil {
		return fmt.Errorf("get all invoices: %v", err)
	}

	invsByUid = make(map[int64][]*Invoice)
	for _, inv := range invs {
		if _, ok := invsByUid[inv.UserId]; !ok {
			invsByUid[inv.UserId] = make([]*Invoice, 0, 10)
		}
		invsByUid[inv.UserId] = append(invsByUid[inv.UserId], inv)
	}

	us, err = users.GetAllUsersAt(locId)
	if err != nil {
		return fmt.Errorf("get all users: %v", err)
	}

	return
}
