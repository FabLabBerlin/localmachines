package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"sort"
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
	Price             float64
	PriceUnit         string `orm:"size(100)"`
	Comments          string `orm:"type(text)"`
	Visible           bool
	ConnectedMachines string `orm:"size(255)"`
	SwitchRefCount    int64
	UnderMaintenance  bool
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

type ExtendedMachine struct {
	NumActivations int64
	MachineData    *Machine
}
type ExtendedMachineList []*ExtendedMachine

func (s ExtendedMachineList) Len() int {
	return len(s)
}
func (s ExtendedMachineList) Less(i, j int) bool {
	return s[i].NumActivations > s[j].NumActivations
}
func (s ExtendedMachineList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func GetAllMachines() ([]*Machine, error) {
	var machines []Machine
	o := orm.NewOrm()
	m := Machine{}
	num, err := o.QueryTable(m.TableName()).All(&machines)
	if err != nil {
		return nil, fmt.Errorf("Failed to get all machines: %v", err)
	}
	beego.Trace("Got num machines:", num)

	var extendedMachines ExtendedMachineList

	// Get sum of activations per machine
	a := Activation{}
	for i := 0; i < len(machines); i++ {
		var numActivations int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE machine_id=?",
			a.TableName())
		//beego.Trace("Counting activations for machine with ID", machines[i].Id)
		err = o.Raw(query, machines[i].Id).QueryRow(&numActivations)
		if err != nil {
			return nil, fmt.Errorf("Failed to read activations: %v", err)
		}
		extendedMachine := ExtendedMachine{}
		extendedMachine.NumActivations = numActivations
		extendedMachine.MachineData = &machines[i]
		extendedMachines = append(extendedMachines, &extendedMachine)
	}

	// Sort the machines by number of activations
	sort.Sort(extendedMachines)
	// Build new machine slice
	var sortedMachines []*Machine
	for i := 0; i < len(extendedMachines); i++ {
		sortedMachines = append(sortedMachines, extendedMachines[i].MachineData)
	}

	return sortedMachines, nil
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

func (this *Machine) On() error {

	// Get current switch reference count
	var err error
	o := orm.NewOrm()
	if err = o.Read(this); err != nil {
		return fmt.Errorf("Failed to get machine switch ref count: %v", err)
	}

	var netSwitchMapping *NetSwitchMapping = nil
	netSwitchMapping, err = GetNetSwitchMapping(this.Id)
	if err != nil {
		beego.Warning("Failed to get NetSwitch mapping:", err)
	} else if netSwitchMapping != nil {
		if err = netSwitchMapping.On(); err != nil {
			return fmt.Errorf("Failed to turn on NetSwitch: %v", err)
		}
	}

	// Increase and update switch reference count
	this.SwitchRefCount += 1
	var num int64
	num, err = o.Update(this)
	if err != nil {
		return fmt.Errorf("Failed to update machine switch ref count: %v", err)
	}
	beego.Trace("Num affected rows during switch ref update:", num)

	return nil
}

func (this *Machine) ReportBroken() error {
	email := NewEmail()
	to := beego.AppConfig.String("trelloemail")
	subject := this.Name + " reported as broken"
	message := "The machine " + this.Name + " seems to be broken."
	return email.Send(to, subject, message)
}

func (this *Machine) SetUnderMaintenance(underMaintenance bool) error {
	this.UnderMaintenance = underMaintenance

	if err := UpdateMachine(this); err != nil {
		return err
	}

	consumerKey := beego.AppConfig.String("maintenancetwitterconsumerkey")
	consumerSecret := beego.AppConfig.String("maintenancetwitterconsumersecret")
	key := beego.AppConfig.String("maintenancetwitteraccesskey")
	secret := beego.AppConfig.String("maintenancetwitteraccesssecret")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(key, secret)
	defer api.Close()

	var msg string
	if underMaintenance {
		msg = "The " + this.Name + " is undergoing maintenance works right now 😟"
	} else {
		msg = "The " + this.Name + " works again!!! 😀"
	}

	// If the tweet fails, we should not worry.
	// This should not abort the maintenance call.
	api.PostTweet(msg, nil)

	return nil
}

func (this *Machine) Off() error {

	// Get current switch reference count
	var err error
	o := orm.NewOrm()
	if err = o.Read(this); err != nil {
		return fmt.Errorf("Failed to get machine switch ref count: %v", err)
	}

	this.SwitchRefCount -= 1
	if this.SwitchRefCount < 0 {
		this.SwitchRefCount = 0
	}

	// Update reference count
	var num int64
	num, err = o.Update(this)
	if err != nil {
		return fmt.Errorf("Failed to update machine switch ref count: %v", err)
	}
	beego.Trace("Num affected rows during switch ref update:", num)

	// Turn off the switch physically only when reference count reaches 0
	if this.SwitchRefCount > 0 {
		return nil
	}

	var netSwitch *NetSwitchMapping = nil
	netSwitch, err = GetNetSwitchMapping(this.Id)
	if err != nil {
		beego.Warning("Failed to get NetSwitch mapping:", err)
	} else if netSwitch != nil {
		if err = netSwitch.Off(); err != nil {
			return fmt.Errorf("Failed to turn off NetSwitch: %v", err)
		}
	}
	return nil
}

func (this *ConnectedMachineList) On() error {

	for _, cm := range this.Data {
		m := Machine{}
		m.Id = cm.Id
		if err := m.On(); err != nil {
			return fmt.Errorf("Failed to turn on connected machine: %v", err)
		}
	}

	return nil
}

func (this *ConnectedMachineList) Off() error {

	for _, cm := range this.Data {
		m := Machine{}
		m.Id = cm.Id
		if err := m.Off(); err != nil {
			return fmt.Errorf("Failed to turn off connected machine: %v", err)
		}
	}

	return nil
}
