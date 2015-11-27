package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Tutor struct {
	Product Product
}

type TutorList struct {
	Data []*Product
}

type TutoringPurchase struct {
	Id         int64
	TutorId    int64
	UserId     int64
	StartTime  time.Time
	EndTime    time.Time
	TotalTime  float64
	TotalPrice float64
	VAT        float64
}

type TutoringPurchaseList struct {
	Data []*TutoringPurchase
}

// Get a list of tutors
func GetTutorList() (*TutorList, error) {
	tutorsAsProducts, err := GetAllProductsOfType(PRODUCT_TYPE_TUTOR)
	if err != nil {
		msg := "Failed to get tutors as products"
		beego.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	tutorList := TutorList{Data: tutorsAsProducts}

	return &tutorList, nil
}

// Get a list of tutoring purchases
func GetTutoringPurchaseList() (*TutoringPurchaseList, error) {
	purchase1 := TutoringPurchase{Id: 1, TutorId: 1, UserId: 1, VAT: 19.0}
	purchase1.StartTime = time.Now()
	purchase1.EndTime = time.Now().Add(time.Duration(time.Hour * 2))
	purchase1.TotalTime = purchase1.EndTime.Sub(purchase1.StartTime).Hours()
	purchase1.TotalPrice = purchase1.TotalTime * 60.0

	purchase2 := TutoringPurchase{Id: 2, TutorId: 1, UserId: 1, VAT: 19.0}
	purchase2.StartTime = time.Now()
	purchase2.EndTime = time.Now().Add(time.Duration(time.Hour * 2))
	purchase2.TotalTime = purchase2.EndTime.Sub(purchase2.StartTime).Hours()
	purchase2.TotalPrice = purchase2.TotalTime * 60.0

	purchaseList := TutoringPurchaseList{}
	purchaseList.Data = append(purchaseList.Data, &purchase1, &purchase2)

	return &purchaseList, nil
}

func CreateTutor(tutor *Tutor) (*Tutor, error) {
	o := orm.NewOrm()

	// In case the type has not been added in upper layers
	tutor.Product.Type = PRODUCT_TYPE_TUTOR

	// Get user name by user ID
	user := User{}
	user.Id = tutor.Product.UserId
	err := o.Read(&user)
	if err != nil {
		beego.Error("Failed to read user:", err)
		return nil, fmt.Errorf("Failed to read user: %v", err)
	}

	tutor.Product.Name = fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	var productId int64
	productId, err = o.Insert(&tutor.Product)
	if err != nil {
		msg := "Failed to insert tutor"
		beego.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	tutor.Product.Id = productId

	return tutor, nil
}

func UpdateTutor(tutor *Tutor) error {
	return nil
}
