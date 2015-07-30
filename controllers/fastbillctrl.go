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
// @Param limit     query   int     false    "Limit of the number of records returned"
// @Param offset    query   int     false    "Offset of the records returned"
// @Param term      query   string  false    "Filter term in one of the given fields: ORGANIZATION, FIRST_NAME, LAST_NAME, ADDRESS, ADDRESS_2, ZIPCODE, EMAIL, TAGS."
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

	filter := models.FastBillCustomerGetFilter{}
	filter.TERM = this.GetString("term")

	var err error
	var limit int64
	limit, err = this.GetInt64("limit")
	if err != nil {
		beego.Warning("Failed to get limit")
		limit = 0
	}

	var offset int64
	offset, err = this.GetInt64("offset")
	if err != nil {
		beego.Warning("Failed to get offset")
		offset = 0
	}

	fastBillCustomers, err := fb.GetCustomers(&filter, limit, offset)
	if err != nil {
		beego.Error("Failed to get FastBill customers:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = fastBillCustomers
	this.ServeJson()
}
