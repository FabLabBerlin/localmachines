package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"errors"
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

func GetMachine(machineId int) (*Machine, error) {
	machine := Machine{Id: machineId}
	o := orm.NewOrm()
	err := o.Read(&machine)
	if err != nil {
		return nil, errors.New("Failed to get machine")
	}
	return &machine, nil
}

func GetAllMachines() ([]*Machine, error) {
	var machines []*Machine
	o := orm.NewOrm()
	num, err := o.QueryTable("machine").All(&machines)
	if err != nil {
		return nil, errors.New("Failed to get all machines")
	}
	beego.Trace("Got num machines:", num)
	return machines, nil
}