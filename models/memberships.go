package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Membership struct {
	Id                    int64  `orm:"auto";"pk"`
	Title                 string `orm:"size(100)"`
	ShortName             string `orm:"size(100)"`
	Duration              int
	Unit                  string `orm:"size(100)"`
	Price                 float32
	MachinePriceDeduction int
	AffectedMachines      string `orm:"type(text)"`
}

type UserMembership struct {
	Id           int64 `orm:"auto";"pk"`
	UserId       int64
	MembershipId int64
	StartDate    time.Time
}

func init() {
	orm.RegisterModel((new(Membership)), new(UserMembership))
}

// Get all memberships from database
func GetAllMemberships() (memberships []*Membership, err error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("membership").All(&memberships)
	if err != nil {
		beego.Error("Failed to get all memberships")
		return nil, errors.New("Failed to get all memberships")
	}
	beego.Trace("Got num memberships:", num)
	return
}

// Get single membership from database using membership unique ID
func GetMembership(membershipId int64) (*Membership, error) {
	membership := &Membership{Id: membershipId}
	o := orm.NewOrm()
	err := o.Read(membership)
	if err != nil {
		return membership, err
	}
	return membership, nil
}

func GetUserMemberships(userId int64) (ums []*UserMembership, err error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("user_membership").Filter("user_id", userId).All(&ums)
	if err != nil {
		beego.Error("Failed to get user memberships")
		return nil, errors.New("Failed to get user memberships")
	}
	beego.Trace("Got num user memberships:", num)
	return
}

// Creates a new membership in the database
func CreateMembership(membershipName string) (int64, error) {
	o := orm.NewOrm()
	membership := Membership{Title: membershipName}
	membership.Unit = "days"
	id, err := o.Insert(&membership)
	if err == nil {
		return id, nil
	} else {
		return 0, err
	}
}

// Updates membership in the database
func UpdateMembership(membership *Membership) error {
	var err error
	var num int64

	o := orm.NewOrm()
	num, err = o.Update(membership)
	if err != nil {
		return err
	}

	beego.Trace("Rows affected:", num)
	return nil
}

// Delete membership from the database
func DeleteMembership(membershipId int64) error {
	var num int64
	var err error
	o := orm.NewOrm()
	num, err = o.Delete(&Membership{Id: membershipId})
	if err != nil {
		return err
	}
	beego.Trace("Deleted num rows:", num)
	return nil
}
