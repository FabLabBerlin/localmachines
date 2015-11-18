package models

import (
	//"github.com/astaxie/beego/orm"
	"time"
)

type Tutor struct {
	Id        int64
	Name      string
	Price     float64
	PriceUnit string
}

type TutorList struct {
	Data []*Tutor
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

func init() {
	//orm.RegisterModel(new(Tutor))
	//orm.RegisterModel(new(TutoringPurchase))
}

// Get a list of tutors
func GetTutorList() (*TutorList, error) {
	tutor1 := Tutor{Id: 1, Name: "Ahmad Taled", Price: 60.0, PriceUnit: "hour"}
	tutor2 := Tutor{Id: 2, Name: "Tina Atari", Price: 60.0, PriceUnit: "hour"}
	tutorList := TutorList{}
	tutorList.Data = append(tutorList.Data, &tutor1, &tutor2)
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
