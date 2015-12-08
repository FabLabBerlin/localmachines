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

func CreateTutoring(tutoring *Tutoring) (int64, error) {
	tutoring.Purchase.Created = time.Now()
	tutoring.Purchase.Type = TYPE_TUTOR
	tutoring.Purchase.TimeStart = time.Now()
	tutoring.Purchase.TimeEndPlanned = time.Now()
	tutoring.Purchase.PriceUnit = "hour"

	o := orm.NewOrm()
	return o.Insert(&tutoring.Purchase)
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

func StartTutoring(id int64) (err error) {
	o := orm.NewOrm()
	tp, err := GetTutoring(id)
	if err != nil {
		return fmt.Errorf("exists: %v", err)
	}
	if tp.TimeEnd.IsZero() {
		t := new(Tutoring)
		_, err = o.QueryTable(t.Purchase.TableName()).
			Filter("id", id).
			Update(orm.Params{
			"running":    true,
			"time_start": time.Now(),
		})
	} else {
		newTp := *tp
		newTp.Id = 0
		newTp.UserId = tp.UserId
		newTp.ProductId = tp.ProductId
		newTp.Running = true
		newTp.TimeStart = time.Now()
		newTp.TimeEndPlanned = tp.TimeEndPlanned
		newTp.TimeEnd = time.Time{}
		newTp.Quantity = 0
		_, err = CreateTutoring(&newTp)
	}
	return

}

func StopTutoring(id int64) (err error) {
	t, err := GetTutoring(id)
	if err != nil {
		return fmt.Errorf("get tutoring purchase: %v", err)
	}
	t.Purchase.Quantity = t.Purchase.quantityFromTimes()
	t.Purchase.Running = false
	t.Purchase.TimeEnd = time.Now()
	return t.Update()
}

func (tutoring *Tutoring) Update() (err error) {
	o := orm.NewOrm()
	if tutoring.ProductId > 0 {
		tutor, err := products.GetTutor(tutoring.ProductId)
		if err != nil {
			return fmt.Errorf("get tutor: %v", err)
		}
		tutoring.PricePerUnit = tutor.Product.Price
		tutoring.PriceUnit = tutor.Product.PriceUnit
	}
	tutoring.Quantity = tutoring.quantityFromTimes()
	_, err = o.Update(&tutoring.Purchase)
	return
}
