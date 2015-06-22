package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(Machine))
}

type Machine struct {
	Id          int64  `orm:"auto";"pk"`
	Name        string `orm:"size(255)"`
	Shortname   string `orm:"size(100)"`
	Description string `orm:"type(text)"`
	Image       string `orm:"size(255)"` // TODO: media and media type tables
	Available   bool
	UnavailMsg  string    `orm:"type(text)"`
	UnavailTill time.Time `orm:"null;type(date)" form:"Date,2006-01-02T15:04:05Z07:00`
	Price       float32
	PriceUnit   string `orm:"size(100)"`
	Comments    string `orm:"type(text)"`
	Visible     bool
}

// Define custom table name as for SQL table with a name "machines"
// makes more sense.
func (u *Machine) TableName() string {
	return "machines"
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
	num, err := o.QueryTable("machines").All(&machines)
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

// Update existing machine in the database
func UpdateMachine(machine *Machine) error {
	var err error
	var num int64

	o := orm.NewOrm()
	num, err = o.Update(machine)
	if err != nil {
		return err
	}

	beego.Trace("Rows affected:", num)
	return nil
}

// Delete machine from the database
func DeleteMachine(machineId int64) error {
	var num int64
	var err error
	o := orm.NewOrm()

	// delete machine
	num, err = o.Delete(&Machine{Id: machineId})
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete machine: %v", err))
	}
	beego.Trace("Deleted num machines:", num)

	// delete switch mapping
	swch := HexabusMapping{}
	num, err = o.QueryTable(swch.TableName()).Filter("machine_id",
		machineId).Delete()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete switch mapping: %v", err))
	}
	beego.Trace("Deleted num switch mappings:", num)

	// Delete activations assigned to machine
	act := Activation{}
	num, err = o.QueryTable(act.TableName()).Filter("machine_id",
		machineId).Delete()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete activations: %v", err))
	}
	beego.Trace("Deleted num activations:", num)

	// Delete user machine permissions of this machine
	perm := Permission{}
	num, err = o.QueryTable(perm.TableName()).Filter("machine_id",
		machineId).Delete()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete machine permissions: %v", err))
	}
	beego.Trace("Deleted num machine permissions:", num)

	return nil
}
