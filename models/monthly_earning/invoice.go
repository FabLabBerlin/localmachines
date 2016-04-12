package monthly_earning

import (
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/users"
)

type Invoice struct {
	User       users.User
	Purchases  purchases.Purchases
	VatPercent float64
}

func (inv *Invoice) byProductNameAndPricePerUnit() map[string]map[float64][]*purchases.Purchase {
	byProductNameAndPricePerUnit := make(map[string]map[float64][]*purchases.Purchase)
	for _, p := range inv.Purchases.Data {
		if _, ok := byProductNameAndPricePerUnit[p.ProductName()]; !ok {
			byProductNameAndPricePerUnit[p.ProductName()] = make(map[float64][]*purchases.Purchase)
		}
		if _, ok := byProductNameAndPricePerUnit[p.ProductName()][p.PricePerUnit]; !ok {
			byProductNameAndPricePerUnit[p.ProductName()][p.PricePerUnit] = make([]*purchases.Purchase, 0, 20)
		}
		byProductNameAndPricePerUnit[p.ProductName()][p.PricePerUnit] = append(byProductNameAndPricePerUnit[p.ProductName()][p.PricePerUnit], p)
	}
	return byProductNameAndPricePerUnit
}
