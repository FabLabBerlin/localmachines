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
	Price                 int
	MachinePriceDeduction int
}

type UserMembership struct {
	Id           int64 `orm:"auto";"pk"`
	UserId       int64
	MembershipId int64
	StartDate    time.Time
}

func init() {
	orm.RegisterModel((new(Membership)))
}

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
