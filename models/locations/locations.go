package locations

import (
	"github.com/astaxie/beego/orm"
)

const TABLE_NAME = "locations"

type Location struct {
	Id   int64
	Name string `orm:"size(100)"`
}

func init() {
	orm.RegisterModel(new(Location))
}

func (this *Location) TableName() string {
	return TABLE_NAME
}

func GetAll() (ls []*Location, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).All(&ls)
	return
}
