package billing

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
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

	list := make([]invutil.Invoice, 0, len(ivs))

	for _, iv := range ivs {
		inv := invutil.Invoice{
			Invoice: *iv,
			User:    usrsById[iv.UserId],
		}
		list = append(list, inv)
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

	fbDraft, _, err := inv.FastbillCreateDraft(true)
	if err != nil {
		beego.Error("Create fastbill draft:", err)
		this.CustomAbort(500, err.Error())
	}

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
	this.assertFinanceUser()
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
		msg := fmt.Sprintf("wrong status to complete invoice: %v", s)
		beego.Error(msg)
		this.CustomAbort(400, msg)
	}

	if inv.Canceled {
		msg := "invoice already canceled"
		beego.Error(msg)
		this.CustomAbort(500, msg)
	}

	if err := inv.FastbillCancel(); err != nil {
		beego.Error("cancel:", err)
		this.CustomAbort(500, err.Error())
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
	this.assertFinanceUser()
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
		msg := fmt.Sprintf("wrong status to complete invoice: %v", s)
		beego.Error(msg)
		this.CustomAbort(400, msg)
	}

	if err := inv.CalculateTotals(); err != nil {
		beego.Error("CalculateTotals:", err)
		this.CustomAbort(500, err.Error())
	}

	if err := inv.FastbillComplete(); err != nil {
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
	this.assertFinanceUser()
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	id, err := this.GetInt64(":id")
	if err != nil {
		this.Abort("400")
	}

	canceled := this.GetString("canceled") == "true"

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

	if canceled {
		if inv.CanceledSent {
			beego.Error("canceled invoice already sent")
			this.CustomAbort(500, "canceled invoice already sent")
		}

		if err := inv.FastbillSendCanceled(); err != nil {
			beego.Error("send canceled invoice via fastbill:", err)
			this.CustomAbort(500, err.Error())
		}
	} else {
		if inv.Sent {
			beego.Error("invoice already sent")
			this.CustomAbort(500, "invoice already sent")
		}

		if err := inv.FastbillSend(); err != nil {
			beego.Error("send invoice via fastbill:", err)
			this.CustomAbort(500, err.Error())
		}
	}

	this.ServeJSON()
}

// @Title SyncFastbillInvoices
// @Description SyncFastbillInvoices
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /users/:uid/sync [get]
func (this *Controller) SyncFastbillInvoices() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.Abort("400")
	}

	u, err := users.GetUser(uid)
	if err != nil {
		beego.Error("get user:", err)
		this.Abort("500")
	}

	err = invutil.FastbillSync(locId, u)
	if err != nil {
		beego.Error("Failed to sync fastbill invoices:", err)
		this.Abort("500")
	}

	this.ServeJSON()
}

func (this *Controller) assertFinanceUser() {
	uid, err := this.GetSessionUserId()
	if err != nil {
		beego.Error(err.Error())
		this.Abort("500")
	}

	switch uid {
	case 6, 19, 28, 336:
		return
	default:
		this.Abort("403")
	}
}
