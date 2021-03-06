package billing

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/invoices/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/settings"
	"github.com/astaxie/beego"
	"io"
	"os"
)

// @Title Get All monthly earnings
// @Description Get all monthly earnings from the database
// @Success 200 []*monthly_earning.MonthlyEarning
// @Failure	403	Failed to get all monthly earnings
// @Failure	401	Not authorized
// @router /monthly_earnings [get]
func (this *Controller) GetAll() {

	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	mes, err := monthly_earning.GetAllAt(locId)
	if err != nil {
		beego.Error("Failed to get all monthly earnings")
		this.Abort("500")
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
// @router /monthly_earnings [post]
func (this *Controller) Create() {

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
// @router /monthly_earnings/:iid/create_drafts [post]
func (this *Controller) CreateDrafts() {

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

	locSettings, err := settings.GetAllAt(me.LocationId)
	if err != nil {
		beego.Error("get settings:", err)
		this.Abort("500")
	}

	var vatPercent float64
	if vat := locSettings.GetFloat(me.LocationId, settings.VAT); vat != nil {
		vatPercent = *vat
	} else {
		vatPercent = 19.0
	}

	creationReport := monthly_earning.CreateFastbillDrafts(me, vatPercent)
	beego.Info("created invoice drafts with IDs", creationReport.Ids)

	this.Data["json"] = creationReport
	this.ServeJSON()
}

func (this *Controller) parseParams() (interval lib.Interval, err error) {
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
// @router /monthly_earnings/:id/download_excel [get]
func (this *Controller) DownloadExcelExport() {

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
		beego.Error("open", me.FilePath, ":", err)
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
