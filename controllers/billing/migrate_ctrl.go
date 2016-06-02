package billing

import (
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego"
	"time"
)

// @Title Migrate
// @Description Migrate purchases/memberships to non-null invoice IDs
// @Success 200 Success
// @Failure	500	Internal Server Error
// @router /migrate [get]
func (this *Controller) Migrate() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	intv := lib.Interval{
		MonthFrom: 1,
		YearFrom:  2015,
		MonthTo:   12,
		YearTo:    2016,
	}

	ps, err := purchases.GetAllBetweenAt(locId, intv)
	if err != nil {
		beego.Error("Failed to get purchases:", err)
		this.Abort("500")
	}

	ums, err := models.GetUserMemberships

	invs, err := invoices.GetAllInvoicesAt(locId)
	if err != nil {
		beego.Error("Failed to get invoices:", err)
		this.Abort("500")
	}

	invsByYearMonth := make(map[int]map[time.Month][]invoices.Invoice)
	for _, year := range []int{2015, 2016} {
		invsByYearMonth[year] = make(map[time.Month][]invoices.Invoice)
		for m := 1; m < 12; m++ {
			month := time.Month(m)
			invsByYearMonth[year][month] = make([]invoices.Invoice, 0, 1)
		}
	}

	for _, inv := range invs {
		m := time.Month(inv.Month)
		y := inv.Year
		invsByYearMonth[y][m] = append(invsByYearMonth[y][m], inv)
	}

	for _, p := range ps {
		if p.InvoiceId > 0 {
			continue
		}

		if inv.Interval.Contains(p.TimeStart) {

		}
	}

	for _, um := range ums {
		if um.InvoiceId > 0 {
			continue
		}

		if m.StartDate.Before(inv.Interval.TimeTo()) && m.EndDate.After(inv.Interval.TimeFrom()) {

		}

		// "Multiply"
		// ...
	}

	this.ServeJSON()
}
