package controllers

import (
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type FastBillController struct {
	Controller
}

// @Title GetCustomers
// @Description Get FastBill customers
// @Param limit     					query   int     false    "Limit of the number of records returned"
// @Param offset    					query   int     false    "Offset of the records returned"
// @Param term      					query   string  false    "Filter term in one of the given fields: ORGANIZATION, FIRST_NAME, LAST_NAME, ADDRESS, ADDRESS_2, ZIPCODE, EMAIL, TAGS."
// @Param customerid      		query   int  		false    "Customer ID"
// @Param customernumber      query   int  		false    "Customer Number"
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
	beego.Trace(filter.TERM)

	var err error
	var customerId int64
	customerId, err = this.GetInt64("customerid")
	if err != nil {
		beego.Warning("Failed to get customer ID.")
	} else {
		filter.CUSTOMER_ID = fmt.Sprintf("%d", customerId)
	}

	var customerNumber int64
	customerNumber, err = this.GetInt64("customernumber")
	if err != nil {
		beego.Warning("Failed to get customer number.")
	} else {
		filter.CUSTOMER_NUMBER = fmt.Sprintf("%d", customerNumber)
	}

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

// @Title CreateCustomer
// @Description Create FastBill customer
// @Param firstname     query   string  true     "First name of the customer"
// @Param lastname      query   string  true     "Last name of the customer"
// @Param email         query   string  true     "Customer email"
// @Param address       query   string  true     "Customer billing address"
// @Param city          query   string  true     "Customer city"
// @Param countrycode   query   string  true     "Customer ISO 3166 ALPHA-2"
// @Param zipcode       query   string  true     "Customer zip code"
// @Param phone         query   string  false    "Customer phone number"
// @Param organization  query   string  false    "Organization of the customer"
// @Success 200 {object} models.FastBillCreateCustomerResponse
// @Failure 500 Internal Server Error
// @Failure 401 Not authorized
// @router  /customer [post]
func (this *FastBillController) CreateCustomer() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	fb := models.FastBill{}
	fb.Email = beego.AppConfig.String("fastbillemail")
	fb.APIKey = beego.AppConfig.String("fastbillapikey")

	customer := models.FastBillCustomer{}
	customer.FIRST_NAME = this.GetString("firstname")
	customer.LAST_NAME = this.GetString("lastname")
	customer.EMAIL = this.GetString("email")
	customer.PHONE = this.GetString("phone")
	customer.ADDRESS = this.GetString("address")
	customer.CITY = this.GetString("city")
	customer.ZIPCODE = this.GetString("zipcode")
	customer.ORGANIZATION = this.GetString("organization")
	customer.COUNTRY_CODE = this.GetString("countrycode")

	if customer.ORGANIZATION == "" {
		customer.CUSTOMER_TYPE = models.FASTBILL_CUSTOMER_TYPE_CONSUMER
	} else {
		customer.CUSTOMER_TYPE = models.FASTBILL_CUSTOMER_TYPE_BUSINESS
	}

	customerId, err := fb.CreateCustomer(&customer)
	if err != nil {
		beego.Error("Failed to create FastBill customer:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	response := models.FastBillCreateCustomerResponse{}
	response.CUSTOMER_ID = customerId

	this.Data["json"] = response
	this.ServeJson()
}

// @Title DeleteCustomer
// @Description Delete existing FastBill customer
// @Param customerid	path	int	true	"Customer ID"
// @Success 200	ok
// @Failure 500 Internal Server Error
// @Failure 401 Not authorized
// @router  /customer/:customerid [delete]
func (this *FastBillController) DeleteCustomer() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	fb := models.FastBill{}
	fb.Email = beego.AppConfig.String("fastbillemail")
	fb.APIKey = beego.AppConfig.String("fastbillapikey")

	customerId, err := this.GetInt64(":customerid")
	if err != nil {
		beego.Error("Failed to get :customerid")
		this.CustomAbort(500, "Internal Server Error")
	}

	err = fb.DeleteCustomer(customerId)
	if err != nil {
		beego.Error("Failed to delete customer:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
