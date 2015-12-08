package purchases

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
	"time"
)

type Activation struct {
	json.Marshaler
	Purchase Purchase
}

func (this *Activation) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Purchase)
}

// Type to be used as response model for the HTTP GET activations request.
type GetActivationsResponse struct {
	NumActivations  int64
	ActivationsPage *[]Activation
}

// Gets filtered activations in a paged manner between start and end time.
// Items per page and page number can be specified. Already invoiced
// activations can be excluded.
func GetActivations(startTime time.Time,
	endTime time.Time,
	userId int64,
	itemsPerPage int64,
	page int64) (*[]Activation, error) {

	if page <= 0 {
		page = 1
	}

	// Get activations from database
	purchases := []*Purchase{}
	act := Activation{}
	usr := models.User{}
	o := orm.NewOrm()

	var pageOffset int64
	pageOffset = itemsPerPage * (page - 1)

	query := fmt.Sprintf("SELECT a.* FROM %s a JOIN %s u ON a.user_id=u.id "+
		"WHERE a.type=? AND a.time_start>? AND a.time_end<? AND a.running=false "+
		"ORDER BY u.first_name ASC, a.time_start DESC "+
		"LIMIT ? OFFSET ?",
		act.Purchase.TableName(),
		usr.TableName())

	_, err := o.Raw(query,
		TYPE_ACTIVATION,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
		itemsPerPage,
		pageOffset).QueryRows(&purchases)

	if err != nil {
		msg := fmt.Sprintf("Failed to get activations: %v", err)
		return nil, errors.New(msg)
	}

	activations := make([]Activation, 0, len(purchases))
	for _, purchase := range purchases {
		act := Activation{
			Purchase: *purchase,
		}
		activations = append(activations, act)
	}

	return &activations, nil
}

// Gets number of matching activations.
// Used together with GetActivations.
func GetNumActivations(startTime time.Time,
	endTime time.Time,
	userId int64) (int64, error) {

	// Count activations matching params
	o := orm.NewOrm()
	act := Activation{}
	cnt, err := o.QueryTable(act.Purchase.TableName()).
		Filter("timeStart__gt", startTime).
		Filter("timeEnd__lt", endTime).
		//Filter("userId", userId).
		Filter("Running", false).
		Filter("type", TYPE_ACTIVATION).
		Count()

	if err != nil {
		msg := fmt.Sprintf("Failed to count activations: %v", err)
		return 0, errors.New(msg)
	}

	beego.Trace("Num activations matches:", cnt)

	return cnt, nil
}

// Gets only running activations (activation_running meaning that users are using
// the machine/resource currently)
func GetActiveActivations() ([]*Activation, error) {
	var purchases []*Purchase
	o := orm.NewOrm()
	act := Activation{}
	num, err := o.QueryTable(act.Purchase.TableName()).
		Filter("running", true).
		Filter("type", TYPE_ACTIVATION).
		All(&purchases)
	if err != nil {
		return nil, fmt.Errorf("Failed to get active activations: %v", err)
	}
	beego.Trace("Got num activations:", num)

	activations := make([]*Activation, 0, len(purchases))
	for _, purchase := range purchases {
		a := &Activation{
			Purchase: *purchase,
		}
		a.Purchase.Quantity = a.Purchase.quantityFromTimes()
		activations = append(activations, a)
	}

	return activations, nil
}

// Creates activation and returns activation ID.
func CreateActivation(machineId, userId int64, startTime time.Time) (
	activationId int64, err error) {

	o := orm.NewOrm()
	mch := models.Machine{Id: machineId}

	if !mch.Exists() {
		activationId = 0
		err = fmt.Errorf("Machine with provided ID does not exist")
		return
	}

	// Check for duplicate activations
	// TODO: Replace this with a more readable helper function
	var dupActivations []Activation
	act := Activation{} // Used to get table name of the model
	query := fmt.Sprintf("SELECT id FROM %s WHERE machine_id = ? "+
		"AND user_id = ? AND running = 1 AND type = ?", act.Purchase.TableName())
	numDuplicates, err := o.Raw(query, machineId, userId, TYPE_ACTIVATION).
		QueryRows(&dupActivations)
	if err != nil {
		return 0, fmt.Errorf("Could not get duplicate activations: %v", err)
	}
	if numDuplicates > 0 {
		beego.Error("Duplicate activations found:", numDuplicates)
		return 0, fmt.Errorf("Duplicate activations found")
	}

	if !mch.IsAvailable() {
		activationId = 0
		err = fmt.Errorf("Machine with provided ID is not available")
		return
	}

	err, _ = mch.Read()
	if err != nil {
		activationId = 0
		err = fmt.Errorf("Failed to read existing machine")
		return
	}

	newActivation := Activation{
		Purchase: Purchase{
			Type:      TYPE_ACTIVATION,
			UserId:    userId,
			MachineId: machineId,
			Running:   true,
			TimeStart: startTime,

			// Save current activation price, currency and price unit (minute, hour, pcs)
			PricePerUnit: mch.Price,
			PriceUnit:    mch.PriceUnit,
		},
	}

	activationId, err = o.Insert(&newActivation.Purchase)
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

// Gets pointer to activation store by activation ID.
func GetActivation(activationId int64) (activation *Activation, err error) {
	activation = &Activation{}
	activation.Purchase.Id = activationId

	o := orm.NewOrm()
	err = o.Read(&activation.Purchase)

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
	activation.Purchase.Running = false
	activation.Purchase.TimeEnd = endTime
	activation.Purchase.Quantity = activation.Purchase.quantityFromTimes()

	err = activation.Update()
	if err != nil {
		beego.Error("Failed to update activation:", err)
		return fmt.Errorf("Failed to update activation: %v", err)
	}

	// Make the machine available again.
	var machine *models.Machine
	machine, err = models.GetMachine(activation.Purchase.MachineId)
	if err != nil {
		beego.Error("Failed to get machine:", err)
		return fmt.Errorf("Failed to get machine: %v", err)
	}

	machine.Available = true
	err = machine.Update()
	if err != nil {
		beego.Error("Failed to update machine:", err)
		return fmt.Errorf("Failed to update machine: %v", err)
	}

	return nil
}

// Updates existing activation by consuming a pointer to
// existing activation store.
func (activation *Activation) Update() error {
	o := orm.NewOrm()
	num, err := o.Update(&activation.Purchase)

	if err != nil {
		beego.Error("Failed to update activation:", err)
		return fmt.Errorf("Failed to update activation: %v", err)
	}

	beego.Trace("UpdateActivation: Affected num rows:", num)

	return nil
}

// Gets the machine ID of a specific activation defined by activation ID.
func GetActivationMachineId(activationId int64) (int64, error) {
	activationModel := Activation{}
	o := orm.NewOrm()
	err := o.QueryTable(activationModel.Purchase.TableName()).
		Filter("id", activationId).
		Filter("type", TYPE_ACTIVATION).
		One(&activationModel.Purchase, "MachineId")
	if err != nil {
		beego.Error("Could not get activation")
		return 0, err
	}
	return activationModel.Purchase.MachineId, nil
}
