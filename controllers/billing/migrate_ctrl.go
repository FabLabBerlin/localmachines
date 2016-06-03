package billing

import (
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

	usrs, err := users.GetAllUsersAt(locId)
	if err != nil {
		beego.Error("Failed to get users:", err)
		this.Abort("500")
	}

	usrsById := make(map[int64]*users.User)
	for _, u := range usrs {
		usrsById[u.Id] = u
	}

	ps, err := purchases.GetAllBetweenAt(locId, intv)
	if err != nil {
		beego.Error("Failed to get purchases:", err)
		this.Abort("500")
	}

	ums, err := models.GetAllUserMembershipsAt(locId)
	if err != nil {
		beego.Error("Failed to get user memberships:", err)
		this.Abort("500")
	}

	invs, err := invoices.GetAllInvoices(locId)
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
		invsByYearMonth[y][m] = append(invsByYearMonth[y][m], *inv)
	}

	newInvoices := make([]*invoices.Invoice, 0, 1000)

	for _, p := range ps {
		if p.InvoiceId > 0 {
			continue
		}

		/*if inv.Interval().Contains(p.TimeStart) {
			p.InvoiceId = inv.Id
		}*/
		month := p.TimeStart.Month()
		year := p.TimeStart.Year()

		invs := invsByYearMonth[year][month]
		switch len(invs) {
		case 0:
			u := usrsById[p.UserId]
			inv := &invoices.Invoice{
				LocationId: locId,
				Month:      int(month),
				Year:       year,
				CustomerNo: u.ClientId,
				UserId:     u.Id,
			}
			newInvoices = append(newInvoices, inv)
			invs = []invoices.Invoice{*inv}
			invsByYearMonth[year][month] = invs
			fallthrough
		case 1:
			if p.InvoiceId == 0 {
				inv := invs[0]
				p.InvoiceId = inv.Id
			}
		default:
			beego.Error("Matched", len(invs), "invoices to purchase", p.Id)
			this.Abort("500")
		}
	}

	for _, um := range ums {
		if um.InvoiceId > 0 {
			continue
		}

		/*if um.StartDate.Before(inv.Interval().TimeTo()) &&
			um.EndDate.After(inv.Interval().TimeFrom()) {

		}*/

		// "Multiply"
		// ...
	}

	// Persist it all
	o := orm.NewOrm()

	if err := o.Begin(); err != nil {
		beego.Error("Begin tx:", err)
		this.Abort("500")
	}

	// 1. Purchases
	// 2. User memberships
	// 3. New invoices

	if err := o.Commit(); err != nil {
		beego.Error("Commit tx:", err)
		this.Abort("500")
	}

	this.ServeJSON()
}
