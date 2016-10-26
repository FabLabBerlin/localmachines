package main

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego/migration"
	"github.com/astaxie/beego/orm"
	"sort"
	"time"
)

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
	StartDate    time.Time `orm:"type(datetime)"`
	EndDate      time.Time `orm:"type(datetime)"`
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
	InitialDurationMonths int
	AutoExtend            bool

	Created time.Time
	Updated time.Time
}

type Month struct {
	InvoiceId          int64
	Month              int
	Year               int
	OldUserMemberships []UserMembershipOld `orm:"-"`
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
type Months []Month

// General methods:
func (ms Months) MaxLenOldUserMemberships() (max int) {
	for _, m := range ms {
		if l := len(m.OldUserMemberships); l > max {
			max = l
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
	/*_, err = o.
	QueryTable("invoices").
	Filter("location_id", locId).
	Filter("user_id", userId).
	All(&ms)*/
	ivs, err := invoices.GetAllOfUserAt(locId, userId)
	if err != nil {
		return fmt.Errorf("GetAllOfUserAt: %v", err)
	}

	ms := make([]Month, 0, len(ivs))

	for _, iv := range ivs {
		ms = append(ms, Month{
			InvoiceId: iv.Id,
			Month:     iv.Month,
			Year:      iv.Year,
		})
	}

	months := Months(ms)

	sort.Sort(months)

	if l := months.MaxLenOldUserMemberships(); l > 1 {
		fmt.Printf("l=%v\n", l)
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
	orm.RegisterModel(new(UserMembershipOld))
	locs, err := locations.GetAll()
	if err != nil {
		panic(err.Error())
	}
	o := orm.NewOrm()
	o.Begin()
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
CREATE TABLE user_memberships (
	id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	location_id INT(11) UNSIGNED,
	membership_id INT(11) UNSIGNED,
	start_date DATETIME,
	termination_date DATETIME,
	intial_duration_months INT(11),
	auto_extend TINY INT(1),
	created DATETIME,
	updated DATETIME
)
	`)
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
