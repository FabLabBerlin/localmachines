package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"sort"
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
		if this.Machine.ReservationPriceHourly != nil {
			return *this.Machine.ReservationPriceHourly / 2
		} else {
			return 0
		}
	} else {
		return float64(this.Machine.Price)
	}
}

func PriceTotalExclDisc(p *Purchase) float64 {
	if p.Activation != nil {
		var pricePerSecond float64
		switch p.Machine.PriceUnit {
		case "minute":
			pricePerSecond = float64(p.Machine.Price) / 60
			break
		case "hour":
			pricePerSecond = float64(p.Machine.Price) / 60 / 60
			break
		}
		return p.MachineUsage.Seconds() * pricePerSecond
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

func (this Purchases) Swap(i, j int) {
	*this.Data[i], *this.Data[j] = *this.Data[j], *this.Data[i]
}

func (this Purchases) SummarizedByMachine() (
	Purchases, error) {

	byMachine := make(map[string]*Purchase)
	for _, activation := range this.Data {
		summary, ok := byMachine[activation.Machine.Name]
		if !ok {
			summary = &Purchase{
				Activation:      &Activation{},
				MachineUsage:    0,
				TotalPrice:      0,
				DiscountedTotal: 0,
				Machine:         activation.Machine,
				Memberships:     activation.Memberships,
			}
			byMachine[activation.Machine.Name] = summary
		}
		summary.MachineUsage += activation.MachineUsage
		summary.TotalPrice += activation.TotalPrice
		summary.DiscountedTotal += activation.DiscountedTotal

	}

	sumPurchasesData := make([]*Purchase, 0, len(byMachine))
	sumPurchases := Purchases{}
	sumPurchases.Data = sumPurchasesData
	for _, summary := range byMachine {
		sumPurchases.Data = append(sumPurchases.Data, summary)
	}
	sort.Stable(sumPurchases)

	return sumPurchases, nil
}
