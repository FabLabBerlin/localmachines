package controllers

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/astaxie/beego"
	"time"
)

type InvoicesController struct {
	Controller
}

// @Title Get All Invoices
// @Description Get all invoices from the database
// @Success 200 {object} models.invoices.Invoice
// @Failure	403	Failed to get all invoices
// @Failure	401	Not authorized
// @router / [get]
func (this *InvoicesController) GetAll() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	invoices, err := invoices.GetAll()
	if err != nil {
		beego.Error("Failed to get all invoices")
		this.CustomAbort(403, "Failed to get all invoices")
	}

	this.Data["json"] = invoices
	this.ServeJSON()
}

// @Title Get All Invoices
// @Description Get all invoices from the database
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
		this.CustomAbort(403, "Failed to delete invoice")
	}

	if err = invoices.Delete(iid); err != nil {
		beego.Error("Failed to delete invoice:", err)
		this.CustomAbort(403, "Failed to delete invoice")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title Create Invoice
// @Description Create invoice from selection of activations
// @Param	startDate		query 	string	true		"Period start date"
// @Param	endDate		query 	string	true		"Period end date"
// @Success 200 {object} models.invoices.Invoice
// @Failure	403	Failed to create invoice
// @Failure	401	Not authorized
// @router / [post]
func (this *InvoicesController) Create() {

	// Only local admin can use this API call
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	startTime, endTime, err := this.parseParams()
	if err != nil {
		beego.Error("Request parameters:", err)
		this.CustomAbort(400, "Bad request")
	}

	invoices, err := invoices.Create(locId, startTime, endTime)
	if err != nil {
		beego.Error("Failed to create invoices:", err)
		this.CustomAbort(403, "Failed to create invoices")
	}

	this.Data["json"] = invoices
	this.ServeJSON()
}

// @Title Create Invoice
// @Description Create invoice from selection of activations
// @Param	startDate		query 	string	true		"Period start date"
// @Param	endDate		query 	string	true		"Period end date"
// @Success 200 {object} models.invoices.Invoice
// @Failure	401	Not authorized
// @Failure	403	Failed to create invoice
// @router /create_drafts [post]
func (this *InvoicesController) CreateDrafts() {

	// Only local admin can use this API call
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	startTime, endTime, err := this.parseParams()
	if err != nil {
		beego.Error("Request parameters:", err)
		this.CustomAbort(400, "Bad request")
	}

	invs, err := invoices.Create(locId, startTime, endTime)
	if err != nil {
		beego.Error("Failed to create invoices:", err)
		this.CustomAbort(500, "Failed to create invoice drafts")
	}

	creationReport := invoices.CreateFastbillDrafts(invs)
	beego.Info("created invoice drafts with IDs", creationReport.Ids)

	this.Data["json"] = creationReport
	this.ServeJSON()
}

func (this *InvoicesController) parseParams() (startTime, endTime time.Time, err error) {
	// Get variables
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

	// Convert / parse string time values as time.Time
	var timeForm = "2006-01-02 15:04:05"

	// Enhance start date
	startDate = fmt.Sprintf("%s 00:00:00", startDate)

	startTime, err = time.Parse(timeForm, startDate)
	if err != nil {
		err = fmt.Errorf("Failed to parse startDate:", err)
		return
	}

	// Enhance end date, make all the day inclusive
	endDate = fmt.Sprintf("%s 23:59:59", endDate)

	endTime, err = time.Parse(timeForm, endDate)
	if err != nil {
		err = fmt.Errorf("Failed to parse endDate:", err)
		return
	}

	return
}
