package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(Activation))
}

type Activation struct {
	Id               int `orm:"auto";"pk"`
	InvoiceId        int `orm:"null"`
	UserId           int64
	MachineId        int
	Active           bool
	TimeStart        time.Time
	TimeEnd          time.Time `orm:"null"`
	TimeTotal        int
	UsedKwh          float32
	DiscountPercents float32
	DiscountFixed    float32
	VatRate          float32
	CommentRef       string `orm:"size(255)"`
	Invoiced         bool
	Changed          bool
}

func (this *Activation) TableName() string {
	return "activations"
}

type GetActivationsResponse struct {
	NumActivations  int64
	ActivationsPage *[]Activation
}

// Get filtered activations in a paged manner
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
	o := orm.NewOrm()
	num, err := o.QueryTable(act.TableName()).
		Filter("timeStart__gt", startTime).
		Filter("timeEnd__lt", endTime).
		Filter("invoiced", includeInvoiced).
		//Filter("userId", userId).
		Filter("active", false).
		OrderBy("userId", "-id").
		Limit(itemsPerPage).
		Offset(itemsPerPage * (page - 1)).
		All(&activations)

	if err != nil {
		msg := fmt.Sprintf("Failed to get activations: %v", err)
		return nil, errors.New(msg)
	}

	beego.Trace("Got num activations:", num)

	return &activations, nil
}

// Get number of matching activations.
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

// Get only active activations (active meaning that users are using
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
	return activations, nil
}

// Creates activation and returns activation ID
func CreateActivation(machineId, userId int64) (int64, error) {

	// Check if machine with machineId exists
	o := orm.NewOrm()
	machineExists := o.QueryTable("machine").Filter("Id", machineId).Exist()
	beego.Trace("machineExists:", machineExists)
	if !machineExists {
		return 0, errors.New("Machine with provided ID does not exist")
	}

	// Check for duplicate activations
	var dupActivations []Activation
	act := Activation{} // Used to get table name of the model
	query := fmt.Sprintf("SELECT id FROM %s WHERE machine_id = ? AND user_id = ? AND active = 1", act.TableName())
	numDuplicates, err := o.Raw(query, machineId, userId).
		QueryRows(&dupActivations)
	if err != nil {
		beego.Error("Could not get duplicate activations")
		return 0, err
	}
	if numDuplicates > 0 {
		beego.Error("Duplicate activations found:", numDuplicates)
		return 0, errors.New("Duplicate activations found")
	}

	// Beego model time stuff is bad, here we use workaround that works.
	// TODO: explore the beego ORM time management,
	// try to fix or use as it should be used.
	var res sql.Result
	query = fmt.Sprintf("INSERT INTO %s VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", act.TableName())
	res, err = o.Raw(query,
		nil, nil, userId, machineId, true,
		time.Now().Format("2006-01-02 15:04:05"),
		nil, 0, 0, 0, 0, 0, "", false, false).Exec()
	if err != nil {
		beego.Error("Failed to insert activation in to DB:", err)
		return 0, err
	}

	// Get the ID of the record we just inserted
	var activationId int64
	activationId, err = res.LastInsertId()
	if err != nil {
		beego.Error("Failed to get insterted activation ID")
		return 0, err
	}
	beego.Trace("Created activation with ID", activationId)

	// Update machine as unavailable
	_, err = o.QueryTable("machine").
		Filter("Id", machineId).
		Update(orm.Params{"available": false})
	if err != nil {
		beego.Error("Failed to update activated machine")
	}

	return activationId, nil
}

// Close running/active activation
func CloseActivation(activationId int64) error {

	// Get activation start time and machine id
	var tempModel struct {
		TimeStart string
		MachineId int
	}
	o := orm.NewOrm()
	act := Activation{}
	query := fmt.Sprintf("SELECT time_start, machine_id FROM %s WHERE active = true AND id = ?", act.TableName())
	err := o.Raw(query, activationId).QueryRow(&tempModel)
	if err != nil {
		beego.Error("Could not get activation:", err)
		return err
	}

	// Calculate activation duration
	const timeForm = "2006-01-02 15:04:05"
	timeStart, _ := time.ParseInLocation(timeForm, tempModel.TimeStart, time.Now().Location())
	timeNow := time.Now() // time.Now().Format("2006-01-02 15:04:05")
	totalDuration := timeNow.Sub(timeStart)

	// Update activation
	query = fmt.Sprintf("UPDATE %s SET active=false, time_end=?, time_total=? WHERE id=?", act.TableName())
	_, err = o.Raw(query,
		timeNow.Format("2006-01-02 15:04:05"),
		totalDuration.Seconds(), activationId).Exec()
	if err != nil {
		beego.Error("Failed to update activation:", err)
		return err
	}

	// Make the machine available
	_, err = o.QueryTable("machine").Filter("Id", tempModel.MachineId).
		Update(orm.Params{"available": true})
	if err != nil {
		beego.Error("Failed to available machine")
		return err
	}
	return nil
}

// Delete activation.
// It might happen that an activation is created by accident.
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

// Get the machine ID of a specific activation
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
