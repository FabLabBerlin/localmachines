package invoices

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/astaxie/beego"
	"time"
)

const (
	IS_GROSS_BRUTTO = "1"
	IS_GROSS_NETTO  = "0"
)

type DraftsCreationReport struct {
	Ids                 []int64
	SuccessUids         []int64
	EmptyUids           []int64
	AlreadyExportedUids []int64
	Errors              []DraftsCreationError
}

type DraftsCreationError struct {
	UserId  int64
	Problem string
}

func CreateFastbillDrafts(inv *Invoice) (report DraftsCreationReport) {
	report.Ids = make([]int64, 0, len(inv.UserSummaries))
	report.SuccessUids = make([]int64, 0, len(inv.UserSummaries))
	report.EmptyUids = make([]int64, 0, len(inv.UserSummaries))
	report.AlreadyExportedUids = make([]int64, 0, len(inv.UserSummaries))
	report.Errors = make([]DraftsCreationError, 0, len(inv.UserSummaries))

	for _, userSummary := range inv.UserSummaries {
		uid := userSummary.User.Id
		if uid != 19 {
			continue
		}
		if r := userSummary.User.GetRole(); (r == user_roles.STAFF || r == user_roles.ADMIN || r == user_roles.SUPER_ADMIN) && uid != 19 {
			e := DraftsCreationError{
				UserId:  uid,
				Problem: "User role is " + r.String(),
			}
			report.Errors = append(report.Errors, e)
			beego.Error("no draft created for user", uid, ":", e.Problem)
		} else {
			fbDraft, empty, err := CreateFastbillDraft(inv, userSummary)
			if err == fastbill.ErrInvoiceAlreadyExported {
				beego.Info("draft for user", uid, "already exported")
				report.AlreadyExportedUids = append(report.AlreadyExportedUids, uid)
				report.SuccessUids = append(report.SuccessUids, uid)
				continue
			} else if err != nil {
				e := DraftsCreationError{
					UserId:  uid,
					Problem: err.Error(),
				}
				report.Errors = append(report.Errors, e)
				beego.Error("create draft for user", uid, ":", err)
				continue
			} else if !empty {
				report.SuccessUids = append(report.SuccessUids, uid)
			}
			if empty {
				report.EmptyUids = append(report.EmptyUids, uid)
				beego.Debug("draft is empty")
				continue
			}
			id := fbDraft.Id
			beego.Info("Draft created with ID", id)
			report.Ids = append(report.Ids, id)
		}
	}
	return
}

func CreateFastbillDraft(inv *Invoice, userSummary *UserSummary) (fbDraft *fastbill.Invoice, empty bool, err error) {
	fbDraft = &fastbill.Invoice{
		CustomerNumber: userSummary.User.ClientId,
		TemplateId:     fastbill.TemplateStandardId,
		Items:          make([]fastbill.Item, 0, 10),
	}
	if fbDraft.Month, fbDraft.Year, err = getFastbillMonthYear(inv); err != nil {
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
		if m.MonthlyPrice > 0 && m.StartDate.Before(inv.PeriodTo()) && m.EndDate.After(inv.PeriodFrom()) {
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

func getFastbillMonthYear(inv *Invoice) (month string, year int, err error) {
	if inv.MonthFrom != inv.MonthTo || inv.YearFrom != inv.YearTo {
		return "", 0, fmt.Errorf("2 months present")
	}
	return time.Month(inv.MonthFrom).String(), inv.YearFrom, nil
}
