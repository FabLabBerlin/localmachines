package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Activation struct {
	json.Marshaler
	purchase Purchase
}

func (this *Activation) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.purchase)
}

/*// Activation type/model to hold information of single activation.
type Activation struct {
	Id int64 `orm:"auto";"pk"`
	//[unused]InvoiceId                   int   `orm:"null"`
	UserId    int64
	MachineId int64
	Active    bool
	TimeStart time.Time
	TimeEnd   time.Time `orm:"null"`
	TimeTotal int64
	//[unused]UsedKwh                     float32
	//[unused]DiscountPercents            float32
	//[unused]DiscountFixed               float32
	//[unused]VatRate                     float32
	//[unused]CommentRef                  string `orm:"size(255)"`
	//[unused]Invoiced                    bool
	//[unused]Changed                     bool
	//CurrentMachinePrice         float64
	//CurrentMachinePriceCurrency string
	//CurrentMachinePriceUnit     string
}

// Returns mysql table name of the table mapped to the Activation model.
func (this *Activation) TableName() string {
	return "activations"
}

func init() {
	orm.RegisterModel(new(Activation))
}*/

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
		"WHERE a.type=? AND a.time_start>? AND a.time_end<? AND a.activation_running=false AND a.user_id=? "+
		"ORDER BY u.first_name ASC, a.time_start DESC",
		act.purchase.TableName(),
		usr.TableName())

	_, err := o.Raw(query,
		PURCHASE_TYPE_ACTIVATION,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
		userId).QueryRows(&activations)

	if err != nil {
		msg := fmt.Sprintf("Failed to get activations: %v", err)
		return nil, errors.New(msg)
	}

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
	itemsPerPage int64,
	page int64) (*[]Activation, error) {

	if page <= 0 {
		page = 1
	}

	// Get activations from database
	purchases := []*Purchase{}
	act := Activation{}
	usr := User{}
	o := orm.NewOrm()

	var pageOffset int64
	pageOffset = itemsPerPage * (page - 1)

	query := fmt.Sprintf("SELECT a.* FROM %s a JOIN %s u ON a.user_id=u.id "+
		"WHERE a.type=? AND a.time_start>? AND a.time_end<? AND a.activation_running=false "+
		"ORDER BY u.first_name ASC, a.time_start DESC "+
		"LIMIT ? OFFSET ?",
		act.purchase.TableName(),
		usr.TableName())

	_, err := o.Raw(query,
		PURCHASE_TYPE_ACTIVATION,
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
			purchase: *purchase,
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
	cnt, err := o.QueryTable(act.purchase.TableName()).
		Filter("timeStart__gt", startTime).
		Filter("timeEnd__lt", endTime).
		//Filter("userId", userId).
		Filter("ActivationRunning", false).
		Filter("type", PURCHASE_TYPE_ACTIVATION).
		Count()

	if err != nil {
		msg := fmt.Sprintf("Failed to count activations: %v", err)
		return 0, errors.New(msg)
	}

	beego.Trace("Num activations matches:", cnt)

	return cnt, nil
}

// Gets only activation_running activations (activation_running meaning that users are using
// the machine/resource currently)
func GetActiveActivations() ([]*Activation, error) {
	var purchases []*Purchase
	o := orm.NewOrm()
	act := Activation{}
	num, err := o.QueryTable(act.purchase.TableName()).
		Filter("activation_running", true).
		Filter("type", PURCHASE_TYPE_ACTIVATION).
		All(&purchases)
	if err != nil {
		return nil, errors.New("Failed to get active activations")
	}
	beego.Trace("Got num activations:", num)

	activations := make([]*Activation, 0, len(purchases))
	for _, purchase := range purchases {
		a := &Activation{
			purchase: *purchase,
		}
		timeNow := time.Now()
		a.purchase.Quantity =
			float64(timeNow.Sub(purchase.TimeStart).Seconds())
		activations = append(activations, a)
	}

	return activations, nil
}

// Creates activation and returns activation ID.
func CreateActivation(machineId, userId int64, startTime time.Time) (
	activationId int64, err error) {

	o := orm.NewOrm()
	mch := Machine{Id: machineId}

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
		"AND user_id = ? AND activation_running = 1 AND type = ?", act.purchase.TableName())
	numDuplicates, err := o.Raw(query, machineId, userId, PURCHASE_TYPE_ACTIVATION).
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

	newActivation := Activation{}
	newActivation.purchase.Type = PURCHASE_TYPE_ACTIVATION
	newActivation.purchase.UserId = userId
	newActivation.purchase.MachineId = machineId
	newActivation.purchase.ActivationRunning = true
	newActivation.purchase.TimeStart = startTime

	// Save current activation price, currency and price unit (minute, hour, pcs)
	newActivation.purchase.PricePerUnit = mch.Price
	newActivation.purchase.PriceUnit = mch.PriceUnit

	activationId, err = o.Insert(&newActivation.purchase)
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
	activation.purchase.Id = activationId

	o := orm.NewOrm()
	err = o.Read(&activation.purchase)

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
	activation.purchase.ActivationRunning = false
	activation.purchase.TimeEnd = endTime
	activation.purchase.Quantity = float64(endTime.Sub(activation.purchase.TimeStart).Seconds())

	err = UpdateActivation(activation)
	if err != nil {
		beego.Error("Failed to update activation:", err)
		return fmt.Errorf("Failed to update activation: %v", err)
	}

	// Make the machine available again.
	var machine *Machine
	machine, err = GetMachine(activation.purchase.MachineId)
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

// Updates existing activation by consuming a pointer to
// existing activation store.
func UpdateActivation(activation *Activation) error {
	o := orm.NewOrm()
	num, err := o.Update(&activation.purchase)

	if err != nil {
		beego.Error("Failed to update activation:", err)
		return fmt.Errorf("Failed to update activation: %v", err)
	}

	beego.Trace("UpdateActivation: Affected num rows:", num)

	return nil
}

// Delete activation matching an activation ID.
func DeleteActivation(activationId int64) error {

	// Set machine of the activation available
	var err error
	var activation Activation
	o := orm.NewOrm()
	err = o.QueryTable(activation.purchase.TableName()).
		Filter("Id", activationId).
		Filter("type", PURCHASE_TYPE_ACTIVATION).
		One(&activation.purchase, "MachineId")
	if err != nil {
		beego.Error("Failed to get machine ID of the activation")
		return err
	}
	m := Machine{}
	_, err = o.QueryTable(m.TableName()).
		Filter("Id", activation.purchase.MachineId).
		Update(orm.Params{"available": true})
	if err != nil {
		beego.Error("Failed to update machine as available")
		return err
	}

	return DeletePurchase(activation.purchase.Id)
}

// Gets the machine ID of a specific activation defined by activation ID.
func GetActivationMachineId(activationId int64) (int64, error) {
	activationModel := Activation{}
	o := orm.NewOrm()
	err := o.QueryTable(activationModel.purchase.TableName()).
		Filter("id", activationId).
		Filter("type", PURCHASE_TYPE_ACTIVATION).
		One(&activationModel.purchase, "MachineId")
	if err != nil {
		beego.Error("Could not get activation")
		return 0, err
	}
	return activationModel.purchase.MachineId, nil
}
