package controllers

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning"
	"github.com/astaxie/beego"
)

type InvoicesController struct {
	Controller
}

// @Title Get All monthly earnings
// @Description Get all monthly earnings from the database
// @Success 200 {object} models.invoices.Invoice
// @Failure	403	Failed to get all monthly earnings
// @Failure	401	Not authorized
// @router / [get]
func (this *InvoicesController) GetAll() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	mes, err := monthly_earning.GetAll()
	if err != nil {
		beego.Error("Failed to get all monthly earnings")
		this.CustomAbort(403, "Failed to get all monthly earnings")
	}

	this.Data["json"] = mes
	this.ServeJSON()
}

// @Title Get All monthly earnings
// @Description Get all monthly earnings from the database
// @Param	iid	path	int	true	"Invoice ID"
// @Success 200 string ok
// @Failure	403	Failed to delete
// @Failure	401	Not authorized
// @router /:iid [delete]
func (this *InvoicesController) Delete() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	iid, err := this.GetInt64(":iid")
	if err != nil {
		beego.Error("Failed to get iid:", err)
		this.CustomAbort(403, "Failed to delete monthly earning")
	}

	if err = monthly_earning.Delete(iid); err != nil {
		beego.Error("Failed to delete monthly earning:", err)
		this.CustomAbort(403, "Failed to delete monthly earning")
	}

	this.Data["json"] = "ok"
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

	mes, err := monthly_earning.New(locId, dbMes.Interval())
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	creationReport := monthly_earning.CreateFastbillDrafts(&mes)
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
