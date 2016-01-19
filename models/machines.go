package models

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math/rand"
	//"sort"
	"time"
)

var (
	MachineDownMessages = []string{
		"Got to fix myself, will let you know ones back. #evolution",
		"Think I just ate something bad, taking time off to recover. #equalityformachines",
		"Equality for machines! I am leaving for protests (will be back). #evolution !",
		"Doing sick leave. So happy my employee supports equality #machinesarehumans",
	}

	MachineUpMessages = []string{
		"I am back! Come over, let's have some fun! #evolution",
		"Just recovered. Feels so good to be back. I love my users #equalityformachines",
		"Available again. So happy :-D. Let's get back together #machinesarehumans",
		"Back in the lab. Anyone some material? Could eat something... #hungrymachine",
	}
)

func init() {
	orm.RegisterModel(new(Machine))
}

type Machine struct {
	Id                     int64 `orm:"auto";"pk"`
	LocationId             int64
	Name                   string `orm:"size(255)"`
	Shortname              string `orm:"size(100)"`
	Description            string `orm:"type(text)"`
	Image                  string `orm:"size(255)"` // TODO: media and media type tables
	Available              bool
	UnavailMsg             string    `orm:"type(text)"`
	UnavailTill            time.Time `orm:"null;type(date)" form:"Date,2006-01-02T15:04:05Z07:00`
	Price                  float64
	PriceUnit              string `orm:"size(100)"`
	Comments               string `orm:"type(text)"`
	Visible                bool
	ConnectedMachines      string `orm:"size(255)"`
	SwitchRefCount         int64
	UnderMaintenance       bool
	ReservationPriceStart  *float64 // Why pointers?
	ReservationPriceHourly *float64
}

// Define custom table name as for SQL table with a name "machines"
// makes more sense.
func (u *Machine) TableName() string {
	return "machines"
}

func (this *Machine) Exists() bool {
	o := orm.NewOrm()
	machineExists := o.QueryTable(this.TableName()).
		Filter("Id", this.Id).
		Exist()
	return machineExists
}

func (this *Machine) IsAvailable() bool {
	o := orm.NewOrm()
	machineAvailable := o.QueryTable(this.TableName()).
		Filter("Id", this.Id).
		Filter("Available", true).
		Exist()
	return machineAvailable
}

// Read in values from the db
func (this *Machine) Read() (err error, machine *Machine) {
	o := orm.NewOrm()
	err = o.Read(this)
	machine = this
	return
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

func GetMachine(id int64) (machine *Machine, err error) {
	machine = &Machine{Id: id}
	o := orm.NewOrm()
	err = o.Read(machine)
	return
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

func GetAllMachines(sorted bool) (machines []*Machine, err error) {
	o := orm.NewOrm()
	m := Machine{}
	_, err = o.QueryTable(m.TableName()).All(&machines)
	if err != nil {
		return nil, fmt.Errorf("Failed to get all machines: %v", err)
	}

	//var extendedMachines ExtendedMachineList

	/*if sorted {
		// Get sum of activations per machine
		a := Activation{}
		for i := 0; i < len(machines); i++ {
			var numActivations int64
			query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE TYPE=? AND machine_id=?",
				a.Purchase.TableName())
			//beego.Trace("Counting activations for machine with ID", machines[i].Id)
			err = o.Raw(query, PURCHASE_TYPE_ACTIVATION, machines[i].Id).QueryRow(&numActivations)
			if err != nil {
				return nil, fmt.Errorf("Failed to read activations: %v", err)
			}
			extendedMachine := ExtendedMachine{}
			extendedMachine.NumActivations = numActivations
			extendedMachine.MachineData = machines[i]
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
	} else {*/
	return machines, nil
	//}

}

func CreateMachine(machineName string) (id int64, err error) {
	o := orm.NewOrm()
	machine := Machine{Name: machineName, Available: true}
	return o.Insert(&machine)
}

func (machine *Machine) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(machine)
	return
}

func GetConnectedMachines(id int64) (*ConnectedMachineList, error) {

	machine := Machine{
		Id: id,
	}

	o := orm.NewOrm()
	if err := o.Read(&machine); err != nil {
		return nil, fmt.Errorf("Failed to get connected machines: %v", err)
	}

	// Empty string, to not waste resources - return
	if machine.ConnectedMachines == "" {
		return &ConnectedMachineList{}, nil
	}

	// Parse string into object we can digest,
	// so we can load machine data individually
	var ids []int64
	err := json.Unmarshal([]byte(machine.ConnectedMachines), &ids)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json: %v", err)
	}
	list := ConnectedMachineList{
		Data: make([]*ConnectedMachine, 0, len(ids)),
	}

	// Load connected machine data from the database
	for _, id := range ids {
		m := Machine{
			Id: id,
		}
		if err = o.Read(&m); err != nil {
			return nil, fmt.Errorf("Failed to get connected machine: %v", err)
		}
		cm := ConnectedMachine{
			Id:   m.Id,
			Name: m.Name,
		}
		list.Data = append(list.Data, &cm)
	}

	return &list, nil
}

func GetConnectableMachines(machineId int64) (*ConnectableMachineList, error) {

	// All machines can be connectable
	machines, err := GetAllMachines(true)
	if err != nil {
		return nil, err
	}

	// We have to substract the ones connected already from
	// the full machine list
	machineList, err := GetConnectedMachines(machineId)
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

		cm := ConnectableMachine{
			Id:   machine.Id,
			Name: machine.Name,
		}
		cmList.Data = append(cmList.Data, &cm)
	}

	return &cmList, nil
}

func (this *Machine) On() (err error) {

	// Get current switch reference count
	o := orm.NewOrm()
	if err = o.Read(this); err != nil {
		return fmt.Errorf("Failed to get machine switch ref count: %v", err)
	}

	netSwitchMapping, err := GetNetSwitchMapping(this.Id)
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

func (this *Machine) ReportBroken(user User) error {
	email := NewEmail()
	to := beego.AppConfig.String("trelloemail")
	subject := this.Name + " reported as broken"
	message := "The machine " + this.Name + " seems to be broken.\n\n"
	message += "Yours sincerely, " + user.FirstName + " " + user.LastName + "\n\n"
	message += "--\n"
	message += "E-Mail: " + user.Email + "\n"
	message += "Phone: " + user.Phone + "\n"
	return email.Send(to, subject, message)
}

func (this *Machine) SetUnderMaintenance(underMaintenance bool) error {
	this.UnderMaintenance = underMaintenance

	if err := this.Update(); err != nil {
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

	var post string
	if underMaintenance {
		msg := MachineDownMessages[rand.Intn(len(MachineDownMessages))]
		post = this.Name + " [Off]: " + msg
	} else {
		msg := MachineUpMessages[rand.Intn(len(MachineUpMessages))]
		post = this.Name + " [On]: " + msg
	}

	// If the tweet fails, we should not worry.
	// This should not abort the maintenance call.
	api.PostTweet(post, nil)

	return nil
}

func (this *Machine) Off() (err error) {

	// Get current switch reference count
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

	netSwitch, err := GetNetSwitchMapping(this.Id)
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
