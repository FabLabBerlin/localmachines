package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(HexabusMapping))
}

// Holds hexabus switch mapping between a machine and switch
type HexabusMapping struct {
	Id        int64 `orm:"auto";"pk"`
	MachineId int64
	SwitchIp  string `orm:"size(100)"`
}

// For now the table name will differ from the model name
func (u *HexabusMapping) TableName() string {
	return "hexaswitch"
}

// Create new mapping and return its ID
func CreateHexabusMapping(machineId int64) (int64, error) {
	mapping := HexabusMapping{}
	mapping.MachineId = machineId
	o := orm.NewOrm()
	created, id, err := o.ReadOrCreate(&mapping, "MachineId")
	if err != nil {
		return 0, err
	}
	if created {
		beego.Info("Created new hexabus mapping, machine id:", id)
	} else {
		beego.Warning("A hexabus mapping exists, machine id:", id)
	}
	return id, nil
}

// Get hexabus mapping by using a machine ID
func GetHexabusMapping(machineId int64) (*HexabusMapping, error) {
	mapping := HexabusMapping{}
	mapping.MachineId = machineId
	o := orm.NewOrm()
	err := o.Read(&mapping, "MachineId")
	if err != nil {
		return &mapping, err
	}
	return &mapping, nil
}

// Delete hexabus mapping by machine ID
func DeleteHexabusMapping(machineId int64) error {

	mapping, err := GetHexabusMapping(machineId)
	if err != nil {
		return err
	}

	var num int64

	o := orm.NewOrm()
	num, err = o.Delete(mapping)
	if err != nil {
		return err
	}

	beego.Trace("Affected num rows:", num)
	return nil
}

// Update hexabus mapping
func UpdateHexabusMapping(mapping *HexabusMapping) error {
	o := orm.NewOrm()
	num, err := o.Update(mapping)
	if err != nil {
		return err
	}
	beego.Trace("Affected num rows while updating:", num)
	return nil
}
