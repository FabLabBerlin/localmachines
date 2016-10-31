package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships/inv_user_memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
)

type PrefetchedData struct {
	LocationId              int64
	UsersById               map[int64]*users.User
	PurchasesByInv          map[int64][]*purchases.Purchase
	InvUserMembershipsByInv map[int64][]*inv_user_memberships.InvoiceUserMembership
	MsById                  map[int64]*machine.Machine
	MbsById                 map[int64]*memberships.Membership
	UmbsByUid               map[int64][]*user_memberships.UserMembership
	IumbsByUid              map[int64][]*inv_user_memberships.InvoiceUserMembership
}

func NewPrefetchedData(locId int64) (data *PrefetchedData) {
	data = &PrefetchedData{
		LocationId:              locId,
		UsersById:               make(map[int64]*users.User),
		PurchasesByInv:          make(map[int64][]*purchases.Purchase),
		InvUserMembershipsByInv: make(map[int64][]*inv_user_memberships.InvoiceUserMembership),
		MsById:                  make(map[int64]*machine.Machine),
		MbsById:                 make(map[int64]*memberships.Membership),
		UmbsByUid:               make(map[int64][]*user_memberships.UserMembership),
		IumbsByUid:              make(map[int64][]*inv_user_memberships.InvoiceUserMembership),
	}

	return
}

func (data *PrefetchedData) Prefetch() (err error) {
	if us, err := users.GetAllUsersAt(data.LocationId); err == nil {
		for _, u := range us {
			data.UsersById[u.Id] = u
		}
	} else {
		return fmt.Errorf("get all users: %v", err)
	}

	if ps, err := purchases.GetAllAt(data.LocationId); err == nil {
		for _, p := range ps {
			if _, ok := data.PurchasesByInv[p.InvoiceId]; !ok {
				data.PurchasesByInv[p.InvoiceId] = make([]*purchases.Purchase, 0, 20)
			}
			data.PurchasesByInv[p.InvoiceId] = append(data.PurchasesByInv[p.InvoiceId], p)
		}
	} else {
		return fmt.Errorf("get all purchases: %v", err)
	}

	if umbs, err := inv_user_memberships.GetAllAt(data.LocationId); err == nil {
		for _, umb := range umbs {
			if _, ok := data.InvUserMembershipsByInv[umb.InvoiceId]; !ok {
				data.InvUserMembershipsByInv[umb.InvoiceId] = make([]*inv_user_memberships.InvoiceUserMembership, 0, 3)
			}
			data.InvUserMembershipsByInv[umb.InvoiceId] = append(data.InvUserMembershipsByInv[umb.InvoiceId], umb)
		}
	} else {
		return fmt.Errorf("get all user memberships: %v", err)
	}

	ms, err := machine.GetAllAt(data.LocationId)
	if err != nil {
		return fmt.Errorf("get all machines at %v: %v", data.LocationId, err)
	}

	for _, m := range ms {
		data.MsById[m.Id] = m
	}

	if mbs, err := memberships.GetAllAt(data.LocationId); err == nil {
		for _, mb := range mbs {
			data.MbsById[mb.Id] = mb
		}
	} else {
		return fmt.Errorf("Failed to get memberships: %v", err)
	}

	if umbs, err := user_memberships.GetAllAt(data.LocationId); err == nil {
		for _, umb := range umbs {
			uid := umb.UserId
			umb.Membership = data.MbsById[umb.MembershipId]
			if _, ok := data.UmbsByUid[uid]; !ok {
				data.UmbsByUid[uid] = []*user_memberships.UserMembership{
					umb,
				}
			} else {
				data.UmbsByUid[uid] = append(data.UmbsByUid[uid], umb)
			}
		}
	} else {
		return fmt.Errorf("Failed to get user memberships: %v", err)
	}

	if iumbs, err := inv_user_memberships.GetAllAt(data.LocationId); err == nil {
		for _, iumb := range iumbs {
			uid := iumb.UserId
			if _, ok := data.IumbsByUid[uid]; !ok {
				data.IumbsByUid[uid] = []*inv_user_memberships.InvoiceUserMembership{
					iumb,
				}
			} else {
				data.IumbsByUid[uid] = append(data.IumbsByUid[uid], iumb)
			}
		}
	} else {
		return fmt.Errorf("Failed to get user memberships: %v", err)
	}

	return
}
