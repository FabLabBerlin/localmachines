package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	PURCHASE_TYPE_ACTIVATION  = "activation"
	PURCHASE_TYPE_RESERVATION = "reservation"
)

// This is a purchase row that appears in the XLSX file
type Purchase struct {
	Id        int64 `orm:"auto";"pk"`
	Type      string
	ProductId int64
	Created   time.Time `orm:"type(datetime)"`

	User   User `orm:"-" json:"-"`
	UserId int64

	TimeStart    time.Time `orm:"type(datetime)"`
	TimeEnd      time.Time `orm:"type(datetime)"`
	Quantity     float64
	PricePerUnit float64
	PriceUnit    string
	Vat          float64

	TotalPrice      float64 `orm:"-"`
	DiscountedTotal float64 `orm:"-"`

	// Activation fields:
	ActivationRunning bool

	// Reservation fields:
	ReservationDisabled bool

	// Activation+Reservation fields:
	Machine   *Machine `orm:"-"`
	MachineId int64

	// Old fields:
	Activation   *Activation   `orm:"-"`
	Reservation  *Reservation  `orm:"-"`
	MachineUsage time.Duration `orm:"-"`
	Memberships  []*Membership `orm:"-"`
}

func (this *Purchase) TableName() string {
	return "purchases"
}

func init() {
	orm.RegisterModel(new(Purchase))
}

func DeletePurchase(id int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Delete(&Purchase{Id: id})
	return
}

func (this *Purchase) MembershipStr() string {
	membershipStr := ""
	for _, membership := range this.Memberships {
		memStr := fmt.Sprintf("%s (%d%%)",
			membership.ShortName,
			membership.MachinePriceDeduction)
		if membershipStr == "" {
			membershipStr = memStr
		} else {
			membershipStr = fmt.Sprintf("%s, %s",
				membershipStr, memStr)
		}
	}
	if membershipStr == "" {
		membershipStr = "None"
	}
	return membershipStr
}

func (this *Purchase) Usage() float64 {
	if this.Activation != nil {
		return this.MachineUsage.Minutes()
	} else {
		return float64(this.Reservation.Slots())
	}
}

func PriceTotalExclDisc(p *Purchase) float64 {
	var pricePerSecond float64
	switch p.PriceUnit {
	case "minute":
		pricePerSecond = float64(p.PricePerUnit) / 60
		break
	case "hour":
		pricePerSecond = float64(p.PricePerUnit) / 60 / 60
		break
	}
	ret := p.Quantity * pricePerSecond
	return ret
}

func PriceTotalDisc(p *Purchase) (float64, error) {
	if p.Reservation != nil {
		return PriceTotalExclDisc(p), nil
	}

	priceTotal := PriceTotalExclDisc(p)
	for _, membership := range p.Memberships {

		// We need to know whether the machine is affected by the base membership
		// as well as the individual activation is affected by the user membership
		isAffected, err := membership.IsMachineAffected(p.Machine.Id)
		if err != nil {
			beego.Error(
				"Failed to check whether machine is affected by membership:", err)
			return 0, fmt.Errorf(
				"Failed to check whether machine is affected by membership")
		}

		machinePriceDeduction := 0.0
		if isAffected {
			machinePriceDeduction = float64(membership.MachinePriceDeduction)
		}
		// Discount total price
		priceTotal = priceTotal - (priceTotal * machinePriceDeduction / 100.0)
	}
	return priceTotal, nil
}

type Purchases struct {
	Data []*Purchase
}

func (this Purchases) Len() int {
	return len(this.Data)
}

func (this Purchases) Less(i, j int) bool {
	var timeStartI = (*this.Data[i]).TimeStart
	var timeStartJ = (*this.Data[j]).TimeStart
	if timeStartI.Before(timeStartJ) {
		return true
	} else if timeStartJ.Before(timeStartI) {
		return false
	} else {
		return (*this.Data[i]).Machine.Name < (*this.Data[j]).Machine.Name
	}
}

func (this Purchase) ProductName() string {
	switch this.Type {
	case PURCHASE_TYPE_ACTIVATION:
		return this.Machine.Name
	case PURCHASE_TYPE_RESERVATION:
		return "Reservation (" + this.Machine.Name + ")"
	}
	return "Unnamed product"
}

func (this Purchases) Swap(i, j int) {
	*this.Data[i], *this.Data[j] = *this.Data[j], *this.Data[i]
}
