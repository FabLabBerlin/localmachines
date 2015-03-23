package models

import (
	"database/sql"
	"errors"
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

func GetActiveActivations() ([]*Activation, error) {
	var activations []*Activation
	o := orm.NewOrm()
	num, err := o.QueryTable("activation").
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
	numDuplicates, err := o.Raw("SELECT id FROM activation WHERE machine_id = ? AND user_id = ? AND active = 1",
		machineId, userId).QueryRows(&dupActivations)
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
	res, err = o.Raw("INSERT INTO activation VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
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

func CloseActivation(activationId int64) error {

	// Get activation start time and machine id
	var tempModel struct {
		TimeStart string
		MachineId int
	}
	o := orm.NewOrm()
	err := o.Raw("SELECT time_start, machine_id FROM activation WHERE active = true AND id = ?",
		activationId).QueryRow(&tempModel)
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
	_, err = o.Raw("UPDATE activation SET active=false, time_end=?, time_total=? WHERE id=?",
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

func DeleteActivation(activationId int64) error {

	// Set machine of the activation available
	var err error
	var activation Activation
	o := orm.NewOrm()
	err = o.QueryTable("activation").
		Filter("Id", activationId).
		One(&activation, "MachineId")
	if err != nil {
		beego.Error("Failed to get machine ID of the activation")
		return err
	}
	_, err = o.QueryTable("machine").
		Filter("Id", activation.MachineId).
		Update(orm.Params{"available": true})
	if err != nil {
		beego.Error("Failed to update machine as available")
		return err
	}

	// Remove the activation
	_, err = o.QueryTable("activation").
		Filter("Id", activationId).Delete()
	if err != nil {
		beego.Error("Failed to delete activation")
		return err
	}
	return nil
}

func GetActivationMachineId(activationId int64) (int64, error) {
	activationModel := Activation{}
	o := orm.NewOrm()
	err := o.QueryTable("activation").
		Filter("id", activationId).
		One(&activationModel, "MachineId")
	if err != nil {
		beego.Error("Could not get activation")
		return 0, err
	}
	return int64(activationModel.MachineId), nil
}
