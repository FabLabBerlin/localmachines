package invoices

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego"
)

func CreateFastbillDrafts(inv *Invoice) (ids []int64, err error) {
	ids = make([]int64, 0, len(inv.UserSummaries))
	for _, userSummary := range inv.UserSummaries {
		if uid := userSummary.User.Id; uid == 19 {
			fbDraft, empty, err := createFastbillDraft(userSummary)
			if err != nil {
				return nil, fmt.Errorf("create draft for user %v: %v", uid, err)
			}
			if empty {
				beego.Debug("draft is empty")
				continue
			}
			id := fbDraft.Id
			beego.Info("Draft created with ID", id)
			ids = append(ids, id)
		}
	}
	return
}

func createFastbillDraft(userSummary *UserSummary) (fbDraft *fastbill.Invoice, empty bool, err error) {
	fbDraft = &fastbill.Invoice{
		CustomerId: userSummary.User.ClientId,
		TemplateId: fastbill.TemplateStandardId,
		Items:      make([]fastbill.Item, 0, 10),
	}
	memberships, err := models.GetUserMemberships(userSummary.User.Id)
	if err != nil {
		return nil, false, fmt.Errorf("GetUserMemberships: %v", err)
	}

	if len(userSummary.Purchases.Data) == 0 &&
		(memberships == nil || len(memberships.Data) == 0) {
		return nil, true, nil
	}

	for _, m := range memberships.Data {
		item := fastbill.Item{
			Description: m.Title,
			Quantity:    1,
			UnitPrice:   m.MonthlyPrice,
			VatPercent:  19,
		}
		fbDraft.Items = append(fbDraft.Items, item)
	}

	byProductNameAndPricePerUnit := userSummary.byProductNameAndPricePerUnit()

	for productName, byPricePerUnit := range byProductNameAndPricePerUnit {
		for pricePerUnit, ps := range byPricePerUnit {
			var quantity float64
			var discPrice float64
			var unitPrice float64
			discount := false
			for _, purchase := range ps {
				quantity += purchase.Quantity
				priceDisc, err := purchases.PriceTotalDisc(purchase)
				if err != nil {
					return nil, false, fmt.Errorf("PriceTotalDisc: %v", err)
				}
				discPrice += priceDisc
				affected, err := purchases.AffectedMemberships(purchase)
				if err != nil {
					return nil, false, fmt.Errorf("affected memberships: %v", err)
				}
				discount = len(affected) > 0
			}
			beego.Info("")
			beego.Info("discount=", discount)
			beego.Info("unitPrice=", unitPrice)
			if discount {
				unitPrice = pricePerUnit
			} else {
				unitPrice = discPrice / quantity
			}

			item := fastbill.Item{
				Description: productName,
				Quantity:    quantity,
				UnitPrice:   unitPrice,
				VatPercent:  19,
			}
			fbDraft.Items = append(fbDraft.Items, item)
		}
	}

	if _, err := fbDraft.Submit(); err != nil {
		return nil, false, fmt.Errorf("submit: %v", err)
	}

	return
}
