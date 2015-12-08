package purchases

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models/products"
	"time"
)

type TutoringPurchase struct {
	json.Marshaler
	json.Unmarshaler
	Purchase
	//TutorId    int64 - Product Id?!
}

func (this *TutoringPurchase) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Purchase)
}

func (this *TutoringPurchase) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.Purchase)
}

type TutoringPurchaseList struct {
	Data []*TutoringPurchase
}

func CreateTutoringPurchase(tutoringPurchase *TutoringPurchase) (int64, error) {
	tutoringPurchase.Purchase.Created = time.Now()
	tutoringPurchase.Purchase.Type = PURCHASE_TYPE_TUTOR
	tutoringPurchase.Purchase.TimeStart = time.Now()
	tutoringPurchase.Purchase.TimeEndPlanned = time.Now()
	tutoringPurchase.Purchase.PriceUnit = "hour"

	o := orm.NewOrm()
	return o.Insert(&tutoringPurchase.Purchase)
}

func GetTutoringPurchase(id int64) (tutoringPurchase *TutoringPurchase, err error) {
	tutoringPurchase = &TutoringPurchase{}
	tutoringPurchase.Purchase.Id = id

	o := orm.NewOrm()
	err = o.Read(&tutoringPurchase.Purchase)

	return
}

// Get a list of tutoring purchases
func GetAllTutoringPurchases() (tutoringPurchases *TutoringPurchaseList, err error) {
	purchases, err := GetAllPurchasesOfType(PURCHASE_TYPE_TUTOR)
	if err != nil {
		return
	}
	tutoringPurchases = &TutoringPurchaseList{
		Data: make([]*TutoringPurchase, 0, len(purchases)),
	}
	for _, purchase := range purchases {
		tutoringPurchase := &TutoringPurchase{
			Purchase: *purchase,
		}
		tutoringPurchases.Data = append(tutoringPurchases.Data, tutoringPurchase)
	}
	return
}

func StartTutoringPurchase(tutoringPurchaseId int64) (err error) {
	o := orm.NewOrm()
	tp, err := GetTutoringPurchase(tutoringPurchaseId)
	if err != nil {
		return fmt.Errorf("exists: %v", err)
	}
	if tp.TimeEnd.IsZero() {
		t := new(TutoringPurchase)
		_, err = o.QueryTable(t.Purchase.TableName()).
			Filter("id", tutoringPurchaseId).
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
		_, err = CreateTutoringPurchase(&newTp)
	}
	return

}

func StopTutoringPurchase(tutoringPurchaseId int64) (err error) {
	t, err := GetTutoringPurchase(tutoringPurchaseId)
	if err != nil {
		return fmt.Errorf("get tutoring purchase: %v", err)
	}
	t.Purchase.Quantity = t.Purchase.quantityFromTimes()
	t.Purchase.Running = false
	t.Purchase.TimeEnd = time.Now()
	return UpdateTutoringPurchase(t)
}

func UpdateTutoringPurchase(tutoringPurchase *TutoringPurchase) (err error) {
	o := orm.NewOrm()
	if tutoringPurchase.ProductId > 0 {
		tutor, err := products.GetTutor(tutoringPurchase.ProductId)
		if err != nil {
			return fmt.Errorf("get tutor: %v", err)
		}
		tutoringPurchase.PricePerUnit = tutor.Product.Price
		tutoringPurchase.PriceUnit = tutor.Product.PriceUnit
	}
	tutoringPurchase.Quantity = tutoringPurchase.quantityFromTimes()
	_, err = o.Update(&tutoringPurchase.Purchase)
	return
}
