package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(User))
}

type User struct {
	Id          int    `orm:"auto";"pk"`
	FirstName   string `orm:"size(100)"`
	LastName    string `orm:"size(100)"`
	Username    string `orm:"size(100)"`
	Email       string `orm:"size(100)"`
	InvoiceAddr int
	ShipAddr    int
	ClientId    int
	B2b         bool
	Company     string `orm:"size(100)"`
	VatUserId   string `orm:"size(100)"`
	VatRate     int
}

