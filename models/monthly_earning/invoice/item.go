package invoice

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego/orm"
	"math"
)

func init() {
	orm.RegisterModel(new(Item))
}

type Item struct {
	Id int64
	// Types
	//  - membership
	//  - activation
	//  - tutoring
	//  - reservation
	Type                 string
	Description          string
	Quantity             float64
	VatPercent           float64
	UnitPriceCents       int64
	TotalPriceCents      int64
	DiscountedPriceCents int64
	IsGross              bool
}

func NewItem(obj interface{}, vatPercent float64) (item *Item, err error) {
	item = &Item{
		IsGross: true,
	}

	switch v := obj.(type) {
	case *models.Membership:
		item.Type = "membership"
		item.Description = v.Title + "Membership (unit: month)"
		item.Quantity = 1
		item.VatPercent = vatPercent
		item.UnitPriceCents = int64(Round(v.MonthlyPrice * 100))
		item.TotalPriceCents = int64(Round(v.MonthlyPrice))
		item.DiscountedPriceCents = int64(Round(v.MonthlyPrice * 100))
		break
	case *purchases.Purchase:
		unit := v.PriceUnit
		pricePerUnit := v.PricePerUnit
		item.Type = v.Type
		item.Description = v.ProductName() + " (unit: " + unit + ")"
		item.Quantity += v.Quantity
		priceDisc, err := purchases.PriceTotalDisc(v)
		if err != nil {
			return nil, fmt.Errorf("PriceTotalDisc: %v", err)
		}
		discPrice := 0.0
		discPrice += priceDisc
		affected, err := purchases.AffectedMemberships(v)
		if err != nil {
			return nil, fmt.Errorf("affected memberships: %v", err)
		}
		if discount := len(affected) > 0; discount {
			item.DiscountedPriceCents = int64(Round(discPrice / item.Quantity * 100))
		} else {
			item.DiscountedPriceCents = int64(Round(pricePerUnit * 100))
		}
		item.TotalPriceCents = int64(Round(pricePerUnit * 100))
		break
	}

	return
}

func (item *Item) TableName() string {
	return "invoice_items"
}

// https://gist.github.com/DavidVaini/10308388
func Round(f float64) float64 {
	return math.Floor(f + .5)
}
