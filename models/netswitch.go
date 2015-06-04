package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type NetSwitchMapping struct {
	Id        int64 `orm:"auto";"pk"`
	MachineId int64
	UrlOn     string `orm:"size(255)"`
	UrlOff    string `orm:"size(255)"`
}

func init() {
	orm.RegisterModel(new(NetSwitchMapping))
}

func (u *NetSwitchMapping) TableName() string {
	return "netswitch"
}

func CreateNetSwitchMapping(machineId int64) (int64, error) {
	mapping := NetSwitchMapping{}
	mapping.MachineId = machineId
	o := orm.NewOrm()
	created, id, err := o.ReadOrCreate(&mapping, "MachineId")
	if err != nil {
		return 0, err
	}
	if created {
		beego.Info("Created new NetSwitch mapping, machine ID:", id)
	} else {
		beego.Warning("A NetSwitch mapping exists, machine ID:", id)
	}
	return id, nil
}
