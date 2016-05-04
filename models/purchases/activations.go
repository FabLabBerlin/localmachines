package purchases

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type Activation struct {
	json.Marshaler
	Purchase Purchase
}

func (this *Activation) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Purchase)
}

func CreateActivation(locationId int64) (activation Activation, err error) {
	o := orm.NewOrm()
	if locationId <= 0 {
		return activation, fmt.Errorf("invalid location id: %v", locationId)
	}
	activation.Purchase.Type = TYPE_ACTIVATION
	activation.Purchase.LocationId = locationId
	activation.Purchase.Id, err = o.Insert(&activation.Purchase)
	return
}

// Gets filtered activations in a paged manner between start and end time.
// Items per page and page number can be specified. Already invoiced
// activations can be excluded.
func GetActivations(locationId int64, interval lib.Interval, search string) (activations []Activation, err error) {
	// Get activations from database
	var purchases []*Purchase
	act := Activation{}
	o := orm.NewOrm()

	if pattern, ok := searchTermToPattern(search); ok {
		query := fmt.Sprintf("SELECT a.* FROM %s a "+
			"LEFT JOIN user u ON a.user_id = u.id "+
			"WHERE a.type=? AND a.time_start>=? AND a.time_start<=? AND a.running=false "+
			"      AND (u.username LIKE ? OR u.first_name LIKE ? OR u.last_name LIKE ?) "+
			"      AND a.location_id=? "+
			"ORDER BY a.time_start DESC ",
			act.Purchase.TableName())

		_, err = o.Raw(query,
			TYPE_ACTIVATION,
			interval.TimeFrom(),
			interval.TimeTo(),
			pattern,
			pattern,
			pattern,
			locationId).QueryRows(&purchases)
	} else {
		query := fmt.Sprintf("SELECT a.* FROM %s a "+
			"WHERE a.type=? AND a.time_start>=? AND a.time_start<=? AND a.running=false "+
			"               AND a.location_id=? "+
			"ORDER BY a.time_start DESC ",
			act.Purchase.TableName())

		_, err = o.Raw(query,
			TYPE_ACTIVATION,
			interval.TimeFrom(),
			interval.TimeTo(),
			locationId).QueryRows(&purchases)

	}

	if err != nil {
		return nil, fmt.Errorf("Failed to get activations: %v", err)
	}

	activations = make([]Activation, 0, len(purchases))
	for _, purchase := range purchases {
		act := Activation{
			Purchase: *purchase,
		}
		activations = append(activations, act)
	}

	return
}

func searchTermToPattern(search string) (pattern string, ok bool) {
	search = strings.TrimSpace(search)
	if search != "" {
		return "%" + search + "%", true
	} else {
		return "", false
	}
}

// Gets only running activations (activation_running meaning that users are using
// the machine/resource currently)
func GetActiveActivations() ([]*Activation, error) {
	var purchases []*Purchase
	o := orm.NewOrm()
	act := Activation{}
	_, err := o.QueryTable(act.Purchase.TableName()).
		Filter("running", true).
		Filter("type", TYPE_ACTIVATION).
		All(&purchases)
	if err != nil {
		return nil, fmt.Errorf("Failed to get active activations: %v", err)
	}

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
func StartActivation(machineId, userId int64, startTime time.Time) (
	activationId int64, err error) {

	o := orm.NewOrm()

	mch, err := machine.Get(machineId)
	if err != nil {
		return 0, fmt.Errorf("get machine: %v", err)
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

	newActivation := Activation{
		Purchase: Purchase{
			LocationId: mch.LocationId,
			Type:       TYPE_ACTIVATION,
			UserId:     userId,
			MachineId:  machineId,
			Running:    true,
			TimeStart:  startTime,

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

	if err := redis.PublishMachinesUpdate(mch.LocationId); err != nil {
		beego.Error("publish machines update:", err)
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
		return nil, fmt.Errorf("Failed to read activation: %v", err)
	}

	return
}

// Close running/active activation.
func (activation *Activation) Close(endTime time.Time) error {
	machine, err := machine.Get(activation.Purchase.MachineId)
	if err != nil {
		beego.Error("Failed to get machine:", err)
		return fmt.Errorf("Failed to get machine: %v", err)
	}

	gracePeriod := machine.GetGracePeriod()

	// Calculate activation duration and update activation.
	activation.Purchase.Running = false
	activation.Purchase.TimeStart = activation.Purchase.TimeStart.Add(gracePeriod)
	activation.Purchase.TimeEnd = endTime
	if activation.Purchase.TimeEnd.Before(activation.Purchase.TimeStart) {
		activation.Purchase.TimeEnd = activation.Purchase.TimeStart
		if gracePeriod == 0 {
			beego.Error("time end before time start?!")
		}
	}

	err = activation.Update()
	if err != nil {
		return fmt.Errorf("Failed to update activation: %v", err)
	}

	// Make the machine available again.
	machine.Available = true
	if err = machine.Update(false); err != nil {
		beego.Error("Failed to update machine:", err)
		return fmt.Errorf("Failed to update machine: %v", err)
	}

	if err := redis.PublishMachinesUpdate(machine.LocationId); err != nil {
		beego.Error("publish machines update:", err)
	}

	return nil
}

// Updates existing activation by consuming a pointer to
// existing activation store.
func (activation *Activation) Update() error {
	o := orm.NewOrm()

	activation.Purchase.Quantity = activation.Purchase.quantityFromTimes()

	_, err := o.Update(&activation.Purchase)

	if err != nil {
		beego.Error("Failed to update activation:", err)
		return fmt.Errorf("Failed to update activation: %v", err)
	}

	return nil
}
