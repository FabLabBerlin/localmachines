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
type InvoiceUserMemberships_20161025_190752 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &InvoiceUserMemberships_20161025_190752{}
	m.Created = "20161025_190752"
	migration.Register("InvoiceUserMemberships_20161025_190752", m)
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
}

func (this *UserMembershipNew) TableName() string {
	return "user_memberships"
}

type Month struct {
	LocationId         int64
	InvoiceId          int64
	Month              int
	Year               int
	OldUserMemberships []*UserMembershipOld `orm:"-"`
}

func (m Month) OldUserMembershipsReversed() (ums []*UserMembershipOld) {
	ums = make([]*UserMembershipOld, len(m.OldUserMemberships))

	for i, um := range m.OldUserMemberships {
		ums[len(ums)-1-i] = um
	}

	return
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
		panic("not implemented yet")
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
	if ms.MaxLenOldUserMemberships() > 1 {
		panic("not implemented yet")
	} else if ms.HasVaryingMemberships() {
		panic("not implemented yet")
	}

	// Case for l = 1:
	um := &UserMembershipNew{}
	for _, m := range ms {
		for _, old := range m.OldUserMembershipsReversed() {
			um.LocationId = m.LocationId
			um.UserId = old.UserId
			um.MembershipId = old.MembershipId
			um.StartDate = old.StartDate
			if old.EndDate.Before(time.Now()) {
				um.TerminationDate = old.EndDate
			}
			um.AutoExtend = old.AutoExtend
			um.Updated = time.Now()
			um.InitialDurationMonths = old.Membership.DurationMonths

			return []*UserMembershipNew{um}
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
		fmt.Printf("l=%v\n", l)
	} else if months.HasVaryingMemberships() {
		fmt.Printf("varying\n")
	} else {
		for _, newUm := range months.NewUserMemberships() {
			if _, err = o.Insert(newUm); err != nil {
				return fmt.Errorf("insert new um: %v", err)
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

// Run the migrations
func (m *InvoiceUserMemberships_20161025_190752) Up() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovering from %v...\n", r)
			orm.NewOrm().Raw("DROP TABLE user_memberships").Exec()
			panic("lol")
		} else {
			fmt.Printf("no recovery possible\n")
		}
	}()

	_, err := orm.NewOrm().Raw(`
CREATE TABLE user_memberships (
	id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	location_id INT(11) UNSIGNED,
	membership_id INT(11) UNSIGNED,
	start_date DATETIME,
	termination_date DATETIME,
	intial_duration_months INT(11),
	auto_extend TINYINT(1),
	created DATETIME,
	updated DATETIME,
	PRIMARY KEY (id)
)
	`).Exec()
	if err != nil {
		panic(err.Error())
	}

	orm.RegisterModel(new(UserMembershipOld), new(UserMembershipNew))

	o := orm.NewOrm()
	o.Begin()

	var mbs []*memberships.Membership
	_, err = o.QueryTable("membership").
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
	panic("123x")
	if err := o.Commit(); err != nil {
		panic(err.Error())
	}

	type User struct {
		Id       int
		Username string
	}

	var user User
	err = o.Raw("SELECT id, username FROM user WHERE id = ?", 19).QueryRow(&user)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("user=%v\n", user)
	panic("foo")
	m.SQL("RENAME TABLE user_membership TO invoice_user_memberships")
	m.SQL(`
INSERT INTO user_memberships (
	location_id
)
SELECT
	location_id
FROM invoice_user_memberships
WHERE
GROUP BY
	`)
}

// Reverse the migrations
func (m *InvoiceUserMemberships_20161025_190752) Down() {
	m.SQL("RENAME TABLE invoice_user_memberships TO user_membership")
	m.SQL("DROP TABLE user_memberships")
}
