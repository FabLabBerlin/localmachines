package billing

import (
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices/invutil"
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

// @Title GetInvoice
// @Description Get invoice
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	403	Forbidden
// @Failure	500	Internal Server Error
// @router /invoices/:id [get]
func (this *Controller) GetInvoice() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	id, err := this.GetInt64(":id")
	if err != nil {
		this.Abort("400")
	}

	inv, err := invutil.Get(id)
	if err != nil {
		beego.Error("invutil get:", err)
		this.Abort("500")
	}

	if inv.LocationId != locId {
		this.Abort("403")
	}

	this.Data["json"] = inv
	this.ServeJSON()
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

	ivs, err := invoices.GetAllInvoicesBetween(locId, year, month)
	if err != nil {
		beego.Error("Failed to get invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	list := make([]invutil.Invoice, len(ivs))

	for i, iv := range ivs {
		list[i] = invutil.Invoice{
			Invoice: *iv,
			User:    usrsById[iv.UserId],
		}
	}

	this.Data["json"] = list
	this.ServeJSON()
}

// @Title Create draft
// @Description Create draft for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /invoices/:id/draft [post]
func (this *Controller) CreateDraft() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	id, err := this.GetInt64(":id")
	if err != nil {
		this.Abort("400")
	}

	inv, err := invutil.Get(id)
	if err != nil {
		beego.Error("invutil get:", err)
		this.Abort("500")
	}

	if inv.LocationId != locId {
		this.Abort("403")
	}

	if s := inv.Status; s != "draft" {
		beego.Error("wrong status to create fastbill draft:", s)
		this.Abort("400")
	}

	if err := inv.CalculateTotals(); err != nil {
		beego.Error("CalculateTotals:", err)
		this.CustomAbort(500, err.Error())
	}

	fbDraft, empty, err := inv.CreateFastbillDraft(true)
	if err != nil {
		beego.Error("Create fastbill draft:", err)
		this.CustomAbort(500, err.Error())
	}
	beego.Info("empty=", empty)
	beego.Info("fbDraft=", fbDraft)

	this.Data["json"] = fbDraft
	this.ServeJSON()
}

// @Title Cancel
// @Description Cancel invoice, here and in fastbill
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /invoices/:id/cancel [post]
func (this *Controller) Cancel() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	id, err := this.GetInt64(":id")
	if err != nil {
		this.Abort("400")
	}

	inv, err := invutil.Get(id)
	if err != nil {
		beego.Error("invutil get:", err)
		this.Abort("500")
	}

	if inv.LocationId != locId {
		this.Abort("403")
	}

	if s := inv.Status; s != "outgoing" {
		beego.Error("wrong status to complete invoice:", s)
		this.Abort("400")
	}

	if inv.Canceled {
		beego.Error("invoice already canceled")
		this.Abort("500")
	}

	if err := inv.Cancel(); err != nil {
		beego.Error("cancel:", err)
		this.Abort("500")
	}

	this.ServeJSON()
}

// @Title Complete
// @Description Complete draft invoice, here and in fastbill
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /invoices/:id/complete [post]
func (this *Controller) Complete() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	id, err := this.GetInt64(":id")
	if err != nil {
		this.Abort("400")
	}

	inv, err := invutil.Get(id)
	if err != nil {
		beego.Error("invutil get:", err)
		this.Abort("500")
	}

	if inv.LocationId != locId {
		this.Abort("403")
	}

	if s := inv.Status; s != "draft" {
		beego.Error("wrong status to complete invoice:", s)
		this.Abort("400")
	}

	if err := inv.CalculateTotals(); err != nil {
		beego.Error("CalculateTotals:", err)
		this.CustomAbort(500, err.Error())
	}

	if err := inv.CompleteFastbill(); err != nil {
		beego.Error("complete fastbill:", err)
		this.CustomAbort(500, err.Error())
	}

	this.ServeJSON()
}

// @Title Send
// @Description Send outging invoice (through fastbill) but remember it here
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /invoices/:id/send [post]
func (this *Controller) Send() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	id, err := this.GetInt64(":id")
	if err != nil {
		this.Abort("400")
	}

	inv, err := invutil.Get(id)
	if err != nil {
		beego.Error("invutil get:", err)
		this.Abort("500")
	}

	if inv.LocationId != locId {
		this.Abort("403")
	}

	if s := inv.Status; s != "outgoing" {
		beego.Error("wrong status to send invoice:", s)
		this.CustomAbort(400, "wrong status to send invoice: "+s)
	}

	if inv.Sent {
		beego.Error("invoice already sent")
		this.CustomAbort(500, "invoice already sent")
	}

	if err := inv.Send(); err != nil {
		beego.Error("send invoice via fastbill:", err)
		this.CustomAbort(500, err.Error())
	}

	this.ServeJSON()
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

	err = invutil.SyncFastbillInvoices(locId, year, time.Month(month))
	if err != nil {
		beego.Error("Failed to sync fastbill invoices:", err)
		this.Abort("500")
	}

	this.ServeJSON()
}
