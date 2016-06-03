package billing

import (
	"github.com/FabLabBerlin/localmachines/models/memberships"
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

	usrs, err := users.GetAllUsersAt(locId)
	if err != nil {
		beego.Error("Failed to get users:", err)
		this.Abort("500")
	}

	usrsById := make(map[int64]*users.User)
	for _, u := range usrs {
		usrsById[u.Id] = u
	}

	ps, err := purchases.GetAllAt(locId)
	if err != nil {
		beego.Error("Failed to get purchases:", err)
		this.Abort("500")
	}

	ums, err := memberships.GetAllUserMembershipsAt(locId)
	if err != nil {
		beego.Error("Failed to get user memberships:", err)
		this.Abort("500")
	}

	invs, err := invoices.GetAllInvoices(locId)
	if err != nil {
		beego.Error("Failed to get invoices:", err)
		this.Abort("500")
	}

	invsByYearMonthUserId := make(map[int]map[time.Month]map[int64][]invoices.Invoice)
	for _, year := range []int{2015, 2016} {
		invsByYearMonthUserId[year] = make(map[time.Month]map[int64][]invoices.Invoice)
		for m := 1; m <= 12; m++ {
			month := time.Month(m)
			invsByYearMonthUserId[year][month] = make(map[int64][]invoices.Invoice)
			for _, u := range usrs {
				invsByYearMonthUserId[year][month][u.Id] = make([]invoices.Invoice, 0, 1)
			}
		}
	}

	for _, inv := range invs {
		m := time.Month(inv.Month)
		y := inv.Year
		invsByYearMonthUserId[y][m][inv.UserId] = append(invsByYearMonthUserId[y][m][inv.UserId], *inv)
	}

	/**********************************/
	//
	// The database is expected to be really clean because the
	// User Memberships will be duplicated for every billing month.
	// Otherwise that might become too complex for the moment.
	//

	for _, p := range ps {
		if p.InvoiceId != 0 {
			panic("purchase found with inv id")
		}
	}

	for _, um := range ums {
		if um.InvoiceId != 0 {
			panic("user membership found with inv id")
		}
	}

	if len(invs) > 0 {
		panic("Expected db to have no invoices")
	}
	/**********************************/

	o := orm.NewOrm()

	if err := o.Begin(); err != nil {
		beego.Error("Begin tx:", err)
		this.Abort("500")
	}

	newInvoices := make([]*invoices.Invoice, 0, 1000)

	newInvoice := func(userId int64, year int, month time.Month) invoices.Invoice {
		beego.Info("new invoice for user id", userId)
		u := usrsById[userId]
		inv := &invoices.Invoice{
			LocationId: locId,
			Month:      int(month),
			Year:       year,
			CustomerNo: u.ClientId,
			UserId:     u.Id,
		}

		if _, err = o.Insert(inv); err != nil {
			beego.Error("Create new user invoice:", err)
			this.Abort("500")
		}

		newInvoices = append(newInvoices, inv)
		return *inv
	}

	for _, p := range ps {
		if p.InvoiceId > 0 {
			continue
		}
		if p.UserId == 0 {
			continue
		}

		beego.Info("purchase", p.Id)

		/* << if inv.Interval().Contains(p.TimeStart) {
			p.InvoiceId = inv.Id >>
		}*/
		month := p.TimeStart.Month()
		year := p.TimeStart.Year()

		invs := invsByYearMonthUserId[year][month][p.UserId]
		switch len(invs) {
		case 0:
			invs = []invoices.Invoice{
				newInvoice(p.UserId, year, month),
			}
			beego.Info("invsByYearMonthUserId[year]=", invsByYearMonthUserId[year])
			beego.Info("invsByYearMonthUserId[year][month]=", invsByYearMonthUserId[year][month])
			beego.Info("year=", year)
			beego.Info("month=", month)
			beego.Info("invsByYearMonthUserId[year][month][p.UserId]=", invsByYearMonthUserId[year][month][p.UserId])
			invsByYearMonthUserId[year][month][p.UserId] = invs
			fallthrough
		case 1:
			if p.InvoiceId == 0 {
				inv := invs[0]
				p.InvoiceId = uint64(inv.Id)
			}
		default:
			beego.Error("Matched", len(invs), "invoices to purchase", p.Id)
			this.Abort("500")
		}
	}

	newUms := make([]*memberships.UserMembership, 0, 1000)

	for _, umCurrent := range ums {
		if umCurrent.InvoiceId > 0 {
			continue
		}

		/* << if um.StartDate.Before(inv.Interval().TimeTo()) &&
			um.EndDate.After(inv.Interval().TimeFrom()) {
			>>
		}*/

		t := umCurrent.EndDate
		startMonth := umCurrent.StartDate.Month()
		startYear := umCurrent.StartDate.Year()

		// "Multiply"
		for i := 0; ; i++ {
			var um *memberships.UserMembership

			if i == 0 {
				um = umCurrent
			} else {
				clone := *umCurrent
				um = &clone
				newUms = append(newUms, um)
			}

			month := t.Month()
			year := t.Year()

			invs := invsByYearMonthUserId[year][month][um.UserId]
			switch len(invs) {
			case 0:
				invs = []invoices.Invoice{
					newInvoice(um.UserId, year, month),
				}
				invsByYearMonthUserId[year][month][um.UserId] = invs
				fallthrough
			case 1:
				if um.InvoiceId == 0 {
					inv := invs[0]
					um.InvoiceId = inv.Id
				}
			default:
				beego.Error("Matched", len(invs), "invoices to user m'ship",
					um.Id)
				this.Abort("500")
			}

			if t.Month() == startMonth && t.Year() == startYear {
				break
			} else {
				t = t.AddDate(0, -1, 0)
			}
		}
	}

	// Persist items

	// 1. Purchases
	for _, p := range ps {
		if _, err = o.Update(p); err != nil {
			beego.Error("Update purchase:", err)
			this.Abort("500")
		}
	}

	// 2. User memberships
	for _, um := range ums {
		if _, err = o.Update(um); err != nil {
			beego.Error("Update user membership:", err)
			this.Abort("500")
		}
	}

	// 3. New User memberships
	for _, newUm := range newUms {
		if _, err = o.Insert(&newUm); err != nil {
			beego.Error("Insert new user membership:", err)
			this.Abort("500")
		}
	}

	if err := o.Commit(); err != nil {
		beego.Error("Commit tx:", err)
		this.Abort("500")
	}

	this.ServeJSON()
}
