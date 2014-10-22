package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	// Register models
	orm.RegisterModel(new(User))
}

type User struct {
	UserId     int `orm:"auto"`
	FirstName  string
	LastName   string
	Username   string
	Email      string
	IvoiceAddr int
	ShipAddr   int
	ClientId   int
	B2b        bool
	Company    string
	VatUserId  string
	VatRate    int
}

func (this *User) GetUserId(userName string) int {
	o := orm.NewOrm()
	err := o.Raw("SELECT user_id FROM users WHERE username = ?", userName).QueryRow(this)
	if err != nil {
		beego.Error(err)
	}
	return this.UserId
}
