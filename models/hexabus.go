package models

import (
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
func CreateHexabusMapping(machineId int64, switchIp string) (int64, error) {
	mapping := HexabusMapping{}
	mapping.MachineId = machineId
	mapping.SwitchIp = switchIp
	o := orm.NewOrm()
	id, err := o.Insert(&mapping)
	if err != nil {
		return 0, err
	}
	return id, nil
}

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
