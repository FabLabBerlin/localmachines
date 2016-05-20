package controllers

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"io"
	"os"
	"time"
)

type InvoicesController struct {
	Controller
}

// @Title Get All monthly earnings
// @Description Get all monthly earnings from the database
// @Success 200 []*monthly_earning.MonthlyEarning
// @Failure	403	Failed to get all monthly earnings
// @Failure	401	Not authorized
// @router / [get]
func (this *InvoicesController) GetAll() {

	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	mes, err := monthly_earning.GetAllAt(locId)
	if err != nil {
		beego.Error("Failed to get all monthly earnings")
		this.CustomAbort(403, "Failed to get all monthly earnings")
	}

	this.Data["json"] = mes
	this.ServeJSON()
}

// @Title Create monthly earning
// @Description Create monthly earning from date range
// @Param	startDate	query 	string	true	"Period start date"
// @Param	endDate		query 	string	true	"Period end date"
// @Success 200 {object}
// @Failure	403	Failed to create monthly earning
// @Failure	401	Not authorized
// @router / [post]
func (this *InvoicesController) Create() {

	// Only local admin can use this API call
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	interval, err := this.parseParams()
	if err != nil {
		beego.Error("Request parameters:", err)
		this.CustomAbort(400, "Bad request")
	}

	mes, err := monthly_earning.Create(locId, interval)
	if err != nil {
		beego.Error("Failed to create monthly earnings:", err)
		this.CustomAbort(403, "Failed to create monthly earnings")
	}

	this.Data["json"] = mes
	this.ServeJSON()
}

// @Title Create Fastbill drafts
// @Description Create fastbill drafts from date range
// @Param	startDate	query 	string	true	"Period start date"
// @Param	endDate		query 	string	true	"Period end date"
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	403	Failed to create Fastbill drafts
// @router /:iid/create_drafts [post]
func (this *InvoicesController) CreateDrafts() {

	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	if locId != 1 {
		this.CustomAbort(400, "Fastbill not activated for this location")
	}

	if !this.IsSuperAdmin() {
		beego.Error("User must be super admin")
		this.CustomAbort(401, "Not authorized")
	}

	iid, err := this.GetInt64(":iid")
	if err != nil {
		beego.Error("Failed to get iid:", err)
		this.CustomAbort(400, "Bad request")
	}

	dbMes, err := monthly_earning.Get(iid)
	if err != nil {
		beego.Error("invoices get:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	me, err := monthly_earning.New(locId, dbMes.Interval())
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	creationReport := monthly_earning.CreateFastbillDrafts(me)
	beego.Info("created invoice drafts with IDs", creationReport.Ids)

	this.Data["json"] = creationReport
	this.ServeJSON()
}

func (this *InvoicesController) parseParams() (interval lib.Interval, err error) {
	startDate := this.GetString("startDate")
	if startDate == "" {
		err = fmt.Errorf("Missing start date")
		return
	}

	endDate := this.GetString("endDate")
	if endDate == "" {
		err = fmt.Errorf("Missing end date")
		return
	}

	interval, err = lib.NewInterval(startDate, endDate)
	if err != nil {
		err = fmt.Errorf("parse interval: %v", err)
		return
	}

	return
}

// @Title DownloadExcelExport
// @Description Download existing excel export
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:id/download_excel [get]
func (this *InvoicesController) DownloadExcelExport() {

	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get id:", err)
		this.CustomAbort(400, "Bad request")
	}

	me, err := monthly_earning.Get(id)
	if err != nil {
		beego.Error("invoices get:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdminAt(me.LocationId) {
		this.CustomAbort(401, "Not authorized")
	}

	f, err := os.Open(me.FilePath)
	if err != nil {
		beego.Error("open ", me.FilePath, ":", err)
		if os.IsNotExist(err) {
			this.CustomAbort(404, "Cannot find Excel Export file")
		} else {
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	defer f.Close()
	this.Ctx.Output.ContentType("xlsx")
	if _, err := io.Copy(this.Ctx.ResponseWriter, f); err != nil {
		beego.Error("copy ", me.FilePath, ":", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Finish()
}

type MonthlySummary struct {
	User   users.User
	Amount float64
}

// @Title GetMonth
// @Description Get overview for a month
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month [get]
func (this *InvoicesController) GetMonth() {
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
		if err := inv.CalculateTotals(); err != nil {
			beego.Error("CalculateTotals:", err)
			this.Abort("500")
		}
		if inv.User.Id == 57 {
			beego.Info("57:")
			beego.Info("sums:", inv.Sums)
		}
		sum := MonthlySummary{
			User:   inv.User,
			Amount: inv.Sums.All.PriceInclVAT,
		}
		sums = append(sums, sum)
	}

	this.Data["json"] = sums
	this.ServeJSON()
}

// @Title GetUser
// @Description Get monthly overview for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid [get]
func (this *InvoicesController) GetUser() {
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
			User: inv.User,
		}
		for _, p := range inv.Purchases.Data {
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

// @Title GetStatus
// @Description Get status for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid/status [get]
func (this *InvoicesController) GetStatus() {
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

	fbInvs, err := inv.FetchExisting()
	if err != nil {
		beego.Error("Failed to fetch existing fastbill invoice:", err)
		this.Abort("500")
	}

	beego.Info("fbInv=", fbInvs)

	this.Data["json"] = fbInvs
	this.ServeJSON()
}
