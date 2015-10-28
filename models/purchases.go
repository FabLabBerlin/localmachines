package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"time"
)

// This is a purchase row that appears in the XLSX file
type Purchase struct {
	Activation      *Activation
	Reservation     *Reservation
	Machine         *Machine
	MachineUsage    time.Duration
	User            User
	Memberships     []*Membership
	TotalPrice      float64
	DiscountedTotal float64
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

func (this *Purchase) PricePerUnit() float64 {
	if this.Reservation != nil {
		if this.Reservation.CurrentPrice != 0 {
			return this.Reservation.CurrentPrice / 2
		} else {
			return 0
		}
	} else {
		return this.Activation.CurrentMachinePrice
	}
}

func (this *Purchase) PriceUnit() string {
	if this.Activation != nil {
		return this.Activation.CurrentMachinePriceUnit
	} else {
		return this.Reservation.PriceUnit()
	}
}

func (this *Purchase) Usage() float64 {
	if this.Activation != nil {
		return this.MachineUsage.Minutes()
	} else {
		return float64(this.Reservation.Slots())
	}
}

func PriceTotalExclDisc(p *Purchase) float64 {
	if p.Activation != nil {
		var pricePerSecond float64
		beego.Trace("p.Activation.CurrentMachinePriceUnit", p.Activation.CurrentMachinePriceUnit)
		beego.Trace("p.Activation.CurrentMachinePrice", p.Activation.CurrentMachinePrice)
		switch p.Activation.CurrentMachinePriceUnit {
		case "minute":
			beego.Trace("minute")
			pricePerSecond = float64(p.Activation.CurrentMachinePrice) / 60
			break
		case "hour":
			beego.Trace("hour")
			pricePerSecond = float64(p.Activation.CurrentMachinePrice) / 60 / 60
			break
		}
		beego.Trace(p.MachineUsage.Seconds())
		beego.Trace(pricePerSecond)
		ret := p.MachineUsage.Seconds() * pricePerSecond
		return ret
	} else {
		return float64(p.Reservation.Slots()) * p.PricePerUnit()
	}
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
	var timeStartI time.Time
	var timeStartJ time.Time
	if (*this.Data[i]).Activation != nil {
		timeStartI = (*this.Data[i]).Activation.TimeStart
	} else {
		timeStartI = (*this.Data[i]).Reservation.TimeStart
	}
	if (*this.Data[j]).Activation != nil {
		timeStartJ = (*this.Data[j]).Activation.TimeStart
	} else {
		timeStartI = (*this.Data[j]).Reservation.TimeStart
	}
	if timeStartI.Before(timeStartJ) {
		return true
	} else if timeStartJ.Before(timeStartI) {
		return false
	} else {
		return (*this.Data[i]).Machine.Name < (*this.Data[j]).Machine.Name
	}
}

func (this Purchase) ProductName() string {
	if this.Activation != nil {
		return this.Machine.Name
	} else {
		return "Reservation (" + this.Machine.Name + ")"
	}
}

func (this Purchases) Swap(i, j int) {
	*this.Data[i], *this.Data[j] = *this.Data[j], *this.Data[i]
}
