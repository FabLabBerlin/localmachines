package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
	"time"
)

type InvoicesController struct {
	Controller
}

// @Title Get All Invoices
// @Description Get all invoices from the database
// @Success 200 {object} models.Invoice
// @Failure	403	Failed to get all invoices
// @Failure	401	Not authorized
// @router / [get]
func (this *InvoicesController) GetAll() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	invoices, err := models.GetAllInvoices()
	if err != nil {
		beego.Error("Failed to get all invoices")
		this.CustomAbort(403, "Failed to get all invoices")
	}

	this.Data["json"] = invoices
	this.ServeJson()
}

// @Title Create Invoice
// @Description Create invoice from selection of activations
// @Param	startDate		query 	string	true		"Period start date"
// @Param	endDate		query 	string	true		"Period end date"
// @Success 200 {object} models.Invoice
// @Failure	403	Failed to create invoice
// @Failure	401	Not authorized
// @router / [post]
func (this *InvoicesController) Create() {

	// Only admin can use this API call
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error
	var startDate string
	var endDate string

	// Get variables
	startDate = this.GetString("startDate")
	if startDate == "" {
		beego.Error("Missing start date")
		this.CustomAbort(403, "Failed to create invoice")
	}

	endDate = this.GetString("endDate")
	if endDate == "" {
		beego.Error("Missing end date")
		this.CustomAbort(403, "Failed to create invoice")
	}

	beego.Trace("startDate:", startDate)
	beego.Trace("endDate:", endDate)

	// Convert / parse string time values as time.Time
	var timeForm = "2006-01-02"
	var startTime time.Time

	startTime, err = time.ParseInLocation(
		timeForm, startDate, time.Now().Location())
	if err != nil {
		beego.Error("Failed to parse startDate:", err)
		this.CustomAbort(403, "Failed to create invoice")
	}

	var endTime time.Time

	endTime, err = time.ParseInLocation(
		timeForm, endDate, time.Now().Location())
	if err != nil {
		beego.Error("Failed to parse endDate:", err)
		this.CustomAbort(403, "Failed to create invoice")
	}

	beego.Trace(startTime)
	beego.Trace(endTime)

	// Create invoice
	var invoice *models.Invoice
	invoice, err = models.CreateInvoice(startTime, endTime)
	if err != nil {
		beego.Error("Failed to create invoice:", err)
		this.CustomAbort(403, "Failed to create invoice")
	}

	this.Data["json"] = invoice
	this.ServeJson()
}
