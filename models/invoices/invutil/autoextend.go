package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/locations"
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
	invsByUid, us, err := taskData(locId)
	if err != nil {
		return fmt.Errorf("load task data: %v", err)
	}

	data := NewPrefetchedData(locId)

	if err := data.Prefetch(); err != nil {
		return fmt.Errorf("prefetch data: %v", err)
	}

	for _, u := range us {
		err := taskAutoExtendFor(locId, u, invsByUid, data)

		if err != nil {
			return fmt.Errorf("auto extend for user %v: %v", u.Id, err)
		}
	}

	return
}

func taskAutoExtendFor(
	locId int64,
	user *users.User,
	invsByUid map[int64][]*invoices.Invoice,
	data *PrefetchedData,
) (err error) {
	var currentIv *invoices.Invoice

	for _, iv := range invsByUid[user.Id] {
		if iv.Month == int(time.Now().Month()) &&
			iv.Year == time.Now().Year() &&
			iv.Status == "draft" {

			currentIv = iv
			break
		}
	}

	if currentIv == nil {
		currentIv, err = invoices.GetDraft(locId, user.Id, time.Now())
		if err != nil {
			return fmt.Errorf("get draft: %v", err)
		}
	}

	inv := Invoice{Invoice: *currentIv}

	if err = inv.load(*data); err != nil {
		return fmt.Errorf("load: %v", err)
	}

	if err := inv.InvoiceUserMemberships(data); err != nil {
		return fmt.Errorf("invoice user memberships: %v", err)
	}

	return
}

func taskData(locId int64) (
	invsByUid map[int64][]*invoices.Invoice,
	us []*users.User,
	err error,
) {
	invs, err := invoices.GetAllInvoices(locId)
	if err != nil {
		return nil, nil, fmt.Errorf("get all invoices: %v", err)
	}

	invsByUid = make(map[int64][]*invoices.Invoice)
	for _, inv := range invs {
		if _, ok := invsByUid[inv.UserId]; !ok {
			invsByUid[inv.UserId] = make([]*invoices.Invoice, 0, 10)
		}
		invsByUid[inv.UserId] = append(invsByUid[inv.UserId], inv)
	}

	us, err = users.GetAllUsersAt(locId)
	if err != nil {
		return nil, nil, fmt.Errorf("get all users: %v", err)
	}

	return
}
