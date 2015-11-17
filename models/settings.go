package models

import (
	"github.com/astaxie/beego/orm"
)

type Setting struct {
	Id          int64  `orm:"auto";"pk"`
	Name        string `orm:"size(100)"`
	ValueInt    *int64
	ValueString *string
	ValueFloat  *float64
}

func init() {
	orm.RegisterModel(new(Setting))
}
