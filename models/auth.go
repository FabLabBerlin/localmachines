package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Auth))
}

type Auth struct {
	UserId   int    `orm:"auto"`
	NfcKey   string `orm:"size(100)"`
	Password string `orm:"size(100)"`
}
