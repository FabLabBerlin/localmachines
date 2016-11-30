package maintenance

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

const TABLE_NAME = "machine_maintenances"

type Maintenance struct {
	Id        int64
	MachineId int64
	Start     time.Time
	End       time.Time
}

func Create(mt *Maintenance) (err error) {
	o := orm.NewOrm()

	if mt.MachineId == 0 {
		return fmt.Errorf("no MachineId")
	}

	mt.Id, err = o.Insert(mt)

	return
}

func On(machineId int64) (mt *Maintenance, err error) {
	mt = &Maintenance{
		MachineId: machineId,
		Start:     time.Now(),
	}

	err = Create(mt)

	return
}

func (mt *Maintenance) TableName() string {
	return TABLE_NAME
}

func (mt *Maintenance) Off() (err error) {
	o := orm.NewOrm()
	mt.End = time.Now()

	_, err = o.Update(mt)

	return
}

func init() {
	orm.RegisterModel(new(Maintenance))
}
