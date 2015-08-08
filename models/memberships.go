package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

func (this *Membership) TableName() string {
	return "membership"
}

type UserMembership struct {
	Id           int64 `orm:"auto";"pk"`
	UserId       int64
	MembershipId int64
	StartDate    time.Time
}

func (this *UserMembership) TableName() string {
	return "user_membership"
}

func init() {
	orm.RegisterModel((new(Membership)), new(UserMembership))
}

// Extended UserMembership type that contains the same fields as UserMembership
// model, but additionaly fields from the membership model.
type UserMembershipCombo struct {
	Id                    int64
	UserId                int64
	MembershipId          int64
	StartDate             time.Time
	Title                 string
	ShortName             string
	Duration              int
	Unit                  string
	Price                 float32
	MachinePriceDeduction int
	AffectedMachines      string
}

type UserMembershipList struct {
	Data []UserMembershipCombo
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

	// Check if the membership has been added to user memberships.
	// Do not allow deletion if so.
	umem := UserMembership{}
	num, err = o.QueryTable(umem.TableName()).
		Filter("membership_id", membershipId).
		Count()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to get user memberships: %v", err))
	}

	if num > 0 {
		return errors.New(
			fmt.Sprintf("Membership has been assigned to %d users", num))
	}

	num, err = o.Delete(&Membership{Id: membershipId})
	if err != nil {
		return err
	}
	beego.Trace("Deleted num rows:", num)
	return nil
}

func CreateUserMembership(userMembership *UserMembership) (umid int64, err error) {
	if userMembership != nil {
		o := orm.NewOrm()

		if umid, err := o.Insert(userMembership); err != nil {
			beego.Error("Error creating new user membership: ", err)
			return umid, err
		}

		return umid, nil
	} else {
		return 0, errors.New("userMembership is nil")
	}
}

func GetUserMemberships(userId int64) (*UserMembershipList, error) {

	o := orm.NewOrm()

	// Use these for the table names
	m := Membership{}
	um := UserMembership{}

	/*
			Id           int64 `orm:"auto";"pk"`
		UserId       int64
		MembershipId int64
		StartDate    time.Time
	*/

	// Joint query, select user memberships and expands them with
	// membership base information.
	sql := fmt.Sprintf("SELECT um.*, m.title, m.short_name, m.duration, m.unit, "+
		"m.price, m.machine_price_deduction, m.affected_machines "+
		"FROM %s AS um "+
		"JOIN %s m ON um.membership_id=m.id "+
		"WHERE um.user_id=?",
		um.TableName(), m.TableName())

	var userMemberships []UserMembershipCombo
	var num int64
	var err error
	num, err = o.Raw(sql, userId).QueryRows(&userMemberships)
	if err != nil {
		beego.Error(err)
		return nil, fmt.Errorf("Failed to get user memberships")
	}
	beego.Trace("Got num user memberships:", num)

	userMembershipList := UserMembershipList{}
	userMembershipList.Data = userMemberships

	return &userMembershipList, nil
}

// Delete membership from the database
func DeleteUserMembership(userMembershipId int64) error {
	var num int64
	var err error
	o := orm.NewOrm()

	beego.Trace("models.DeleteUserMembership: userMembershipId:",
		userMembershipId)

	userMembership := &UserMembership{}
	userMembership.Id = userMembershipId

	num, err = o.QueryTable(userMembership.TableName()).
		Filter("id", userMembershipId).Delete()
	if err != nil {
		beego.Error("Failed to delete user membership")
		return errors.New("Failed to delete user membership")
	}

	/*
		num, err = o.Delete(userMembership)
		if err != nil {
			return errors.New(
				fmt.Sprintf("Failed to delete user membership: %v", err))
		}
	*/

	beego.Trace("Deleted num rows:", num)
	return nil
}
