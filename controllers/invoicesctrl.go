package controllers

import (
	"github.com/astaxie/beego"
	//"github.com/kr15h/fabsmith/models"
)

type InvoicesController struct {
	Controller
}

// @Title Create Invoice
// @Description Create invoice from selection of activations
// @Param	startDate		query 	string	true		"Period start date"
// @Param	endDate		query 	string	true		"Period end date"
// @Param	userId		query 	int	true		"User ID"
// @Param	includeInvoiced		query 	bool	true		"Whether to include already invoiced activations"
// @Success 200 {object} models.Invoice
// @Failure	403	Failed to create invoice
// @Failure	401	Not authorized
// @router / [post]
func (this *InvoicesController) Create() {

	beego.Info("Invoices controller")

}
