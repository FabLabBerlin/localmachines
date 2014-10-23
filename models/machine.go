package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(Machine))
}

type Machine struct {
	Id           int    `orm:"auto";"pk"`
	Name         string `orm:"size(255)"`
	Description  string `orm:"type(text)"`
	Available    bool
	UnavailMsg   string `orm:"type(text)"`
	UnavailTill  time.Time
	CalcByEnergy bool
	CalcByTime   bool
	CostsPerKwh  float32
	CostsPerMin  float32
}
