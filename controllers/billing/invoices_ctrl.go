package billing

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"time"
)

type MonthlySummary struct {
	User          users.User
	InvoiceNumber string
	InvoiceStatus string
	Amount        float64
	CustomerId    int64
}

// @Title GetMonth
// @Description Get overview for a month
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month [get]
func (this *Controller) GetMonth() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	usrs, err := users.GetAllUsersAt(locId)
	if err != nil {
		beego.Error("Failed to get users:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	usrsById := make(map[int64]*users.User)
	for _, u := range usrs {
		usrsById[u.Id] = u
	}

	ivs, err := invoices.GetAllInvoicesAt(locId, year, month)
	if err != nil {
		beego.Error("Failed to get invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	for _, iv := range ivs {
		iv.User = usrsById[iv.UserId]
	}

	this.Data["json"] = ivs
	this.ServeJSON()
}

// @Title GetUser
// @Description Get monthly overview for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid [get]
func (this *Controller) GetUser() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt64(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.CustomAbort(400, "Bad request")
	}

	interval := lib.Interval{
		MonthFrom: int(month),
		YearFrom:  int(year),
		MonthTo:   int(month),
		YearTo:    int(year),
	}

	me, err := monthly_earning.New(locId, interval)
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	sums := make([]MonthlySummary, 0, len(me.Invoices))
	for _, inv := range me.Invoices {
		sum := MonthlySummary{
			User: *inv.User,
		}
		for _, p := range inv.Purchases {
			sum.Amount += p.DiscountedTotal
		}
		sums = append(sums, sum)
	}

	var userInv *invoices.Invoice

	for _, inv := range me.Invoices {
		if inv.User.Id == uid {
			userInv = inv
		}
	}

	if err := userInv.CalculateTotals(); err != nil {
		beego.Error("CalculateTotals:", err)
		this.Abort("500")
	}

	this.Data["json"] = userInv
	this.ServeJSON()
}

// @Title GetStatuses
// @Description Get statuses for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid/statuses [get]
func (this *Controller) GetStatuses() {
	_, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt64(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.CustomAbort(400, "Bad request")
	}

	user, err := users.GetUser(uid)
	if err != nil {
		beego.Error("Failed to get user:", err)
		this.Abort("500")
	}

	inv := fastbill.Invoice{
		Month:          time.Month(month).String(),
		Year:           int(year),
		CustomerNumber: user.ClientId,
	}

	var existingMonth fastbill.ExistingMonth

	key := fmt.Sprintf("/months/%v/%v/users/%v/statuses", year, month, user.Id)
	redis.Cached(key, 3600, &existingMonth, func() interface{} {
		ivs, err := inv.FetchExisting()
		if err != nil {
			beego.Error("Failed to fetch existing fastbill invoice:", err)
			this.Abort("500")
		}
		return *ivs
	})

	this.Data["json"] = existingMonth
	this.ServeJSON()
}

// @Title Create draft
// @Description Create draft for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid/draft [post]
func (this *Controller) CreateDraft() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt64(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.CustomAbort(400, "Bad request")
	}

	interval := lib.Interval{
		MonthFrom: int(month),
		YearFrom:  int(year),
		MonthTo:   int(month),
		YearTo:    int(year),
	}

	me, err := monthly_earning.New(locId, interval)
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	sums := make([]MonthlySummary, 0, len(me.Invoices))
	for _, inv := range me.Invoices {
		sum := MonthlySummary{
			User: *inv.User,
		}
		for _, p := range inv.Purchases {
			sum.Amount += p.DiscountedTotal
		}
		sums = append(sums, sum)
	}

	var userInv *invoices.Invoice

	for _, inv := range me.Invoices {
		if inv.User.Id == uid {
			userInv = inv
		}
	}

	if err := userInv.CalculateTotals(); err != nil {
		beego.Error("CalculateTotals:", err)
		this.Abort("500")
	}

	fbDraft, empty, err := monthly_earning.CreateFastbillDraft(me, userInv)
	if err != nil {
		beego.Error("Create fastbill draft:", err)
		this.Abort("500")
	}
	beego.Info("empty=", empty)
	beego.Info("fbDraft=", fbDraft)

	this.Data["json"] = fbDraft
	this.ServeJSON()
}

// @Title Send user invoicing data
// @Description Send invoicing data for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid/invoices/:id/send [post]
func (this *Controller) Send() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt64(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.CustomAbort(400, "Bad request")
	}
	_ = locId
	_ = year
	_ = month
	_ = uid
	//invoices.Send()

	beego.Error("Not implemented")
	this.CustomAbort(500, "Not implemented")
}

// @Title Update user invoicing data
// @Description Update invoicing data for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid/update [post]
func (this *Controller) Update() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt64(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.CustomAbort(400, "Bad request")
	}

	interval := lib.Interval{
		MonthFrom: int(month),
		YearFrom:  int(year),
		MonthTo:   int(month),
		YearTo:    int(year),
	}

	me, err := monthly_earning.New(locId, interval)
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	sums := make([]MonthlySummary, 0, len(me.Invoices))
	for _, inv := range me.Invoices {
		sum := MonthlySummary{
			User: *inv.User,
		}
		for _, p := range inv.Purchases {
			sum.Amount += p.DiscountedTotal
		}
		sums = append(sums, sum)
	}

	var userInv *invoices.Invoice

	for _, inv := range me.Invoices {
		if inv.User.Id == uid {
			userInv = inv
		}
	}

	if err := userInv.CalculateTotals(); err != nil {
		beego.Error("CalculateTotals:", err)
		this.Abort("500")
	}

	beego.Error("Not implemented")
	this.CustomAbort(500, "Not implemented")
}

// @Title SyncFastbillInvoices
// @Description SyncFastbillInvoices
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/sync [get]
func (this *Controller) SyncFastbillInvoices() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	l, err := fastbill.ListInvoices(year, time.Month(month))
	if err != nil {
		beego.Error("Failed to get invoice list from fastbill:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	usrs, err := users.GetAllUsersAt(locId)
	if err != nil {
		beego.Error("Failed to get user list:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	for _, fbInv := range l {
		inv := invoices.Invoice{
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
		if _, err := invoices.CreateOrUpdate(&inv); err != nil {
			beego.Error("Failed to create or update inv:", err)
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	this.ServeJSON()
}
