package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(Machine))
}

type Machine struct {
	Id           int64  `orm:"auto";"pk"`
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

func GetMachine(machineId int64) (*Machine, error) {
	beego.Trace(machineId)
	machine := Machine{Id: machineId}
	o := orm.NewOrm()
	err := o.Read(&machine)
	if err != nil {
		return nil, err
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

func CreateMachine(machineName string) (int64, error) {
	o := orm.NewOrm()
	machine := Machine{Name: machineName, Available: true}
	id, err := o.Insert(&machine)
	if err == nil {
		return id, nil
	} else {
		return 0, err
	}
}
