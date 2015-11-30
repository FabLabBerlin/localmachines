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

type TutorList struct {
	Data []*Product
}

// Get a list of tutors
func GetAllTutors() (*TutorList, error) {
	tutorsAsProducts, err := GetAllProductsOfType(PRODUCT_TYPE_TUTOR)
	if err != nil {
		msg := "Failed to get tutors as products"
		beego.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	tutorList := TutorList{Data: tutorsAsProducts}

	return &tutorList, nil
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
	tutoringPurchase.Purchase = Purchase{
		Created:   time.Now(),
		Type:      PURCHASE_TYPE_TUTOR,
		TimeStart: time.Now(),
		TimeEnd:   time.Now(),
		PriceUnit: "hour",
	}

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

func UpdateTutoringPurchase(tutoringPurchase *TutoringPurchase) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&tutoringPurchase.Purchase)
	return
}
