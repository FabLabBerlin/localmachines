package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/settings"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"time"
)

type Invoice struct {
	invoices.Invoice
	User            *users.User                     `json:",omitempty"`
	UserMemberships *memberships.UserMembershipList `json:",omitempty"`
	Purchases       purchases.Purchases             `json:",omitempty"`
	Sums            *Sums                           `json:",omitempty"`
}

// Send an invoice transactionally. This includes:
// 1. Send invoice through Fastbill
// 2. Synchronize Fastbill
// 3. Propagate Fastbill sync changes to associated purchases
func (inv *Invoice) Send() (err error) {
	return fmt.Errorf("not implemented")
}

type Sums struct {
	Memberships struct {
		PriceInclVAT float64
		PriceExclVAT float64
		PriceVAT     float64
		Undiscounted float64
	}
	Purchases struct {
		PriceInclVAT float64
		PriceExclVAT float64
		PriceVAT     float64
		Undiscounted float64
	}
	All struct {
		PriceInclVAT float64
		PriceExclVAT float64
		PriceVAT     float64
		Undiscounted float64
	}
}

func (inv *Invoice) ByProductNameAndPricePerUnit() map[string]map[float64][]*purchases.Purchase {
	byProductNameAndPricePerUnit := make(map[string]map[float64][]*purchases.Purchase)
	for _, p := range inv.Purchases {
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

func (inv *Invoice) CalculateTotals() (err error) {
	inv.Sums = &Sums{}

	for _, purchase := range inv.Purchases {
		inv.Sums.Purchases.Undiscounted += purchase.TotalPrice
		inv.Sums.Purchases.PriceInclVAT += purchase.DiscountedTotal
	}
	p := (100.0 + inv.VatPercent) / 100.0
	inv.Sums.Purchases.PriceExclVAT = inv.Sums.Purchases.PriceInclVAT / p
	inv.Sums.Purchases.PriceVAT = inv.Sums.Purchases.PriceInclVAT - inv.Sums.Purchases.PriceExclVAT

	ms, err := memberships.GetUserMembershipsForInvoice(inv.Id)
	if err != nil {
		return fmt.Errorf("GetUserMemberships: %v", err)
	}

	for _, m := range ms.Data {
		if m.StartDate.Before(inv.Interval().TimeFrom()) &&
			m.EndDate.After(inv.Interval().TimeTo()) {
			inv.Sums.Memberships.Undiscounted += m.MonthlyPrice
			inv.Sums.Memberships.PriceInclVAT += m.MonthlyPrice
		}
		inv.Sums.Memberships.PriceExclVAT = inv.Sums.Memberships.PriceInclVAT / p
		inv.Sums.Memberships.PriceVAT = inv.Sums.Memberships.PriceInclVAT - inv.Sums.Memberships.PriceExclVAT
	}
	inv.Sums.All.Undiscounted = inv.Sums.Purchases.Undiscounted + inv.Sums.Memberships.Undiscounted
	inv.Sums.All.PriceInclVAT = inv.Sums.Purchases.PriceInclVAT + inv.Sums.Memberships.PriceInclVAT
	inv.Sums.All.PriceExclVAT = inv.Sums.Purchases.PriceExclVAT + inv.Sums.Memberships.PriceExclVAT
	inv.Sums.All.PriceVAT = inv.Sums.Purchases.PriceVAT + inv.Sums.Memberships.PriceVAT

	return
}

func (inv *Invoice) Interval() lib.Interval {
	if inv.Month == 0 || inv.Year == 0 {
		panic(fmt.Sprintf("inv.Month=%v, inv.Year=%v", inv.Month, inv.Year))
	}

	return lib.Interval{
		MonthFrom: inv.Month,
		YearFrom:  inv.Year,
		MonthTo:   inv.Month,
		YearTo:    inv.Year,
	}
}

func (inv *Invoice) Load() (err error) {
	if inv.User, err = users.GetUser(inv.UserId); err != nil {
		return fmt.Errorf("get user(id=%v): %v", inv.UserId, err)
	}
	if inv.Purchases, err = purchases.GetByInvoiceId(inv.Id); err != nil {
		return fmt.Errorf("get purchases by invoice id: %v", err)
	}
	inv.UserMemberships, err = memberships.GetUserMembershipsForInvoice(inv.Id)
	if err != nil {
		return fmt.Errorf("get user memberships for invoice: %v", err)
	}
	return
}

func (inv *Invoice) SplitByMonths() (invs []*Invoice, err error) {
	var tMin time.Time
	invs = make([]*Invoice, 0, 10)

	if len(inv.Purchases) == 0 {
		return
	}

	for _, p := range inv.Purchases {
		if tMin.IsZero() || p.TimeStart.Before(tMin) {
			tMin = p.TimeStart
		}
	}

	i := 0
	for t := tMin; ; t = t.AddDate(0, 1, 0) {
		i++
		iv := &Invoice{}
		iv.Month = int(t.Month())
		iv.Year = t.Year()
		iv.Purchases = make([]*purchases.Purchase, 0, 20)
		iv.User = inv.User
		iv.VatPercent = inv.VatPercent
		invs = append(invs, iv)

		if i > 100 {
			return nil, fmt.Errorf("loop not finishing")
		}
		if t.Month() == time.Now().Month() && t.Year() == time.Now().Year() {
			break
		}
	}

	for _, iv := range invs {
		for _, p := range inv.Purchases {
			if iv.Interval().Contains(p.TimeStart) {
				iv.Purchases = append(iv.Purchases, p)
			}
		}
		if err := iv.CalculateTotals(); err != nil {
			return nil, fmt.Errorf("CalculateTotals: %v", err)
		}
	}

	return
}

func Get(id int64) (inv *Invoice, err error) {
	iv, err := invoices.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get invoice entity: %v", err)
	}
	tmp, err := toUtilInvoices(iv.LocationId, []*invoices.Invoice{
		iv,
	})

	return tmp[0], err
}

func GetAllOfUserAt(locId, userId int64) (invs []*Invoice, err error) {
	ivs, err := invoices.GetAllOfUserAt(locId, userId)
	if err != nil {
		return nil, fmt.Errorf("invoices.GetAllUserAt: %v", err)
	}

	if invs, err = toUtilInvoices(locId, ivs); err != nil {
		return nil, fmt.Errorf("to util invoices: %v", err)
	}

	return
}

func GetAllOfMonthAt(locId int64, year int, m time.Month) ([]*Invoice, error) {
	var invs []*Invoice

	ivs, err := invoices.GetAllInvoicesBetween(locId, year, int(m))
	if err != nil {
		return nil, fmt.Errorf("get all invoices between: %v", err)
	}

	if invs, err = toUtilInvoices(locId, ivs); err != nil {
		return nil, fmt.Errorf("to util invoices: %v", err)
	}

	return invs, err
}

func toUtilInvoices(locId int64, ivs []*invoices.Invoice) (invs []*Invoice, err error) {
	invs = make([]*Invoice, 0, len(ivs))

	for _, iv := range ivs {
		inv := &Invoice{
			Invoice: *iv,
		}
		if err = inv.Load(); err != nil {
			return nil, fmt.Errorf("load: %v", err)
		}
		invs = append(invs, inv)
	}

	ms, err := machine.GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("get all machines at %v: %v", locId, err)
	}
	msById := make(map[int64]*machine.Machine)
	for _, m := range ms {
		msById[m.Id] = m
	}

	umbsByUid := make(map[int64][]*memberships.UserMembership)
	if umbs, err := memberships.GetAllUserMembershipsAt(locId); err == nil {
		for _, umb := range umbs {
			uid := umb.UserId
			if _, ok := umbsByUid[uid]; !ok {
				umbsByUid[uid] = []*memberships.UserMembership{
					umb,
				}
			} else {
				umbsByUid[uid] = append(umbsByUid[uid], umb)
			}
		}
	} else {
		return nil, fmt.Errorf("Failed to get user memberships: %v", err)
	}

	mbsById := make(map[int64]*memberships.Membership)
	if mbs, err := memberships.GetAllMembershipsAt(locId); err == nil {
		for _, mb := range mbs {
			mbsById[mb.Id] = mb
		}
	} else {
		return nil, fmt.Errorf("Failed to get memberships: %v", err)
	}

	for _, inv := range invs {
		for _, p := range inv.Purchases {
			if m, ok := msById[p.MachineId]; ok {
				p.Machine = m
			}

			umbs, ok := umbsByUid[p.UserId]
			if !ok {
				umbs = []*memberships.UserMembership{}
			}
			for _, umb := range umbs {
				mbId := umb.MembershipId
				mb, ok := mbsById[mbId]
				if !ok {
					return nil, fmt.Errorf("Unknown membership id: %v", mbId)
				}
				if umb.EndDate.IsZero() {
					return nil, fmt.Errorf("end date is zero")
				}
				if umb.Interval().Contains(p.TimeStart) &&
					umb.InvoiceId == inv.Id {
					p.Memberships = append(p.Memberships, mb)
				}
			}
		}

		locSettings, err := settings.GetAllAt(inv.LocationId)
		if err != nil {
			return nil, fmt.Errorf("get settings: %v", err)
		}
		var vatPercent float64
		if vat := locSettings.GetFloat(inv.LocationId, settings.VAT); vat != nil {
			vatPercent = *vat
		} else {
			vatPercent = 19.0
		}
		for _, p := range inv.Purchases {
			p.TotalPrice = purchases.PriceTotalExclDisc(p)
			p.DiscountedTotal, err = purchases.PriceTotalDisc(p)
			if err != nil {
				return nil, fmt.Errorf("price total disc (purchase %v): %v", p.Id, err)
			}
			percent := (100.0 + vatPercent) / 100.0
			p.PriceExclVAT = p.DiscountedTotal / percent
			p.PriceVAT = p.DiscountedTotal - p.PriceExclVAT
		}
		if err = inv.CalculateTotals(); err != nil {
			return nil, fmt.Errorf("calculate totals: %v", err)
		}
	}

	return
}

func SyncFastbillInvoices(locId int64, year int, month time.Month) (err error) {
	l, err := fastbill.ListInvoices(year, time.Month(month))
	if err != nil {
		return fmt.Errorf("Failed to get invoice list from fastbill: %v", err)
	}

	usrs, err := users.GetAllUsersAt(locId)
	if err != nil {
		return fmt.Errorf("Failed to get user list: %v", err)
	}

	for _, fbInv := range l {
		inv := Invoice{}
		inv.LocationId = locId
		inv.FastbillId = fbInv.Id
		inv.FastbillNo = fbInv.InvoiceNumber
		inv.CustomerId = fbInv.CustomerId
		inv.Status = fbInv.Type
		inv.Total = fbInv.Total
		inv.VatPercent = fbInv.VatPercent
		inv.Canceled = fbInv.Canceled()
		inv.DueDate = fbInv.DueDate()
		inv.InvoiceDate = fbInv.InvoiceDate()
		inv.PaidDate = fbInv.PaidDate()
		inv.Month, inv.Year, inv.CustomerNo, err = fbInv.ParseTitle()
		if err != nil {
			beego.Error("Cannot parse", fbInv.InvoiceTitle)
			err = nil
			continue
		}
		for _, u := range usrs {
			if u.ClientId == inv.CustomerNo {
				inv.UserId = u.Id
				break
			}
		}
		if inv.UserId == 0 {
			beego.Error("Cannot find user for customer number", inv.CustomerNo)
			continue
		}
		beego.Info("Adding invoice of user", inv.UserId)
		if _, err := invoices.CreateOrUpdate(&inv.Invoice); err != nil {
			return fmt.Errorf("Failed to create or update inv: %v", err)
		}
	}
	return
}
