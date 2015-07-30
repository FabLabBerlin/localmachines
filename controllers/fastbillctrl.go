package controllers

import (
	//"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type FastBillController struct {
	Controller
}

// @Title GetCustomers
// @Description Get FastBill customers
// @Success 200 {object} models.FastBillCustomerList
// @Failure 500 Internal Server Error
// @Failure 401 Not authorized
// @router  /customer [get]
func (this *FastBillController) GetCustomers() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	fb := models.FastBill{}
	fb.Email = beego.AppConfig.String("fastbillemail")
	fb.APIKey = beego.AppConfig.String("fastbillapikey")

	/*
		filter := models.FastBillCustomerGetFilter{}
		filter.CUSTOMER_ID = "1556512"
	*/

	fastBillCustomers, err := fb.GetCustomers()
	if err != nil {
		beego.Error("Failed to get FastBill customers:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = fastBillCustomers
	this.ServeJson()
}
