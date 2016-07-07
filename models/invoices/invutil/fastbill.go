package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models/coupons"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/settings"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

const (
	IS_GROSS_BRUTTO = "1"
	IS_GROSS_NETTO  = "0"
)

func (inv *Invoice) fastbillBeforeCheck() (err error) {
	if inv.User == nil {
		return fmt.Errorf("nil user")
	}
	if inv.User.NoAutoInvoicing {
		return fmt.Errorf("no auto invoicing is true for user")
	}
	return
}

func (inv *Invoice) FastbillCancel() (err error) {
	if err := inv.fastbillBeforeCheck(); err != nil {
		return fmt.Errorf("fastbill before check: %v", err)
	}

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

	if err := FastbillSync(inv.LocationId, inv.User); err != nil {
		beego.Error("Error syncing fastbill invoices of user")
	}

	return
}

// FastbillComplete invoice. Data must be synchronized, so better to do it
// too often than to seldomly.
func (inv *Invoice) FastbillComplete() (err error) {
	if err := inv.fastbillBeforeCheck(); err != nil {
		return fmt.Errorf("fastbill before check: %v", err)
	}

	if inv.Year > time.Now().Year() ||
		(inv.Year == time.Now().Year() &&
			int(inv.Month) >= int(time.Now().Month())) {
		return fmt.Errorf("invoices must be from a past month")
	}

	_, empty, err := inv.FastbillCreateDraft(true)
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

func (inv *Invoice) FastbillCreateDraft(overwriteExisting bool) (fbDraft *fastbill.Invoice, empty bool, err error) {
	if err := inv.fastbillBeforeCheck(); err != nil {
		return nil, false, fmt.Errorf("fastbill before check: %v", err)
	}

	locSettings, err := settings.GetAllAt(inv.LocationId)
	if err != nil {
		beego.Error("FastbillCreateDraft: get settings:", err)
		return nil, false, fastbill.ErrObtainingValidTemplateId
	}

	templateId := locSettings.GetInt(inv.LocationId, settings.FASTBILL_TEMPLATE_ID)
	if templateId == nil || *templateId <= 0 {
		beego.Error("FastbillCreateDraft: invalid template id:", templateId)
		return nil, false, fastbill.ErrObtainingValidTemplateId
	}

	fbDraft = &fastbill.Invoice{
		CustomerNumber: inv.User.ClientId,
		TemplateId:     *templateId,
		Items:          make([]fastbill.Item, 0, 10),
	}
	intv := inv.Interval()
	if fbDraft.Month, fbDraft.Year, err = getFastbillMonthYear(&intv); err != nil {
		return nil, false, fmt.Errorf("get fastbill month: %v", err)
	}
	ms, err := user_memberships.GetForInvoice(inv.Id)
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

func (inv *Invoice) FastbillSend() (err error) {
	if err := inv.fastbillBeforeCheck(); err != nil {
		return fmt.Errorf("fastbill before check: %v", err)
	}

	if err := fastbill.SendInvoiceByEmail(inv.FastbillId, inv.User); err != nil {
		return fmt.Errorf("fastbill send invoice by email: %v", err)
	}

	inv.Sent = true

	if err := inv.Save(); err != nil {
		return fmt.Errorf("save: %v", err)
	}

	return
}

func (inv *Invoice) FastbillSendCanceled() (err error) {
	if err := inv.fastbillBeforeCheck(); err != nil {
		return fmt.Errorf("fastbill before check: %v", err)
	}

	if err := FastbillSync(inv.LocationId, inv.User); err != nil {
		beego.Error("Error syncing fastbill invoices of user")
	}

	if err := fastbill.SendInvoiceByEmail(inv.CanceledFastbillId, inv.User); err != nil {
		return fmt.Errorf("fastbill send canceled invoice by email: %v", err)
	}

	inv.CanceledSent = true

	if err := inv.Save(); err != nil {
		return fmt.Errorf("save: %v", err)
	}

	return
}

func FastbillSync(locId int64, u *users.User) (err error) {
	fbCustId, err := fastbill.GetCustomerId(*u)
	if err != nil {
		return fmt.Errorf("get customer id: %v", err)
	}

	l, err := fastbill.ListInvoices(fbCustId)
	if err != nil {
		return fmt.Errorf("Failed to get invoice list from fastbill: %v", err)
	}

	invs, err := invoices.GetAllOfUserAt(locId, u.Id)
	if err != nil {
		return fmt.Errorf("get invoices of user at location: %v", err)
	}

	// Sync draft and outgoing data
	for _, fbInv := range l {
		var inv *invoices.Invoice

		for _, iv := range invs {
			if iv.FastbillId == fbInv.Id {
				inv = iv
				break
			}
		}

		if inv == nil {
			continue
		}

		inv.Total = fbInv.Total
		inv.VatPercent = fbInv.VatPercent
		inv.Canceled = fbInv.Canceled()
		inv.DueDate = fbInv.DueDate()
		inv.InvoiceDate = fbInv.InvoiceDate()
		inv.PaidDate = fbInv.PaidDate()
		if fbInv.Type == "draft" || fbInv.Type == "outgoing" {
			inv.Status = fbInv.Type
		}

		if err = inv.Save(); err != nil {
			return fmt.Errorf("save invoice: %v", err)
		}
	}

	// Sync canceled/credit data
	for _, fbInv := range l {
		if fbInv.Type != "credit" {
			continue
		}

		var inv *invoices.Invoice

		for _, iv := range invs {
			if len(strings.TrimSpace(iv.FastbillNo)) < 3 {
				continue
			}
			if strings.Contains(fbInv.InvoiceNumber, iv.FastbillNo) {
				inv = iv
				break
			}
		}

		if inv == nil {
			continue
		}

		inv.CanceledFastbillId = fbInv.Id
		inv.CanceledFastbillNo = fbInv.InvoiceNumber

		if err = inv.Save(); err != nil {
			return fmt.Errorf("save invoice: %v", err)
		}
	}

	return
}
