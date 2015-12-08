package purchases

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models/products"
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
func CreateTutoring(tutoringPurchase *Tutoring) (int64, error) {
	tutoringPurchase.Purchase.Created = time.Now()
	tutoringPurchase.Purchase.Type = TYPE_TUTOR
	tutoringPurchase.Purchase.TimeStart = time.Now()
	tutoringPurchase.Purchase.TimeEnd = tutoringPurchase.Purchase.TimeStart.
		Add(time.Duration(1) * time.Hour)
	//tutoringPurchase.Purchase.TimeEndPlanned = time.Now()
	tutoringPurchase.Purchase.PriceUnit = "hour"

	o := orm.NewOrm()
	return o.Insert(&tutoringPurchase.Purchase)
}

func GetTutoring(id int64) (tutoring *Tutoring, err error) {
	tutoring = &Tutoring{}
	tutoring.Purchase.Id = id

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

// Start the tutoring purchase timer.
func StartTutoring(tutoringPurchaseId int64) (err error) {
	o := orm.NewOrm()

	// Let's just use full names by convention as code is being copy/pasted
	// and there is too much human error involved. At the end of the day
	// the variable names just do not make sense.
	tutoringPurchase, err := GetTutoring(tutoringPurchaseId)
	if err != nil {
		return fmt.Errorf("Failed to get tutoring purchase: %v", err)
	}

	// Mark the timer as running and store current time
	// so we can make calculation when we stop the timer.
	tutoringPurchase.Purchase.Running = true
	tutoringPurchase.Purchase.TimerTimeStart = time.Now()
	_, err = o.Update(&tutoringPurchase.Purchase)
	if err != nil {
		return fmt.Errorf("Failed to update: %v", err)
	}

	return
}

// Stop tutoring purchase timer.
func StopTutoringPurchase(tutoringPurchaseId int64) (err error) {
	tutoringPurchase, err := GetTutoring(tutoringPurchaseId)
	if err != nil {
		return fmt.Errorf("get tutoring purchase: %v", err)
	}

	// Get existing tutoring
	if tutoringPurchase.ProductId > 0 {
		tutor, err := products.GetTutor(tutoringPurchase.ProductId)
		if err != nil {
			return fmt.Errorf("get tutor: %v", err)
		}
		tutoringPurchase.PricePerUnit = tutor.Product.Price
		tutoringPurchase.PriceUnit = tutor.Product.PriceUnit // this should be hour
	}

	// OK, now we have to deal with the quantity of the tutoring
	// purchase we are stoping.
	var lastTimerSessionQuantity float64
	switch tutoringPurchase.PriceUnit {
	case "hour":
		lastTimerSessionQuantity = time.
			Since(tutoringPurchase.TimerTimeStart).Hours()
		break
	case "minute":
		lastTimerSessionQuantity = time.
			Since(tutoringPurchase.TimerTimeStart).Minutes()
		break
	case "second":
		lastTimerSessionQuantity = time.
			Since(tutoringPurchase.TimerTimeStart).Seconds()
		break
	default:
		lastTimerSessionQuantity = time.
			Since(tutoringPurchase.TimerTimeStart).Hours()
	}

	// Add up elapsed time since last timer start time to the existing quantity
	tutoringPurchase.Purchase.Quantity += lastTimerSessionQuantity

	// The timer is not running anymore, so we set the running flag to false.
	tutoringPurchase.Purchase.Running = false

	o := orm.NewOrm()
	_, err = o.Update(&tutoringPurchase.Purchase)

	return
}

func UpdateTutoringPurchase(tutoringPurchase *Tutoring) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&tutoringPurchase.Purchase)
	return
}
