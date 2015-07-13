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

type MembershipResponse struct {
	Id                    int64  `orm:"auto";"pk"`
	Title                 string `orm:"size(100)"`
	ShortName             string `orm:"size(100)"`
	Duration              int
	Unit                  string `orm:"size(100)"`
	Price                 float32
	MachinePriceDeduction int
	AffectedMachines      string `orm:"type(text)"`
	StartDate             time.Time
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

func GetUserMemberships(userId int64) (membs *[]MembershipResponse, err error) {
	o := orm.NewOrm()
	membershipsResponse := []MembershipResponse{}
	mem := Membership{}
	usr_mem := UserMembership{}

	query := fmt.Sprintf("SELECT m.*, um.start_date FROM %s m JOIN %s um ON um.membership_id=m.id "+
		"WHERE um.user_id=?",
		mem.TableName(),
		usr_mem.TableName())

	beego.Trace(query)

	num, err := o.Raw(query, userId).QueryRows(&membershipsResponse)
	if err != nil {
		beego.Error(err.Error())
		return nil, errors.New("Failed to get user memberships")
	}
	membs = &membershipsResponse
	beego.Trace("Got num user memberships:", num)
	return

	// --- OLD ---
	// o := orm.NewOrm()
	// num, err := o.QueryTable("user_membership").Filter("user_id", userId).All(&ums)
	// if err != nil {
	// 	beego.Error("Failed to get user memberships")
	// 	return nil, errors.New("Failed to get user memberships")
	// }
	// beego.Trace("Got num user memberships:", num)
	// return
}

// Delete membership from the database
func DeleteUserMembership(userMembershipId int64) error {
	var num int64
	var err error
	o := orm.NewOrm()

	num, err = o.Delete(&UserMembership{Id: userMembershipId})
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete user membership: %v", err))
	}

	beego.Trace("Deleted num rows:", num)
	return nil
}
