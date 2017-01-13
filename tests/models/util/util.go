package util

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"time"
)

var TIME_START = time.Now().AddDate(0, -1, -1)

func CreateTestPurchase(machineId, categoryId int64, machineName string,
	minutes time.Duration, pricePerMinute float64) *purchases.Purchase {

	m := machine.Machine{}
	m.Id = machineId
	m.Name = machineName
	m.PriceUnit = "minute"
	m.Price = pricePerMinute
	m.TypeId = categoryId

	invAct := &purchases.Purchase{
		LocationId:   1,
		Type:         purchases.TYPE_ACTIVATION,
		TimeStart:    TIME_START,
		PricePerUnit: pricePerMinute,
		PriceUnit:    "minute",
		Quantity:     minutes.Minutes(),
		Machine:      &m,
		MachineId:    machineId,
		Memberships: []*memberships.Membership{
			{
				Id:                    42,
				Title:                 "Half price",
				ShortName:             "HP",
				MachinePriceDeduction: 50,
				AffectedCategories:    fmt.Sprintf("[%v]", categoryId),
			},
		},
	}
	invAct.TotalPrice = purchases.PriceTotalExclDisc(invAct)
	var err error
	invAct.DiscountedTotal, err = purchases.PriceTotalDisc(invAct)
	if err != nil {
		panic(err.Error())
	}
	return invAct
}
