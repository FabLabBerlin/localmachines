package invoices

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"time"
)

const TABLE_NAME = "invoices"

// Invoice represents an actual or future invoice. Future invoices do not
// have a FastbillId.
type Invoice struct {
	Id          int64  `json:",omitempty"`
	LocationId  int64  `json:",omitempty"`
	FastbillId  int64  `json:",omitempty"`
	FastbillNo  string `json:",omitempty"`
	Month       int
	Year        int
	CustomerId  int64
	CustomerNo  int64
	UserId      int64
	Total       float64
	Status      string
	User        *users.User         `orm:"-" json:",omitempty"`
	Purchases   purchases.Purchases `orm:"-" json:",omitempty"`
	Sums        *Sums               `orm:"-" json:",omitempty"`
	VatPercent  float64             `json:",omitempty"`
	Canceled    bool
	InvoiceDate time.Time
	PaidDate    time.Time
	DueDate     time.Time
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

	memberships, err := models.GetUserMemberships(inv.User.Id)
	if err != nil {
		return fmt.Errorf("GetUserMemberships: %v", err)
	}
	if inv.User.Id == 57 {
		beego.Info("CalculateTotals: len(memberships) = ", len(memberships.Data))
	}
	for _, m := range memberships.Data {
		if inv.User.Id == 57 {
			beego.Info("m=", m)
			beego.Info("inv.Interval=", inv.Interval)
		}
		if m.StartDate.Before(inv.Interval().TimeFrom()) &&
			m.EndDate.After(inv.Interval().TimeTo()) {
			if inv.User.Id == 57 {
				beego.Info("if => true")
			}
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
	if inv.User.Id == 57 {
		beego.Info("CalculateTotals: inv.Sums.Purchases.PriceInclVAT=", inv.Sums.Purchases.PriceInclVAT)
		beego.Info("CalculateTotals: inv.Sums.Memberships.PriceInclVAT=", inv.Sums.Memberships.PriceInclVAT)
		beego.Info("CalculateTotals: inv.Sums.All.PriceInclVAT=", inv.Sums.All.PriceInclVAT)
	}

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
		iv := &Invoice{
			Month:      int(t.Month()),
			Year:       t.Year(),
			Purchases:  make([]*purchases.Purchase, 0, 20),
			User:       inv.User,
			VatPercent: inv.VatPercent,
		}
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
		inv := Invoice{
			LocationId:  locId,
			FastbillId:  fbInv.Id,
			FastbillNo:  fbInv.InvoiceNumber,
			CustomerId:  fbInv.CustomerId,
			Status:      fbInv.Type,
			Total:       fbInv.Total,
			VatPercent:  fbInv.VatPercent,
			Canceled:    fbInv.Canceled(),
			DueDate:     fbInv.DueDate(),
			InvoiceDate: fbInv.InvoiceDate(),
			PaidDate:    fbInv.PaidDate(),
		}
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
		if _, err := CreateOrUpdate(&inv); err != nil {
			return fmt.Errorf("Failed to create or update inv: %v", err)
		}
	}
	return
}
