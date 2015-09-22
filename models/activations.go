package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// Activation type/model to hold information of single activation.
type Activation struct {
	Id               int64 `orm:"auto";"pk"`
	InvoiceId        int   `orm:"null"`
	UserId           int64
	MachineId        int64
	Active           bool
	TimeStart        time.Time
	TimeEnd          time.Time `orm:"null"`
	TimeTotal        int64
	UsedKwh          float32
	DiscountPercents float32
	DiscountFixed    float32
	VatRate          float32
	CommentRef       string `orm:"size(255)"`
	Invoiced         bool
	Changed          bool
}

// Returns mysql table name of the table mapped to the Activation model.
func (this *Activation) TableName() string {
	return "activations"
}

func init() {
	orm.RegisterModel(new(Activation))
}

// Type to be used as response model for the HTTP GET activations request.
type GetActivationsResponse struct {
	NumActivations  int64
	ActivationsPage *[]Activation
}

// Gets activations of a specific user between specified start and end time.
func GetUserActivations(startTime time.Time,
	endTime time.Time,
	userId int64) (*[]Activation, error) {

	// Get activations from database
	activations := []Activation{}
	act := Activation{}
	usr := User{}
	o := orm.NewOrm()

	query := fmt.Sprintf("SELECT a.* FROM %s a JOIN %s u ON a.user_id=u.id "+
		"WHERE a.time_start>? AND a.time_end<? AND a.active=false AND a.user_id=? "+
		"ORDER BY u.first_name ASC, a.time_start DESC",
		act.TableName(),
		usr.TableName())

	beego.Trace(query)

	num, err := o.Raw(query,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
		userId).QueryRows(&activations)

	if err != nil {
		msg := fmt.Sprintf("Failed to get activations: %v", err)
		return nil, errors.New(msg)
	}

	beego.Trace("Got num activations:", num)

	return &activations, nil
}

// Gets activations of a specific user by consuming user ID.
func GetUserActivationsStartTime(userId int64) (startTime time.Time, err error) {
	query := "SELECT min(time_start) FROM activations WHERE user_id = ?"
	o := orm.NewOrm()
	err = o.Raw(query, userId).QueryRow(&startTime)
	return
}

// Gets filtered activations in a paged manner between start and end time.
// Items per page and page number can be specified. Already invoiced
// activations can be excluded.
func GetActivations(startTime time.Time,
	endTime time.Time,
	userId int64,
	includeInvoiced bool,
	itemsPerPage int64,
	page int64) (*[]Activation, error) {

	if page <= 0 {
		page = 1
	}

	// Get activations from database
	activations := []Activation{}
	act := Activation{}
	usr := User{}
	o := orm.NewOrm()

	var pageOffset int64
	pageOffset = itemsPerPage * (page - 1)

	query := fmt.Sprintf("SELECT a.* FROM %s a JOIN %s u ON a.user_id=u.id "+
		"WHERE a.time_start>? AND a.time_end<? AND a.invoiced=? AND a.active=false "+
		"ORDER BY u.first_name ASC, a.time_start DESC "+
		"LIMIT ? OFFSET ?",
		act.TableName(),
		usr.TableName())

	num, err := o.Raw(query,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
		includeInvoiced,
		itemsPerPage,
		pageOffset).QueryRows(&activations)

	if err != nil {
		msg := fmt.Sprintf("Failed to get activations: %v", err)
		return nil, errors.New(msg)
	}

	beego.Trace("Got num activations:", num)

	return &activations, nil
}

// Gets number of matching activations.
// Used together with GetActivations.
func GetNumActivations(startTime time.Time,
	endTime time.Time,
	userId int64,
	includeInvoiced bool) (int64, error) {

	// Count activations matching params
	o := orm.NewOrm()
	act := Activation{}
	cnt, err := o.QueryTable(act.TableName()).
		Filter("timeStart__gt", startTime).
		Filter("timeEnd__lt", endTime).
		Filter("invoiced", includeInvoiced).
		//Filter("userId", userId).
		Filter("active", false).
		Count()

	if err != nil {
		msg := fmt.Sprintf("Failed to count activations: %v", err)
		return 0, errors.New(msg)
	}

	beego.Trace("Num activations matches:", cnt)

	return cnt, nil
}

// Gets only active activations (active meaning that users are using
// the machine/resource currently)
func GetActiveActivations() ([]*Activation, error) {
	var activations []*Activation
	o := orm.NewOrm()
	act := Activation{}
	num, err := o.QueryTable(act.TableName()).
		Filter("active", true).All(&activations)
	if err != nil {
		return nil, errors.New("Failed to get active activations")
	}
	beego.Trace("Got num activations:", num)

	// Calculate total time for all activations
	for i := 0; i < len(activations); i++ {
		timeNow := time.Now()
		activations[i].TimeTotal =
			int64(timeNow.Sub(activations[i].TimeStart).Seconds())
	}

	return activations, nil
}

// Creates activation and returns activation ID.
func CreateActivation(machineId, userId int64, startTime time.Time) (activationId int64, err error) {

	// Check if machine with machineId exists
	// TODO: Replace this with a more readable helper function
	o := orm.NewOrm()
	mch := Machine{}
	machineExists := o.QueryTable(mch.TableName()).Filter("Id", machineId).Exist()
	beego.Trace("machineExists:", machineExists)
	if !machineExists {
		return 0, fmt.Errorf("Machine with provided ID does not exist")
	}

	// Check for duplicate activations
	// TODO: Replace this with a more readable helper function
	var dupActivations []Activation
	act := Activation{} // Used to get table name of the model
	query := fmt.Sprintf("SELECT id FROM %s WHERE machine_id = ? "+
		"AND user_id = ? AND active = 1", act.TableName())
	numDuplicates, err := o.Raw(query, machineId, userId).
		QueryRows(&dupActivations)
	if err != nil {
		return 0, fmt.Errorf("Could not get duplicate activations: %v", err)
	}
	if numDuplicates > 0 {
		beego.Error("Duplicate activations found:", numDuplicates)
		return 0, fmt.Errorf("Duplicate activations found")
	}

	// Check if the machine is available
	machineAvailable := o.QueryTable(mch.TableName()).
		Filter("Id", machineId).
		Filter("Available", true).
		Exist()

	beego.Trace("machineAvailable:", machineAvailable)

	if !machineAvailable {
		return 0, fmt.Errorf("Machine ID %s not available", machineId)
	}

	newActivation := Activation{}
	newActivation.UserId = userId
	newActivation.MachineId = machineId
	newActivation.Active = true
	newActivation.TimeStart = startTime
	activationId, err = o.Insert(&newActivation)
	if err != nil {
		beego.Error("Failed to insert activation:", err)
		return 0, fmt.Errorf("Failed to insert activation %v", err)
	}
	beego.Trace("Created activation with ID", activationId)

	// Update machine as unavailable
	_, err = o.QueryTable(mch.TableName()).
		Filter("Id", machineId).
		Update(orm.Params{"available": false})
	if err != nil {
		beego.Error("Failed to update activated machine")
	}

	return activationId, nil
}

func GetActivation(activationId int64) (activation *Activation, err error) {
	activation = &Activation{}
	activation.Id = activationId

	o := orm.NewOrm()
	err = o.Read(activation)

	if err != nil {
		beego.Error("Failed to read activation:", err)
		return nil, fmt.Errorf("Failed to read activation: %v", err)
	}

	return
}

// Close running/active activation.
func CloseActivation(activationId int64, endTime time.Time) error {

	activation, err := GetActivation(activationId)
	if err != nil {
		beego.Error("Failed to get activation:", err)
		return fmt.Errorf("Failed to get activation: %v", err)
	}

	// Calculate activation duration and update activation.
	activation.TimeEnd = endTime
	activation.TimeTotal = int64(endTime.Sub(activation.TimeStart).Seconds())

	err = UpdateActivation(activation)
	if err != nil {
		beego.Error("Failed to update activation:", err)
		return fmt.Errorf("Failed to update activation: %v", err)
	}

	// Make the machine available again.
	var machine *Machine
	machine, err = GetMachine(activation.MachineId)
	if err != nil {
		beego.Error("Failed to get machine:", err)
		return fmt.Errorf("Failed to get machine: %v", err)
	}

	machine.Available = true
	err = UpdateMachine(machine)
	if err != nil {
		beego.Error("Failed to update machine:", err)
		return fmt.Errorf("Failed to update machine: %v", err)
	}

	return nil
}

func UpdateActivation(activation *Activation) error {
	o := orm.NewOrm()
	num, err := o.Update(activation)

	if err != nil {
		beego.Error("Failed to update activation:", err)
		return fmt.Errorf("Failed to update activation: %v", err)
	}

	beego.Trace("Affected num rows:", num)

	return nil
}

// Delete activation matching an activation ID.
func DeleteActivation(activationId int64) error {

	// Set machine of the activation available
	var err error
	var activation Activation
	o := orm.NewOrm()
	err = o.QueryTable(activation.TableName()).
		Filter("Id", activationId).
		One(&activation, "MachineId")
	if err != nil {
		beego.Error("Failed to get machine ID of the activation")
		return err
	}
	m := Machine{}
	_, err = o.QueryTable(m.TableName()).
		Filter("Id", activation.MachineId).
		Update(orm.Params{"available": true})
	if err != nil {
		beego.Error("Failed to update machine as available")
		return err
	}

	// Remove the activation
	_, err = o.QueryTable(activation.TableName()).
		Filter("Id", activationId).Delete()
	if err != nil {
		beego.Error("Failed to delete activation")
		return err
	}
	return nil
}

// Gets the machine ID of a specific activation defined by activation ID.
func GetActivationMachineId(activationId int64) (int64, error) {
	activationModel := Activation{}
	o := orm.NewOrm()
	err := o.QueryTable(activationModel.TableName()).
		Filter("id", activationId).
		One(&activationModel, "MachineId")
	if err != nil {
		beego.Error("Could not get activation")
		return 0, err
	}
	return int64(activationModel.MachineId), nil
}
