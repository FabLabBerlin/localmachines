package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Permission))
}

type Permission struct {
	UserId    int `orm:"auto"`
	MachineId int
}
