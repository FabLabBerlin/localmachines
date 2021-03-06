package purchases

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/products"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Tutoring struct {
	json.Marshaler
	json.Unmarshaler
	Purchase
	//TutorId    int64 - Product Id?!
}

func (this *Tutoring) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Purchase)
}

func (this *Tutoring) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.Purchase)
}

type TutoringList struct {
	Data []*Tutoring
}

// Create a planned tutoring reservation in the future with blank values.
func CreateTutoring(tp *Tutoring) (id int64, err error) {
	tp.Created = time.Now()
	tp.Type = TYPE_TUTOR
	tp.TimeStart = time.Now()
	//tp.TimeEndPlanned = time.Now()
	tp.PriceUnit = "hour"

	err = Create(&tp.Purchase)
	id = tp.Purchase.Id

	return
}

func GetTutoring(id int64) (tutoring *Tutoring, err error) {
	tutoring = &Tutoring{}
	tutoring.Id = id

	o := orm.NewOrm()
	err = o.Read(&tutoring.Purchase)

	return
}

// Get a list of tutoring purchases
func GetAllTutorings() (tutorings *TutoringList, err error) {
	purchases, err := GetAllOfType(TYPE_TUTOR)
	if err != nil {
		return
	}
	tutorings = &TutoringList{
		Data: make([]*Tutoring, 0, len(purchases)),
	}
	for _, purchase := range purchases {
		t := &Tutoring{
			Purchase: *purchase,
		}
		tutorings.Data = append(tutorings.Data, t)
	}
	return
}

// Get a list of tutoring purchases
func GetAllTutoringsAt(locationId int64) (tutorings *TutoringList, err error) {
	purchases, err := GetAllOfTypeAt(locationId, TYPE_TUTOR)
	if err != nil {
		return
	}
	tutorings = &TutoringList{
		Data: make([]*Tutoring, 0, len(purchases)),
	}
	for _, purchase := range purchases {
		t := &Tutoring{
			Purchase: *purchase,
		}
		tutorings.Data = append(tutorings.Data, t)
	}
	return
}

// Start the tutoring purchase timer.
func StartTutoring(id int64) (err error) {
	o := orm.NewOrm()

	// Let's just use full names by convention as code is being copy/pasted
	// and there is too much human error involved. At the end of the day
	// the variable names just do not make sense.
	tp, err := GetTutoring(id)
	if err != nil {
		return fmt.Errorf("Failed to get tutoring purchase: %v", err)
	}

	// Mark the timer as running and store current time
	// so we can make calculation when we stop the timer.
	tp.Running = true
	tp.TimerTimeStart = time.Now()
	_, err = o.Update(&tp.Purchase)
	if err != nil {
		return fmt.Errorf("Failed to update: %v", err)
	}

	if err := redis.PublishMachinesUpdate(redis.MachinesUpdate{
		LocationId: tp.LocationId,
	}); err != nil {
		beego.Error("tutoring start: publish machines update:", err)
	}

	return
}

// Stop tutoring purchase timer.
func StopTutoring(id int64) (err error) {
	tp, err := GetTutoring(id)
	if err != nil {
		return fmt.Errorf("get tutoring purchase: %v", err)
	}

	// Get existing tutoring
	if tp.ProductId > 0 {
		tutor, err := products.GetTutor(tp.ProductId)
		if err != nil {
			return fmt.Errorf("get tutor: %v", err)
		}
		tp.PricePerUnit = tutor.Product.Price
		tp.PriceUnit = tutor.Product.PriceUnit // this should be hour
	}

	// OK, now we have to deal with the quantity of the tutoring
	// purchase we are stoping.
	var lastTimerSessionQuantity float64
	switch tp.PriceUnit {
	case "hour":
		lastTimerSessionQuantity = time.
			Since(tp.TimerTimeStart).Hours()
		break
	case "minute":
		lastTimerSessionQuantity = time.
			Since(tp.TimerTimeStart).Minutes()
		break
	case "second":
		lastTimerSessionQuantity = time.
			Since(tp.TimerTimeStart).Seconds()
		break
	default:
		lastTimerSessionQuantity = time.
			Since(tp.TimerTimeStart).Hours()
	}

	// Add up elapsed time since last timer start time to the existing quantity
	tp.Purchase.Quantity += lastTimerSessionQuantity

	// The timer is not running anymore, so we set the running flag to false.
	tp.Purchase.Running = false

	o := orm.NewOrm()
	if _, err = o.Update(&tp.Purchase); err != nil {
		return fmt.Errorf("update: %v", err)
	}

	if err := redis.PublishMachinesUpdate(redis.MachinesUpdate{
		LocationId: tp.LocationId,
	}); err != nil {
		beego.Error("tutoring stop: publish machines update:", err)
	}

	return
}

func (tp *Tutoring) Update() (err error) {
	return Update(&tp.Purchase)
}
