package main

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego/migration"
	"github.com/astaxie/beego/orm"
	"sort"
	"time"
)

var membershipsById = make(map[int64]*memberships.Membership)

/*
At the moment:

u_i: um_1, um_2, ..., um_n

In the end:

u_i: um_1, ..., um_m, m <= n
*/

// DO NOT MODIFY
type PopulateInvoiceUserMemberships_20161028_134609 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PopulateInvoiceUserMemberships_20161028_134609{}
	m.Created = "20161028_134609"
	migration.Register("PopulateInvoiceUserMemberships_20161028_134609", m)
}

type UserMembershipOld struct {
	Id           int64
	UserId       int64
	MembershipId int64
	Membership   *memberships.Membership `orm:"-"`
	StartDate    time.Time               `orm:"type(datetime)"`
	EndDate      time.Time               `orm:"type(datetime)"`
	AutoExtend   bool

	InvoiceId     int64
	InvoiceStatus string
}

func (this *UserMembershipOld) TableName() string {
	return "user_membership"
}

type UserMembershipNew struct {
	Id           int64
	LocationId   int64
	UserId       int64
	MembershipId int64
	//Membership            *memberships.Membership `orm:"-" json:",omitempty"`
	StartDate             time.Time `orm:"type(datetime)"`
	TerminationDate       time.Time `orm:"type(datetime)"`
	InitialDurationMonths int64
	AutoExtend            bool

	Created time.Time
	Updated time.Time

	InvUserMemberships []*InvoiceUserMembership `orm:"-"`
}

func (this *UserMembershipNew) TableName() string {
	return "user_memberships"
}

type InvoiceUserMembership struct {
	Id                    int64
	LocationId            int64
	UserId                int64
	MembershipId          int64
	UserMembershipId      int64
	StartDate             time.Time `orm:"type(datetime)"`
	TerminationDate       time.Time `orm:"type(datetime)"`
	InitialDurationMonths int64

	Created time.Time
	Updated time.Time

	InvoiceId     int64
	InvoiceStatus string
}

func NewInvoiceUserMembership(old UserMembershipOld) (ium *InvoiceUserMembership) {
	ium = &InvoiceUserMembership{
		LocationId:    old.Membership.LocationId,
		UserId:        old.UserId,
		MembershipId:  old.MembershipId,
		Updated:       time.Now(),
		InvoiceId:     old.InvoiceId,
		InvoiceStatus: old.InvoiceStatus,
	}

	return
}

func (this *InvoiceUserMembership) TableName() string {
	return "invoice_user_memberships"
}

type Month struct {
	LocationId         int64
	InvoiceId          int64
	Month              int
	Year               int
	OldUserMemberships []*UserMembershipOld `orm:"-"`
}

func (m Month) toInt() int {
	return m.Year*100 + m.Month
}

//       Months
//
//       Jan 2016
//          |
//          \/
//       Feb 2016
//          |
//          \/
//         ...
//
//
type Months []*Month

// General methods:
func (ms Months) HasVaryingMemberships() bool {
	l := ms.MaxLenOldUserMemberships()
	if l == 0 {
		return false
	} else if l > 1 {
		return true
	}

	var membershipId int64
	for _, m := range ms {
		if len(m.OldUserMemberships) > 0 {
			um := m.OldUserMemberships[0]

			if membershipId == 0 {
				membershipId = um.MembershipId
			} else if membershipId != um.MembershipId {
				return true
			}
		}
	}
	return false
}

func (ms Months) MaxLenOldUserMemberships() (max int) {
	for _, m := range ms {
		if l := len(m.OldUserMemberships); l > max {
			max = l
		}
	}
	return
}

func (ms Months) NewUserMemberships() (ums []*UserMembershipNew) {
	var lastMonthProcessed []*UserMembershipNew

	for _, m := range ms {
		processed := make([]*UserMembershipNew, 0, 2)

		for _, old := range m.OldUserMemberships {
			var um *UserMembershipNew

			for i, umm := range lastMonthProcessed {
				if old.MembershipId == umm.MembershipId &&
					old.StartDate.Equal(umm.StartDate) {

					um = lastMonthProcessed[i]
					lastMonthProcessed = append(lastMonthProcessed[:i], lastMonthProcessed[i+1:]...)
					break
				}
			}

			if um == nil {
				um = &UserMembershipNew{}
				um.LocationId = m.LocationId
				um.UserId = old.UserId
				um.MembershipId = old.MembershipId
				um.StartDate = old.StartDate
				um.Updated = time.Now()
				um.InitialDurationMonths = old.Membership.DurationMonths
				ums = append(ums, um)
			}

			dbg := func(msg string) {
				if um.UserId == 179 {
					fmt.Printf("%v\n", msg)
				}
			}

			if old.EndDate.Before(time.Now()) {
				dbg(fmt.Sprintf("a\n"))
				um.TerminationDate = old.EndDate
			} else if m.Year < time.Now().Year() || m.Month < int(time.Now().Month()) {
				dbg(fmt.Sprintf("b, %v/%v\n", m.Month, m.Year))
				um.TerminationDate = time.Date(m.Year, time.Month(m.Month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
			} else {
				dbg(fmt.Sprintf("c\n"))
				um.AutoExtend = old.AutoExtend
				um.TerminationDate = time.Time{}
			}

			// Create corresponding invoice user membership
			ium := NewInvoiceUserMembership(*old)
			um.InvUserMemberships = append(um.InvUserMemberships, ium)
			ium.StartDate = old.StartDate
			//if old.EndDate.Before(time.Now()) {
			//	ium.TerminationDate = old.EndDate
			//}
			ium.TerminationDate = um.TerminationDate
			ium.InitialDurationMonths = um.InitialDurationMonths

			processed = append(processed, um)
		}

		lastMonthProcessed = processed
	}

	for _, um := range ums {
		for i := len(um.InvUserMemberships) - 1; i >= 0; i-- {
			ium := um.InvUserMemberships[i]

			ium.TerminationDate = um.TerminationDate
		}
	}

	return
}

// Implementation of sort.Interface:
func (ms Months) Len() int {
	return len(ms)
}

func (ms Months) Less(i, j int) bool {
	return ms[i].toInt() < ms[j].toInt()
}

func (ms Months) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

type UserId int64

func userUp(o orm.Ormer, locId, userId int64) (err error) {
	ivs, err := invoices.GetAllOfUserAt(locId, userId)
	if err != nil {
		return fmt.Errorf("GetAllOfUserAt: %v", err)
	}

	ms := make([]*Month, 0, len(ivs))

	for _, iv := range ivs {
		ms = append(ms, &Month{
			LocationId: locId,
			InvoiceId:  iv.Id,
			Month:      iv.Month,
			Year:       iv.Year,
		})
	}

	months := Months(ms)

	sort.Sort(months)

	var ums []*UserMembershipOld

	if _, err = o.
		QueryTable("user_membership").
		Filter("user_id", userId).
		All(&ums); err != nil {

		return fmt.Errorf("query user_membership: %v", err)
	}

	for _, um := range ums {
		um.Membership = membershipsById[um.MembershipId]

		for _, m := range months {
			if um.InvoiceId == m.InvoiceId {
				m.OldUserMemberships = append(m.OldUserMemberships, um)
				break
			}
		}
	}

	if l := months.MaxLenOldUserMemberships(); l > 1 {
		fmt.Printf("l=%v for uid=%v\n", l, userId)
	}
	if months.HasVaryingMemberships() {
		fmt.Printf("varying for uid=%v\n", userId)
	}
	for _, newUm := range months.NewUserMemberships() {
		if _, err = o.Insert(newUm); err != nil {
			return fmt.Errorf("insert new um: %v", err)
		}
		for _, ium := range newUm.InvUserMemberships {
			ium.UserMembershipId = newUm.Id
			if _, err = o.Insert(ium); err != nil {
				return fmt.Errorf("insert ium: %v", err)
			}
		}
	}

	return
}

func locationUp(o orm.Ormer, locId int64) (err error) {
	us, err := users.GetAllUsersAt(locId)
	if err != nil {
		return fmt.Errorf("get all at %v: %v", locId, err)
	}

	for _, u := range us {
		if err := userUp(o, locId, u.Id); err != nil {
			return fmt.Errorf("user %v up: %v", u.Id, err)
		}
	}

	return
}

func init() {
	orm.RegisterModel(new(UserMembershipOld), new(UserMembershipNew), new(InvoiceUserMembership))
}

// Run the migrations
func (m *PopulateInvoiceUserMemberships_20161028_134609) Up() {
	o := orm.NewOrm()
	o.Begin()

	var mbs []*memberships.Membership
	_, err := o.QueryTable("membership").
		All(&mbs)

	for _, mb := range mbs {
		membershipsById[mb.Id] = mb
	}

	locs, err := locations.GetAll()
	if err != nil {
		panic(err.Error())
	}

	for _, loc := range locs {
		if err := locationUp(o, loc.Id); err != nil {
			o.Rollback()
			panic(err.Error())
		}
	}

	if err := o.Commit(); err != nil {
		panic(err.Error())
	}
}

// Reverse the migrations
func (m *PopulateInvoiceUserMemberships_20161028_134609) Down() {
	m.SQL("DELETE FROM user_memberships")
	m.SQL("DELETE FROM invoice_user_memberships")
}
