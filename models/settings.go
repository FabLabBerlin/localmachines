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

func (this *Setting) TableName() string {
	return "settings"
}

func CreateSetting(setting *Setting) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(setting)
}

func GetAllSettings() (settings []*Setting, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(Setting).TableName()).All(&settings)
	return
}

func UpdateSetting(setting *Setting) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(setting)
	return
}

func init() {
	orm.RegisterModel(new(Setting))
}
