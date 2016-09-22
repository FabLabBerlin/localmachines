package purchases

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/easylab-lib/xmpp/commands"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
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

// Gets filtered activations in a paged manner between start and end time.
// Items per page and page number can be specified. Already invoiced
// activations can be excluded.
func GetActivations(locationId int64, interval lib.Interval, search string) (activations []Activation, err error) {
	beego.Info("GetActivations(.................)")
	// Get activations from database
	var purchases []*Purchase
	act := Activation{}
	o := orm.NewOrm()

	if pattern, ok := searchTermToPattern(search); ok {
		query := fmt.Sprintf("SELECT a.* FROM %s a "+
			"LEFT JOIN user u ON a.user_id = u.id "+
			"WHERE a.type=? AND a.time_start>=? AND a.time_start<=? AND a.running=false "+
			"      AND (u.username LIKE ? OR u.first_name LIKE ? OR u.last_name LIKE ?) "+
			"      AND ((invoice_status IS NULL) OR invoice_status = '' OR invoice_status = 'draft')"+
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
			"               AND ((invoice_status IS NULL) OR invoice_status = '' OR invoice_status = 'draft')"+
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
func StartActivation(m *machine.Machine, uid int64, start time.Time) (
	activationId int64, err error) {

	o := orm.NewOrm()

	// Check for duplicate activations
	// TODO: Replace this with a more readable helper function
	var dupActivations []*Purchase
	numDuplicates, err := o.QueryTable(TABLE_NAME).
		Filter("machine_id", m.Id).
		Filter("user_id", uid).
		Filter("running", 1).
		Filter("type", TYPE_ACTIVATION).
		All(&dupActivations)
	if err != nil {
		return 0, fmt.Errorf("Could not get duplicate activations: %v", err)
	}
	if numDuplicates == 1 {
		return dupActivations[0].Id, nil
	} else if numDuplicates > 1 {
		beego.Error("Duplicate activations found:", numDuplicates)
		return 0, fmt.Errorf("Duplicate activations found")
	}

	if !m.IsAvailable() {
		activationId = 0
		err = fmt.Errorf("Machine with provided ID is not available")
		return
	}

	inv, err := invoices.GetDraft(m.LocationId, uid, time.Now())
	if err != nil {
		return 0, fmt.Errorf("current invoice: %v", err)
	}

	newActivation := Activation{
		Purchase: Purchase{
			LocationId: m.LocationId,
			Type:       TYPE_ACTIVATION,
			UserId:     uid,
			MachineId:  m.Id,
			Running:    true,
			TimeStart:  start,
			InvoiceId:  inv.Id,

			// Save current activation price, currency and price unit (minute, hour, pcs)
			PricePerUnit: m.Price,
			PriceUnit:    m.PriceUnit,
		},
	}

	if err = Create(&newActivation.Purchase); err != nil {
		beego.Error("Failed to insert activation:", err)
		return 0, fmt.Errorf("Failed to insert activation %v", err)
	}

	// Update machine as unavailable
	m.Available = false
	if m.Update(false); err != nil {
		beego.Error("Failed to update activated machine")
	}

	return newActivation.Purchase.Id, nil
}

// Gets pointer to activation store by activation ID.
func GetActivation(id int64) (a *Activation, err error) {
	a = &Activation{}
	a.Purchase.Id = id

	err = orm.NewOrm().Read(&a.Purchase)

	return
}

// Close running/active activation.
func (a *Activation) Close(endTime time.Time) error {
	m, err := machine.Get(a.Purchase.MachineId)
	if err != nil {
		beego.Error("Failed to get machine:", err)
		return fmt.Errorf("Failed to get machine: %v", err)
	}

	gracePeriod := m.GetGracePeriod()

	// Calculate activation duration and update activation.
	a.Purchase.Running = false
	a.Purchase.TimeStart = a.Purchase.TimeStart.Add(gracePeriod)
	a.Purchase.TimeEnd = endTime
	if a.Purchase.TimeEnd.Before(a.Purchase.TimeStart) {
		a.Purchase.TimeEnd = a.Purchase.TimeStart
		if gracePeriod == 0 {
			beego.Error("time end before time start?!")
		}
	}

	if err = a.Update(); err != nil {
		return fmt.Errorf("Failed to update activation: %v", err)
	}

	m.Available = true
	if err = m.Update(false); err != nil {
		return fmt.Errorf("Failed to update machine: %v", err)
	}

	if err := redis.PublishMachinesUpdate(redis.MachinesUpdate{
		LocationId: m.LocationId,
	}); err != nil {
		beego.Error("publish machines update:", err)
	}

	return nil
}

// Updates existing activation by consuming a pointer to
// existing activation store.
func (a *Activation) Update() error {
	if a.Purchase.InvoiceId <= 0 {
		return fmt.Errorf("undefined invoice id")
	}

	if mid := a.Purchase.MachineId; mid != 0 {
		m, err := machine.Get(mid)
		if err != nil {
			return err
		}
		a.Purchase.PriceUnit = m.PriceUnit
		a.Purchase.PricePerUnit = m.Price
	}
	a.Purchase.Quantity = floor10(a.Purchase.quantityFromTimes())

	return Update(&a.Purchase)
}

func floor10(x float64) float64 {
	return math.Floor(x*10) / 10
}

func applyToBilling(update redis.MachinesUpdate) {
	switch update.Command {
	case commands.GATEWAY_SUCCESS_ON:
		// Continue with creating activation
		startTime := time.Now()
		m, err := machine.Get(update.MachineId)
		if err != nil {
			beego.Error(err.Error())
			return
		}
		if _, err = StartActivation(
			m,
			update.UserId,
			startTime,
		); err != nil {
			beego.Error("Failed to create activation:", err)
		}
	case commands.GATEWAY_SUCCESS_OFF:
		as, err := GetActiveActivations()
		if err != nil {
			beego.Error(err.Error())
			return
		}
		for _, a := range as {
			if a.Purchase.LocationId == update.LocationId &&
				a.Purchase.MachineId == update.MachineId {
				if err = a.Close(time.Now()); err != nil {
					beego.Error(err.Error())
				}
			}
		}
	}
}

func init() {
	machine.ApplyToBilling = applyToBilling
}
