/*
invutil package facilitates high-level invoicing functions.
*/
package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/settings"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"time"
)

type Invoice struct {
	invoices.Invoice
	User            *users.User            `json:",omitempty"`
	UserMemberships *user_memberships.List `json:",omitempty"`
	Purchases       purchases.Purchases    `json:",omitempty"`
	Sums            *Sums                  `json:",omitempty"`
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
	ms, err := user_memberships.GetForInvoice(inv.Id)
	if err != nil {
		return fmt.Errorf("GetUserMemberships: %v", err)
	}

	return inv.calculateTotals(ms)
}

func (inv *Invoice) calculateTotals(ms *user_memberships.List) (err error) {
	inv.Sums = &Sums{}

	for _, purchase := range inv.Purchases {
		inv.Sums.Purchases.Undiscounted += purchase.TotalPrice
		inv.Sums.Purchases.PriceInclVAT += purchase.DiscountedTotal
	}
	p := (100.0 + inv.VatPercent) / 100.0
	inv.Sums.Purchases.PriceExclVAT = inv.Sums.Purchases.PriceInclVAT / p
	inv.Sums.Purchases.PriceVAT = inv.Sums.Purchases.PriceInclVAT - inv.Sums.Purchases.PriceExclVAT

	for _, m := range ms.Data {
		if m.Interval().Contains(inv.Interval().TimeTo()) {
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

	if inv.Invoice.Total != inv.Sums.All.PriceInclVAT {
		inv.Invoice.Total = inv.Sums.All.PriceInclVAT
		go func() {
			inv.SaveTotal()
		}()
	}

	return
}

func (inv *Invoice) Load() (err error) {
	usersById := make(map[int64]*users.User)
	purchasesByInv := make(map[int64][]*purchases.Purchase)
	userMembershipsByInv := make(map[int64]*user_memberships.List)

	usersById[inv.UserId], err = users.GetUser(inv.UserId)
	if err != nil {
		return fmt.Errorf("get user(id=%v): %v", inv.UserId, err)
	}

	purchasesByInv[inv.Id], err = purchases.GetByInvoiceId(inv.Id)
	if err != nil {
		return fmt.Errorf("get purchases by invoice id: %v", err)
	}

	userMembershipsByInv[inv.Id], err = user_memberships.GetForInvoice(inv.Id)
	if err != nil {
		return fmt.Errorf("get user memberships for invoice: %v", err)
	}

	return inv.load(usersById, purchasesByInv, userMembershipsByInv)
}

func (inv *Invoice) load(
	usersById map[int64]*users.User,
	purchasesByInv map[int64][]*purchases.Purchase,
	userMembershipsByInv map[int64]*user_memberships.List,
) (err error) {
	var ok bool

	if inv.User, ok = usersById[inv.UserId]; !ok {
		return fmt.Errorf("user # %v not found", inv.UserId)
	}
	if inv.Purchases, ok = purchasesByInv[inv.Id]; !ok {
		inv.Purchases = []*purchases.Purchase{}
	}
	if inv.UserMemberships, ok = userMembershipsByInv[inv.Id]; !ok {
		inv.UserMemberships = &user_memberships.List{
			Data: []*user_memberships.Combo{},
		}
	}
	for _, umb := range inv.UserMemberships.Data {
		bill := umb.Interval().Contains(inv.Interval().TimeTo())
		umb.Bill = &bill
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
	if tmp[0].VatPercent < 0.01 {
		return nil, fmt.Errorf("detected zero vat")
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

func toUtilInvoices(locId int64, ivs []*invoices.Invoice) (invs []*Invoice, err error) {
	t0 := time.Now()
	invs = make([]*Invoice, 0, len(ivs))

	usersById := make(map[int64]*users.User)
	purchasesByInv := make(map[int64][]*purchases.Purchase)
	userMembershipsByInv := make(map[int64]*user_memberships.List)

	if us, err := users.GetAllUsersAt(locId); err == nil {
		for _, u := range us {
			usersById[u.Id] = u
		}
	} else {
		return nil, fmt.Errorf("get all users: %v", err)
	}

	if ps, err := purchases.GetAllAt(locId); err == nil {
		for _, p := range ps {
			if _, ok := purchasesByInv[p.InvoiceId]; !ok {
				purchasesByInv[p.InvoiceId] = make([]*purchases.Purchase, 0, 20)
			}
			purchasesByInv[p.InvoiceId] = append(purchasesByInv[p.InvoiceId], p)
		}
	} else {
		return nil, fmt.Errorf("get all purchases: %v", err)
	}

	if umbs, err := user_memberships.GetAllAtList(locId); err == nil {
		for _, umb := range umbs.Data {
			if _, ok := userMembershipsByInv[umb.InvoiceId]; !ok {
				userMembershipsByInv[umb.InvoiceId] = &user_memberships.List{
					Data: make([]*user_memberships.Combo, 0, 3),
				}
			}
			userMembershipsByInv[umb.InvoiceId].Data = append(userMembershipsByInv[umb.InvoiceId].Data, umb)
		}
	} else {
		return nil, fmt.Errorf("get all user memberships: %v", err)
	}

	for _, iv := range ivs {
		inv := &Invoice{
			Invoice: *iv,
		}
		err = inv.load(usersById, purchasesByInv, userMembershipsByInv)
		if err != nil {
			return nil, fmt.Errorf("load: %v", err)
		}
		invs = append(invs, inv)
	}
	beego.Info("STAGE 1:", time.Now().Sub(t0))
	t1 := time.Now()
	ms, err := machine.GetAllAt(locId)
	if err != nil {
		return nil, fmt.Errorf("get all machines at %v: %v", locId, err)
	}
	msById := make(map[int64]*machine.Machine)
	for _, m := range ms {
		msById[m.Id] = m
	}

	umbsByUid := make(map[int64][]*user_memberships.UserMembership)
	if umbs, err := user_memberships.GetAllAt(locId); err == nil {
		for _, umb := range umbs {
			uid := umb.UserId
			if _, ok := umbsByUid[uid]; !ok {
				umbsByUid[uid] = []*user_memberships.UserMembership{
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
	if mbs, err := memberships.GetAllAt(locId); err == nil {
		for _, mb := range mbs {
			mbsById[mb.Id] = mb
		}
	} else {
		return nil, fmt.Errorf("Failed to get memberships: %v", err)
	}
	beego.Info("stage 2:", time.Now().Sub(t1))

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
		for _, p := range inv.Purchases {
			if m, ok := msById[p.MachineId]; ok {
				p.Machine = m
			}

			umbs, ok := umbsByUid[p.UserId]
			if !ok {
				umbs = []*user_memberships.UserMembership{}
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
		/*if inv.UserMemberships, ok = userMembershipsByInv[inv.Id]; !ok {
			inv.UserMemberships = &user_memberships.List{
				Data: []*user_memberships.Combo{},
			}
		}*/
		if err = inv.calculateTotals(inv.UserMemberships); err != nil {
			return nil, fmt.Errorf("calculate totals: %v", err)
		}
	}

	return
}

func AssureUsersHaveInvoiceFor(locId int64, year int, month time.Month) error {
	us, err := users.GetAllUsersAt(locId)
	if err != nil {
		return fmt.Errorf("get all users: %v", err)
	}

	ivs, err := invoices.GetAllInvoicesBetween(locId, year, int(month))
	if err != nil {
		return fmt.Errorf("invoices.GetAllInvoicesBetween: %v", err)
	}

	ivsByUserId := make(map[int64]invoices.Invoice)
	for _, iv := range ivs {
		ivsByUserId[iv.UserId] = *iv
	}

	for _, u := range us {
		if _, ok := ivsByUserId[u.Id]; !ok {
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
		}
	}

	return nil
}
