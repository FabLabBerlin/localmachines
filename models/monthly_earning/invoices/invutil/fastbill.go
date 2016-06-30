package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models/coupons"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	IS_GROSS_BRUTTO = "1"
	IS_GROSS_NETTO  = "0"
)

func (inv *Invoice) Cancel() (err error) {
	if inv.Canceled {
		return fmt.Errorf("invoice already marked as canceled in EASY LAB")
	}

	if err = fastbill.CancelInvoice(inv.FastbillId); err != nil {
		return fmt.Errorf("fastbill cancel invoice: %v", err)
	}

	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		return fmt.Errorf("begin tx: %v", err)
	}

	inv.Canceled = true
	if err := inv.SaveOrm(o); err != nil {
		return fmt.Errorf("error saving invoice changes: %v", err)
	}

	draft := &Invoice{}
	draft.LocationId = inv.LocationId
	draft.Month = inv.Month
	draft.Year = inv.Year
	draft.CustomerId = inv.CustomerId
	draft.CustomerNo = inv.CustomerNo
	draft.UserId = inv.UserId
	draft.Total = inv.Total
	draft.Status = "draft"
	draft.VatPercent = inv.VatPercent
	if _, err = o.Insert(&draft.Invoice); err != nil {
		return fmt.Errorf("insert draft: %v", err)
	}

	for _, p := range inv.Purchases {
		if err := p.CloneOrm(o, &draft.Invoice); err != nil {
			return fmt.Errorf("clone purchase: %v", err)
		}
	}

	for _, um := range inv.UserMemberships.Data {
		err := um.UserMembership().CloneOrm(o, draft.Invoice.Id, draft.Invoice.Status)
		if err != nil {
			return fmt.Errorf("clone user membership: %v", err)
		}
	}

	if err := o.Commit(); err != nil {
		return fmt.Errorf("commit tx: %v", err)
	}

	return
}

// CompleteFastbill invoice. Data must be synchronized, so better to do it
// too often than to seldomly.
func (inv *Invoice) CompleteFastbill() (err error) {
	if inv.Year > time.Now().Year() ||
		(inv.Year == time.Now().Year() &&
			int(inv.Month) >= int(time.Now().Month())) {
		return fmt.Errorf("invoices must be from a past month")
	}

	_, empty, err := inv.CreateFastbillDraft(true)
	if err != nil {
		return fmt.Errorf("create fastbill draft: %v", err)
	}
	if !empty {
		fbNumber, err := fastbill.CompleteInvoice(inv.FastbillId)
		if err != nil {
			return fmt.Errorf("fastbill complete invoice: %v", err)
		}
		inv.FastbillNo = fbNumber
	}
	inv.Status = "outgoing"
	if err := inv.Save(); err != nil {
		return fmt.Errorf("error saving invoice changes: %v", err)
	}
	return
}

func (inv *Invoice) CreateFastbillDraft(overwriteExisting bool) (fbDraft *fastbill.Invoice, empty bool, err error) {
	fbDraft = &fastbill.Invoice{
		CustomerNumber: inv.User.ClientId,
		TemplateId:     fastbill.TemplateMakeaIndustriesId,
		Items:          make([]fastbill.Item, 0, 10),
	}
	intv := inv.Interval()
	if fbDraft.Month, fbDraft.Year, err = getFastbillMonthYear(&intv); err != nil {
		return nil, false, fmt.Errorf("get fastbill month: %v", err)
	}
	ms, err := memberships.GetUserMembershipsForInvoice(inv.Id)
	if err != nil {
		return nil, false, fmt.Errorf("GetUserMemberships: %v", err)
	}

	fbDraft.CustomerId, err = fastbill.GetCustomerId(*inv.User)
	if err != nil {
		return nil, false, fmt.Errorf("Fastbill customer id: %v", err)
	}

	if len(inv.Purchases) == 0 &&
		(ms == nil || len(ms.Data) == 0) {
		return nil, true, nil
	}

	invoiceValue := 0.0

	// Add Memberships
	for _, m := range ms.Data {
		if m.MonthlyPrice > 0 && m.StartDate.Before(inv.Interval().TimeTo()) && m.EndDate.After(inv.Interval().TimeFrom()) {
			item := fastbill.Item{
				Description: m.Title + " Membership (unit: month)",
				Quantity:    1,
				UnitPrice:   m.MonthlyPrice,
				IsGross:     IS_GROSS_BRUTTO,
				VatPercent:  inv.VatPercent,
			}
			invoiceValue += item.UnitPrice
			fbDraft.Items = append(fbDraft.Items, item)
		}
	}

	// Add Product Purchases

	for _, p := range inv.Purchases {
		discount := false
		priceDisc, err := purchases.PriceTotalDisc(p)
		if err != nil {
			return nil, false, fmt.Errorf("PriceTotalDisc: %v", err)
		}
		affected, err := purchases.AffectedMemberships(p)
		if err != nil {
			return nil, false, fmt.Errorf("affected memberships: %v", err)
		}
		discount = len(affected) > 0
		var unitPrice float64
		if discount {
			unitPrice = priceDisc / p.Quantity
		} else {
			unitPrice = p.PricePerUnit
		}

		item := fastbill.Item{
			Description: p.ProductName() + " (unit: " + p.PriceUnit + ")",
			Quantity:    p.Quantity,
			UnitPrice:   unitPrice,
			IsGross:     IS_GROSS_BRUTTO,
			VatPercent:  inv.VatPercent,
		}

		if v := item.UnitPrice * item.Quantity; v > 0.01 {
			invoiceValue += v
			fbDraft.Items = append(fbDraft.Items, item)
		}
	}

	// Add Coupons
	cs, err := coupons.GetAllCouponsOf(inv.LocationId, inv.User.Id)
	if err != nil {
		return nil, false, fmt.Errorf("get all coupons: %v", err)
	}
	rebateValue := 0.0
	for _, c := range cs {
		usage, err := c.UseForInvoice(invoiceValue-rebateValue, time.Month(inv.Interval().MonthFrom), inv.Interval().YearFrom)
		if err != nil {
			return nil, false, fmt.Errorf("use for invoice: %v", err)
		}
		if usage != nil {
			rebateValue += usage.Value
		}
	}
	fbDraft.CashDiscountPercent = fmt.Sprintf("%v", rebateValue/invoiceValue*100)

	if _, err := fbDraft.Submit(overwriteExisting); err == fastbill.ErrInvoiceAlreadyExported {
		return nil, false, fastbill.ErrInvoiceAlreadyExported
	} else if err != nil {
		return nil, false, fmt.Errorf("submit: %v", err)
	} else {
		inv.FastbillId = fbDraft.Id
		fmt.Printf("inv.FastbillId <- %v", fbDraft.Id)
		if err := inv.Save(); err != nil {
			return nil, false, fmt.Errorf("invoice save fastbill db id: %v", err)
		}
	}

	return
}

func getFastbillMonthYear(i *lib.Interval) (month string, year int, err error) {
	if i.MonthFrom != i.MonthTo || i.YearFrom != i.YearTo {
		return "", 0, fmt.Errorf("2 months present")
	}
	return time.Month(i.MonthFrom).String(), i.YearFrom, nil
}
