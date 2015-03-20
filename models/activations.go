package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"errors"
)

func init() {
	orm.RegisterModel(new(Activation))
}

type Activation struct {
	Id               int `orm:"auto";"pk"`
	InvoiceId        int `orm:"null"`
	UserId           int
	MachineId        int
	Active           bool
	TimeStart        time.Time
	TimeEnd          time.Time `orm:"null"`
	TimeTotal        int
	UsedKwh          float32
	DiscountPercents float32
	DiscountFixed    float32
	VatRate          float32
	CommentRef       string `orm:"size(255)"`
	Invoiced         bool
	Changed          bool
}

func GetActiveActivations() ([]*Activation, error){
	var activations []*Activation
	o := orm.NewOrm()
	num, err := o.QueryTable("activation").
		Filter("active", true).All(&activations)
	if err != nil {
		return nil, errors.New("Failed to get active activations")
	}
	beego.Trace("Got num activations:", num)
	return activations, nil
}
