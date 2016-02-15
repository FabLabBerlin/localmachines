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

	// Get variables
	startDate := this.GetString("startDate")
	if startDate == "" {
		beego.Error("Missing start date")
		this.CustomAbort(403, "Failed to create invoice")
	}

	endDate := this.GetString("endDate")
	if endDate == "" {
		beego.Error("Missing end date")
		this.CustomAbort(403, "Failed to create invoice")
	}

	// Convert / parse string time values as time.Time
	var timeForm = "2006-01-02 15:04:05"

	// Enhance start date
	startDate = fmt.Sprintf("%s 00:00:00", startDate)

	startTime, err := time.Parse(timeForm, startDate)
	if err != nil {
		beego.Error("Failed to parse startDate:", err)
		this.CustomAbort(403, "Failed to create invoice")
	}

	// Enhance end date, make all the day inclusive
	endDate = fmt.Sprintf("%s 23:59:59", endDate)

	endTime, err := time.Parse(timeForm, endDate)
	if err != nil {
		beego.Error("Failed to parse endDate:", err)
		this.CustomAbort(403, "Failed to create invoice")
	}

	invoice, err := invoices.Create(locId, startTime, endTime)
	if err != nil {
		beego.Error("Failed to create invoice:", err)
		this.CustomAbort(403, "Failed to create invoice")
	}

	this.Data["json"] = invoice
	this.ServeJSON()
}
