package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(UserRoles))
}

type UserRoles struct {
	UserId int  `orm:"auto"`
	Admin  bool `orm:"size(100)"`
	Staff  bool `orm:"size(100)"`
	Member bool `orm:"size(100)"`
}
