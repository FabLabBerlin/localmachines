package invoices

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego"
)

const (
	IS_GROSS_BRUTTO = "1"
	IS_GROSS_NETTO  = "0"
)

func CreateFastbillDrafts(inv *Invoice) (ids []int64, err error) {
	ids = make([]int64, 0, len(inv.UserSummaries))
	for _, userSummary := range inv.UserSummaries {
		if uid := userSummary.User.Id; uid == 19 {
			fbDraft, empty, err := CreateFastbillDraft(userSummary)
			if err == fastbill.ErrInvoiceAlreadyExported {
				beego.Info("draft for user", uid, "already exported")
				continue
			} else if err != nil {
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

func CreateFastbillDraft(userSummary *UserSummary) (fbDraft *fastbill.Invoice, empty bool, err error) {
	fbDraft = &fastbill.Invoice{
		CustomerNumber: userSummary.User.ClientId,
		TemplateId:     fastbill.TemplateStandardId,
		Items:          make([]fastbill.Item, 0, 10),
	}
	if fbDraft.Month, fbDraft.Year, err = getFastbillMonthYear(userSummary); err != nil {
		return nil, false, fmt.Errorf("get fastbill month: %v", err)
	}
	memberships, err := models.GetUserMemberships(userSummary.User.Id)
	if err != nil {
		return nil, false, fmt.Errorf("GetUserMemberships: %v", err)
	}

	fbDraft.CustomerId, err = fastbill.GetCustomerId(userSummary.User)
	if err != nil {
		return nil, false, fmt.Errorf("error getting fastbill customer id: %v", err)
	}

	if len(userSummary.Purchases.Data) == 0 &&
		(memberships == nil || len(memberships.Data) == 0) {
		return nil, true, nil
	}

	for _, m := range memberships.Data {
		if m.MonthlyPrice > 0 {
			item := fastbill.Item{
				Description: m.Title + " Membership",
				Quantity:    1,
				UnitPrice:   m.MonthlyPrice,
				IsGross:     IS_GROSS_BRUTTO,
				VatPercent:  19,
			}
			fbDraft.Items = append(fbDraft.Items, item)
		}
	}

	byProductNameAndPricePerUnit := userSummary.byProductNameAndPricePerUnit()

	for productName, byPricePerUnit := range byProductNameAndPricePerUnit {
		for pricePerUnit, ps := range byPricePerUnit {
			var quantity float64
			var discPrice float64
			var unitPrice float64
			var unit string
			discount := false
			for _, purchase := range ps {
				unit = purchase.PriceUnit
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
				unitPrice = discPrice / quantity
			} else {
				unitPrice = pricePerUnit
			}

			item := fastbill.Item{
				Description: productName + " (unit: " + unit + ")",
				Quantity:    quantity,
				UnitPrice:   unitPrice,
				IsGross:     IS_GROSS_BRUTTO,
				VatPercent:  19,
			}

			if item.UnitPrice > 0 {
				fbDraft.Items = append(fbDraft.Items, item)
			}
		}
	}

	if _, err := fbDraft.Submit(); err == fastbill.ErrInvoiceAlreadyExported {
		return nil, false, fastbill.ErrInvoiceAlreadyExported
	} else if err != nil {
		return nil, false, fmt.Errorf("submit: %v", err)
	}

	return
}

func getFastbillMonthYear(userSummary *UserSummary) (month string, year int, err error) {
	for _, p := range userSummary.Purchases.Data {
		m := p.TimeStart.Month().String()
		if month == "" {
			month = m
			year = p.TimeStart.Year()
		} else if month != m {
			return "", 0, fmt.Errorf("2 months present: %v and %v", m, month)
		}
	}
	return
}
