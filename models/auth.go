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
	Hash     string `orm:"size(300)"`
	Salt     string `orm:"size(100)"`
}
