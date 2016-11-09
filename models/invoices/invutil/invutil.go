/*
invutil package facilitates high-level invoicing functions.
*/
package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/settings"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships/inv_user_memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"time"
)

type Invoice struct {
	invoices.Invoice
	User               *users.User                                   `json:",omitempty"`
	InvUserMemberships []*inv_user_memberships.InvoiceUserMembership `json:",omitempty"`
	Purchases          purchases.Purchases                           `json:",omitempty"`
	Sums               *Sums                                         `json:",omitempty"`
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
	ms, err := inv_user_memberships.GetForInvoice(inv.Id)
	if err != nil {
		return fmt.Errorf("GetUserMemberships: %v", err)
	}

	return inv.calculateTotals(ms)
}

func (inv *Invoice) calculateTotals(ms []*inv_user_memberships.InvoiceUserMembership) (err error) {
	inv.Sums = &Sums{}

	for _, purchase := range inv.Purchases {
		inv.Sums.Purchases.Undiscounted += purchase.TotalPrice
		inv.Sums.Purchases.PriceInclVAT += purchase.DiscountedTotal
	}
	p := (100.0 + inv.VatPercent) / 100.0
	inv.Sums.Purchases.PriceExclVAT = inv.Sums.Purchases.PriceInclVAT / p
	inv.Sums.Purchases.PriceVAT = inv.Sums.Purchases.PriceInclVAT - inv.Sums.Purchases.PriceExclVAT

	for _, m := range ms {
		if inv.UserMembershipGetsBilledHere(m.UserMembership) {
			inv.Sums.Memberships.Undiscounted += m.Membership().MonthlyPrice
			inv.Sums.Memberships.PriceInclVAT += m.Membership().MonthlyPrice
		}
		inv.Sums.Memberships.PriceExclVAT = inv.Sums.Memberships.PriceInclVAT / p
		inv.Sums.Memberships.PriceVAT = inv.Sums.Memberships.PriceInclVAT - inv.Sums.Memberships.PriceExclVAT
	}
	inv.Sums.All.Undiscounted = inv.Sums.Purchases.Undiscounted + inv.Sums.Memberships.Undiscounted
	inv.Sums.All.PriceInclVAT = inv.Sums.Purchases.PriceInclVAT + inv.Sums.Memberships.PriceInclVAT
	inv.Sums.All.PriceExclVAT = inv.Sums.Purchases.PriceExclVAT + inv.Sums.Memberships.PriceExclVAT
	inv.Sums.All.PriceVAT = inv.Sums.Purchases.PriceVAT + inv.Sums.Memberships.PriceVAT

	if inv.Invoice.Total != inv.Sums.All.PriceInclVAT {
		inv.Invoice.Total = inv.Sums.All.PriceInclVAT
		go func() {
			inv.SaveTotal()
		}()
	}

	return
}

func (inv *Invoice) InvoiceUserMemberships(data *PrefetchedData) (err error) {
	beego.Info("InvoiceUserMemberships(..)")
	beego.Info("  len(inv.InvUserMemberships)=", len(inv.InvUserMemberships))
	umbs, ok := data.UmbsByUid[inv.UserId]
	if !ok {
		umbs = []*user_memberships.UserMembership{}
	}

	for _, um := range umbs {
		invoiced := false

		//fmt.Printf("\ninv[%v/%v] ", inv.Month, inv.Year)
		if !inv.UserMembershipActiveHere(um) {
			//fmt.Printf(" not active here\n")
			continue
		} else {
			//fmt.Printf(" active here\n")
		}

		for _, ium := range inv.InvUserMemberships {
			if ium.UserMembershipId == um.Id {
				invoiced = true
				break
			}
		}

		if !invoiced {
			if inv.Status != "draft" {
				beego.Error("invoice doesn't have status draft but not all user memberships are associated")
				continue
			}

			ium, err := inv_user_memberships.Create(um, inv.Id)
			if err != nil {
				return fmt.Errorf("inv_user_memberships.Create: %v", err)
			}

			inv.InvUserMemberships = append(inv.InvUserMemberships, ium)
			data.IumbsByUid[ium.UserId] = append(data.IumbsByUid[ium.UserId], ium)
		}
	}

	return
}

func (inv *Invoice) UserMembershipActiveHere(um *user_memberships.UserMembership) bool {
	return um.ActiveAt(inv.Interval().TimeFrom()) ||
		um.ActiveAt(inv.Interval().TimeFrom().AddDate(0, 0, 15)) ||
		um.ActiveAt(inv.Interval().TimeTo())
}

func (inv *Invoice) UserMembershipGetsBilledHere(um *user_memberships.UserMembership) bool {
	//invFrom := inv.Interval().TimeFrom()
	quarterDay := -6 * time.Hour
	invTo := inv.Interval().TimeTo().Add(quarterDay)
	//fmt.Printf("'invTo'=%v\n", invTo)

	if um.StartDate.AddDate(0, 0, 2).After(invTo) {
		return false
	}

	if um.TerminationDateDefined() {
		//fmt.Printf("aaaa\n")
		if um.TerminationDate.Before(invTo) {
			//fmt.Printf("bbbb\n")
			return false
		}
	}

	//fmt.Printf("cccc\n")
	//fmt.Printf("um.TerminationDate=%v\n", um.TerminationDate)
	return um.ActiveAt(invTo)
}

func (inv *Invoice) Load() (err error) {
	data := NewPrefetchedData(inv.LocationId)

	data.UsersById[inv.UserId], err = users.GetUser(inv.UserId)
	if err != nil {
		return fmt.Errorf("get user(id=%v): %v", inv.UserId, err)
	}

	data.PurchasesByInv[inv.Id], err = purchases.GetByInvoiceId(inv.Id)
	if err != nil {
		return fmt.Errorf("get purchases by invoice id: %v", err)
	}

	data.InvUserMembershipsByInv[inv.Id], err = inv_user_memberships.GetForInvoice(inv.Id)
	if err != nil {
		return fmt.Errorf("get user memberships for invoice: %v", err)
	}

	return inv.load(*data)
}

func (inv *Invoice) load(data PrefetchedData) (err error) {
	var ok bool

	if inv.User, ok = data.UsersById[inv.UserId]; !ok {
		return fmt.Errorf("user # %v not found", inv.UserId)
	}
	if inv.Purchases, ok = data.PurchasesByInv[inv.Id]; !ok {
		inv.Purchases = []*purchases.Purchase{}
	}
	if inv.InvUserMemberships, ok = data.InvUserMembershipsByInv[inv.Id]; !ok {
		inv.InvUserMemberships = make([]*inv_user_memberships.InvoiceUserMembership, 0, 5)
	}
	/*for _, umb := range inv.InvUserMemberships {
		bill := inv.membershipGetsBilledHere(umb)
		umb.Bill = &bill
	}*/
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

func GetDraft(locId, uid int64, t time.Time) (inv *Invoice, err error) {
	iv, err := invoices.GetDraft(locId, uid, t)
	if err != nil {
		return nil, fmt.Errorf("get invoice draft entity: %v", err)
	}
	tmp, err := toUtilInvoices(iv.LocationId, []*invoices.Invoice{
		iv,
	})
	if tmp[0].VatPercent < 0.01 {
		beego.Error("detected zero vat")
	}

	return tmp[0], err
}

func Get(id int64) (inv *Invoice, err error) {
	iv, err := invoices.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get invoice entity: %v", err)
	}
	tmp, err := toUtilInvoices(iv.LocationId, []*invoices.Invoice{
		iv,
	})
	if tmp[0].VatPercent < 0.01 {
		beego.Error("detected zero vat")
	}

	return tmp[0], err
}

func GetAllAt(locId int64) (invs []*Invoice, err error) {
	ivs, err := invoices.GetAllInvoices(locId)
	if err != nil {
		return nil, fmt.Errorf("invoices.GetAllInvoices: %v", err)
	}

	if invs, err = toUtilInvoices(locId, ivs); err != nil {
		return nil, fmt.Errorf("to util invoices: %v", err)
	}

	return
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

// toUtilInvoices converts []invoices.Invoice into []invutil.Invoice with all
// fields deeply populated.
//
// What it does:
//
// 1. Bulk load associated data
// 2. Go through each invoice
// 2a. Attach data
// 2b. Update monthly total
// (TODO: 2c. Update invoice_user_memberships according to user_memberships)
//
// TODO: convert into multiple functions
func toUtilInvoices(locId int64, ivs []*invoices.Invoice) (invs []*Invoice, err error) {
	invs = make([]*Invoice, 0, len(ivs))

	data := NewPrefetchedData(locId)

	if err := data.Prefetch(); err != nil {
		return nil, fmt.Errorf("prefetch data: %v", err)
	}

	for _, iv := range ivs {
		inv := &Invoice{
			Invoice: *iv,
		}
		if _, userInLocation := data.UsersById[inv.UserId]; !userInLocation {
			beego.Error("user", inv.UserId, "is not in location", inv.LocationId,
				"(referenced through invoice", inv.Id, ") - skipping invoice!")
			continue
		}

		if err = inv.load(*data); err != nil {
			return nil, fmt.Errorf("load (invoice %v): %v", inv.Id, err)
		}
		invs = append(invs, inv)
	}

	locSettings, err := settings.GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("get settings: %v", err)
	}
	var vatPercent float64
	if vat := locSettings.GetFloat(locId, settings.VAT); vat != nil {
		vatPercent = *vat
	} else {
		vatPercent = 19.0
	}

	for _, inv := range invs {
		if err = inv.InvoiceUserMemberships(data); err != nil {
			return nil, fmt.Errorf("invoice user memberships: %v", err)
		}

		for _, p := range inv.Purchases {
			if m, ok := data.MsById[p.MachineId]; ok {
				p.Machine = m
			}

			iumbs, ok := data.IumbsByUid[p.UserId]
			if !ok {
				iumbs = []*inv_user_memberships.InvoiceUserMembership{}
			}

			for _, iumb := range iumbs {
				mbId := iumb.MembershipId
				mb, ok := data.MbsById[mbId]
				if !ok {
					return nil, fmt.Errorf("Unknown membership id: %v", mbId)
				}

				if iumb.UserMembership.ActiveAt(p.TimeStart) &&
					iumb.InvoiceId == inv.Id {
					p.Memberships = append(p.Memberships, mb)
				}
			}
		}

		inv.VatPercent = vatPercent
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

		if err = inv.calculateTotals(inv.InvUserMemberships); err != nil {
			return nil, fmt.Errorf("calculate totals: %v", err)
		}
	}

	return
}

func AssureUserHasDraftFor(locId int64, u *users.User, year int, month time.Month) error {
	ivs, err := invoices.GetAllOfUserAt(locId, u.Id)
	if err != nil {
		return fmt.Errorf("invoices.GetAllOfUserAt: %v", err)
	}
	var draft *invoices.Invoice

	for _, iv := range ivs {
		if iv.Year == year && iv.Month == int(month) {
			draft = iv
		}
	}

	if draft == nil {
		var newIv invoices.Invoice

		newIv.Year = year
		newIv.Month = int(month)
		newIv.Status = "draft"
		newIv.UserId = u.Id
		newIv.LocationId = locId

		if _, err := invoices.Create(&newIv); err != nil {
			return fmt.Errorf("invoices.Create for user %v: %v", u.Id, err)
		}

		if year == time.Now().Year() && month == time.Now().Month() {
			if err := newIv.SetCurrent(); err != nil {
				return fmt.Errorf("set current: %v", err)
			}
		}
	} else {
		if year == time.Now().Year() && month == time.Now().Month() {
			if err := draft.SetCurrent(); err != nil {
				return fmt.Errorf("set current: %v", err)
			}
		}
	}

	return nil
}
