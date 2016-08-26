package billing

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships/auto_extend"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

// @Title MigrateFbNos
// @Description Migrate fastbill numbers for invoices
// @Success 200 Success
// @Failure	500	Internal Server Error
// @router /migrate_fb_nos [get]
func (this *Controller) MigrateFbNos() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}
	if err := this.migrateFbNos(locId); err != nil {
		beego.Error(err.Error())
		this.Abort("500")
	}
	this.ServeJSON()
}

func (this *Controller) migrateFbNos(locId int64) (err error) {
	us, err := users.GetAllUsersAt(locId)
	if err != nil {
		return
	}
	for _, u := range us {
		if u.ClientId > 0 {
			fbId, err := fastbill.GetCustomerId(*u)
			if err != nil {
				return fmt.Errorf("get customer id: %v", err)
			}
			l, err := fastbill.ListInvoices(fbId)
			if err != nil {
				return fmt.Errorf("list invoices: %v", err)
			}
			invs, err := invoices.GetAllOfUserAt(locId, u.Id)
			if err != nil {
				return fmt.Errorf("get invoices of user at location: %v", err)
			}
			for _, inv := range invs {
				for _, fbIv := range l {
					if m, y, _, err := fbIv.ParseTitle(); err == nil {
						if inv.Month == m && inv.Year == y && inv.Status == fbIv.Type {
							inv.FastbillId = fbIv.Id
							inv.FastbillNo = fbIv.InvoiceNumber
							if err := inv.Save(); err == nil {
								beego.Info("Successfully update inv", inv.Id)
							} else {
								beego.Error("Failed syncing", fbIv.InvoiceTitle)
							}
							if err = invutil.FastbillSync(locId, u); err != nil {
								beego.Error("Failing to sync rest of the fields")
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// @Title Migrate
// @Description Migrate purchases/memberships to non-null invoice IDs
// @Success 200 Success
// @Failure	500	Internal Server Error
// @router /migrate [get]
func (this *Controller) Migrate() {
	auto_extend.Lock()
	defer auto_extend.Unlock()

	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	beego.Info("Loading data...")

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

	ums, err := user_memberships.GetAllAt(locId)
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
		u, ok := usrsById[userId]
		if !ok {
			beego.Error("didn't find user with id", userId)
		}
		inv := &invoices.Invoice{
			LocationId: locId,
			Month:      int(month),
			Year:       year,
			CustomerNo: u.ClientId,
			UserId:     u.Id,
			Status:     "draft",
		}

		if inv.Year < 2016 || inv.Month < 3 {
			inv.Status = "outgoing"
		}

		if inv.Year == time.Now().Year() &&
			inv.Month == int(time.Now().Month()) {
			inv.Current = true
		}

		if _, err = o.Insert(inv); err != nil {
			beego.Error("Create new user invoice:", err)
			this.Abort("500")
		}

		newInvoices = append(newInvoices, inv)
		return *inv
	}

	beego.Info("Checking each purchase...")
	for _, p := range ps {
		if p.InvoiceId > 0 {
			continue
		}
		if p.UserId == 0 {
			continue
		}
		if p.TimeStart.IsZero() {
			continue
		}

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
			invsByYearMonthUserId[year][month][p.UserId] = invs
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

	beego.Info("Checking each user membership...")
	for _, umOrig := range ums {
		fmt.Printf("o")
		if umOrig.InvoiceId > 0 {
			continue
		}

		if _, ok := usrsById[umOrig.UserId]; !ok {
			beego.Error("no user has ID", umOrig.UserId)
			continue
		}

		if umOrig.StartDate.After(umOrig.EndDate) && umOrig.Id == 301 {
			beego.Error("ignoring outdated h.m. membership with start > end")
			continue
		}

		/* << if um.StartDate.Before(inv.Interval().TimeTo()) &&
			um.EndDate.After(inv.Interval().TimeFrom()) {
			>>
		}*/

		t := umOrig.EndDate
		startMonth := umOrig.StartDate.Month()
		startYear := umOrig.StartDate.Year()

		// "Multiply"
		for i := 0; ; i++ {
			fmt.Printf(".")
			future := false
			if t.Year() > 2016 {
				future = true
			} else if t.Year() == 2016 &&
				int(t.Month()) > int(time.Now().Month()) {
				future = true
			}

			if !future {
				newUm := &user_memberships.UserMembership{
					UserId:       umOrig.UserId,
					MembershipId: umOrig.MembershipId,
					StartDate:    umOrig.StartDate,
					EndDate:      umOrig.EndDate,
					AutoExtend:   umOrig.AutoExtend,
				}

				month := t.Month()
				year := t.Year()

				invs := invsByYearMonthUserId[year][month][newUm.UserId]
				switch len(invs) {
				case 0:
					invs = []invoices.Invoice{
						newInvoice(newUm.UserId, year, month),
					}
					invsByYearMonthUserId[year][month][newUm.UserId] = invs
					fallthrough
				case 1:
					if newUm.InvoiceId == 0 {
						inv := invs[0]
						newUm.InvoiceId = inv.Id
					}
				default:
					beego.Error("Matched", len(invs), "invoices to user m'ship",
						newUm.Id)
					this.Abort("500")
				}

				if _, err = o.Insert(newUm); err != nil {
					beego.Error("Insert new user membership:", err)
					this.Abort("500")
				}
			}

			if t.Month() == startMonth && t.Year() == startYear {
				break
			} else {
				t = t.AddDate(0, -1, 0)
			}

			if i > 100 {
				panic(fmt.Sprintf("i = %v", i))
			}
		}

		if _, err := o.Delete(umOrig); err != nil {
			beego.Error("error deleting original user membership:", err)
			this.Abort("500")
		}
	}

	// Persist items

	// 1. Purchases
	beego.Info("Persisting purchases...")
	for _, p := range ps {
		if _, err = o.Update(p); err != nil {
			beego.Error("Update purchase:", err)
			this.Abort("500")
		}
	}

	if err := o.Commit(); err != nil {
		beego.Error("Commit tx:", err)
		this.Abort("500")
	}

	this.ServeJSON()
}
