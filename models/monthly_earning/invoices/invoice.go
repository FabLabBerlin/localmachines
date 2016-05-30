package invoices

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

const TABLE_NAME = "invoices"

// Invoice represents an actual or future invoice. Future invoices do not
// have a FastbillId.
type Invoice struct {
	Id         int64 `json:",omitempty"`
	FastbillId int64 `json:",omitempty"`
	FastbillNo int64 `json:",omitempty"`
	Month      int
	Year       int
	UserId     int64
	Interval   lib.Interval        `orm:"-" json:",omitempty"`
	User       *users.User         `orm:"-" json:",omitempty"`
	Purchases  purchases.Purchases `orm:"-" json:",omitempty"`
	Sums       *Sums               `orm:"-" json:",omitempty"`
	VatPercent float64             `json:",omitempty"`
}

func init() {
	orm.RegisterModel(new(Invoice))
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
		if m.StartDate.Before(inv.Interval.TimeFrom()) &&
			m.EndDate.After(inv.Interval.TimeTo()) {
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
			Interval: lib.Interval{
				MonthFrom: int(t.Month()),
				YearFrom:  t.Year(),
				MonthTo:   int(t.Month()),
				YearTo:    t.Year(),
			},
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
			if iv.Interval.Contains(p.TimeStart) {
				iv.Purchases = append(iv.Purchases, p)
			}
		}
		if err := iv.CalculateTotals(); err != nil {
			return nil, fmt.Errorf("CalculateTotals: %v", err)
		}
	}

	return
}

func (inv *Invoice) TableName() string {
	return TABLE_NAME
}

// GetIdsAndStatuses in Invoice struct.  Leaving other fields empty.
func GetIdsAndStatuses(locId int64, year, month int) ([]*Invoice, error) {
	interval := lib.Interval{
		MonthFrom: month,
		YearFrom:  year,
		MonthTo:   month,
		YearTo:    year,
	}

	query := `
SELECT *
FROM
  (SELECT user_id,
          invoice_id,
          invoice_status
   FROM purchases
   WHERE ? < time_end
     AND time_end < ?
     AND location_id = ?
     AND user_id <> 0
   UNION SELECT user_id,
                invoice_id,
                invoice_status
   FROM user_membership
   JOIN membership ON user_membership.membership_id = membership.id
   WHERE ? > start_date
     AND end_date > ?
     AND location_id = ?
     AND user_id <> 0) AS id_data
GROUP BY concat(user_id, '-', COALESCE(invoice_id, ''));
	`

	var invs []*Invoice

	o := orm.NewOrm()
	_, err := o.Raw(query,
		interval.TimeFrom(),
		interval.TimeTo(),
		locId,
		interval.TimeTo(),
		interval.TimeFrom(),
		locId).QueryRows(&invs)

	for _, inv := range invs {
		inv.Month = month
		inv.Year = year
		inv.Interval = interval
	}

	return invs, err
}
