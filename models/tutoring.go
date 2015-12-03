package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

///////////////////////
//      Tutors       //
///////////////////////

type Tutor struct {
	Product Product
}

// Get a list of tutors
func GetAllTutors() ([]*Tutor, error) {
	tutorsAsProducts, err := GetAllProductsOfType(PRODUCT_TYPE_TUTOR)
	if err != nil {
		msg := "Failed to get tutors as products"
		beego.Error(msg)
		return nil, fmt.Errorf(msg+"%v: ", err)
	}

	tutorList := make([]*Tutor, 0, len(tutorsAsProducts))
	for _, product := range tutorsAsProducts {
		tutor := &Tutor{
			Product: *product,
		}
		tutorList = append(tutorList, tutor)
	}

	return tutorList, nil
}

func CreateTutor(tutor *Tutor) (*Tutor, error) {
	o := orm.NewOrm()

	// In case the type has not been added in upper layers
	tutor.Product.Type = PRODUCT_TYPE_TUTOR

	var productId int64
	productId, err := o.Insert(&tutor.Product)
	if err != nil {
		msg := "Failed to insert tutor"
		beego.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	tutor.Product.Id = productId

	return tutor, nil
}

func GetTutor(id int64) (tutor *Tutor, err error) {
	tutor = &Tutor{}
	tutor.Product.Id = id

	o := orm.NewOrm()
	err = o.Read(&tutor.Product)

	return
}

func UpdateTutor(tutor *Tutor) error {
	if tutor.Product.UserId != 0 {
		o := orm.NewOrm()
		// Get user name by user ID
		user := User{}
		user.Id = tutor.Product.UserId
		err := o.Read(&user)
		if err != nil {
			beego.Error("Failed to read user:", err)
			return fmt.Errorf("Failed to read user: %v", err)
		}
		tutor.Product.Name = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}
	return UpdateProduct(&tutor.Product)
}

///////////////////////
//     Tutorings     //
///////////////////////

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
	o := orm.NewOrm()
	t := new(TutoringPurchase)
	_, err = o.QueryTable(t.Purchase.TableName()).
		Filter("id", tutoringPurchaseId).
		Update(orm.Params{
		"quantity": t.Purchase.quantityFromTimes(),
		"running":  false,
		"time_end": time.Now(),
	})
	return
}

func UpdateTutoringPurchase(tutoringPurchase *TutoringPurchase) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&tutoringPurchase.Purchase)
	return
}
