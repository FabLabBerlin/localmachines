package models

import (
	"encoding/json"
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
	Id                int64  `orm:"auto";"pk"`
	Name              string `orm:"size(255)"`
	Shortname         string `orm:"size(100)"`
	Description       string `orm:"type(text)"`
	Image             string `orm:"size(255)"` // TODO: media and media type tables
	Available         bool
	UnavailMsg        string    `orm:"type(text)"`
	UnavailTill       time.Time `orm:"null;type(date)" form:"Date,2006-01-02T15:04:05Z07:00`
	Price             float32
	PriceUnit         string `orm:"size(100)"`
	Comments          string `orm:"type(text)"`
	Visible           bool
	ConnectedMachines string `orm:"size(255)"`
}

// Define custom table name as for SQL table with a name "machines"
// makes more sense.
func (u *Machine) TableName() string {
	return "machines"
}

type ConnectedMachine struct {
	Id   int64
	Name string
}

type ConnectedMachineList struct {
	Data []*ConnectedMachine
}

type ConnectableMachine struct {
	Id   int64
	Name string
}

type ConnectableMachineList struct {
	Data []*ConnectableMachine
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

func GetConnectedMachines(machineId int64) (*ConnectedMachineList, error) {

	machine := Machine{}
	machine.Id = machineId

	o := orm.NewOrm()
	err := o.Read(&machine)
	if err != nil {
		return nil, fmt.Errorf("Failed to get connected machines: %v", err)
	}

	beego.Trace("connected machines:", machine.ConnectedMachines)

	machineList := ConnectedMachineList{}

	// Empty string, to not waste resources - return
	if machine.ConnectedMachines == "" {
		return &machineList, nil
	}

	// Parse string into object we can digest,
	// so we can load machine data individually
	var machineIds []int64
	err = json.Unmarshal([]byte(machine.ConnectedMachines), &machineIds)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json: %v", err)
	}

	// Load connected machine data from the database
	for _, val := range machineIds {
		m := Machine{}
		m.Id = val
		err = o.Read(&m)
		if err != nil {
			beego.Error("Failed to get connected machine from db, ID", val)
			return nil, fmt.Errorf("Failed to get connected machine: %v", err)
		}
		cm := ConnectedMachine{}
		cm.Id = m.Id
		cm.Name = m.Name
		machineList.Data = append(machineList.Data, &cm)
	}

	//machineList.Data = append(machineList.Data, &machine1)
	//machineList.Data = append(machineList.Data, &machine2)

	return &machineList, nil
}

func GetConnectableMachines(machineId int64) (*ConnectableMachineList, error) {

	// All machines can be connectable
	var machines []*Machine
	var err error
	machines, err = GetAllMachines()
	if err != nil {
		return nil, err
	}

	// We have to substract the ones connected already from
	// the full machine list
	var machineList *ConnectedMachineList
	machineList, err = GetConnectedMachines(machineId)
	if err != nil {
		return nil, err
	}

	cmList := ConnectableMachineList{}

MachineLoop:
	for _, machine := range machines {
		if machine.Id == machineId {
			continue MachineLoop
		}

		for _, connMachine := range machineList.Data {
			if machine.Id == connMachine.Id {
				continue MachineLoop
			}
		}

		cm := ConnectableMachine{}
		cm.Id = machine.Id
		cm.Name = machine.Name
		cmList.Data = append(cmList.Data, &cm)
	}

	return &cmList, nil
}
