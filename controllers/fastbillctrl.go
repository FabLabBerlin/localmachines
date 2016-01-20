package controllers

import (
	"github.com/FabLabBerlin/localmachines/models/fastbill"
	"github.com/astaxie/beego"
	"strconv"
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
// @Success 200 {object}
// @Failure 500 Internal Server Error
// @Failure 401 Not authorized
// @router  /customer [get]
func (this *FastBillController) GetCustomers() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	fb := fastbill.New()

	filter := fastbill.CustomerGetFilter{
		TERM: this.GetString("term"),
	}

	customerId, err := this.GetInt64("customerid")
	if err != nil {
		beego.Warning("Failed to get customer ID.")
	} else {
		filter.CUSTOMER_ID = strconv.FormatInt(customerId, 10)
	}

	customerNumber, err := this.GetInt64("customernumber")
	if err != nil {
		beego.Warning("Failed to get customer number.")
	} else {
		filter.CUSTOMER_NUMBER = strconv.FormatInt(customerNumber, 10)
	}

	limit, err := this.GetInt64("limit")
	if err != nil {
		beego.Warning("Failed to get limit")
		limit = 0
	}

	offset, err := this.GetInt64("offset")
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
// @Success 200 {object}
// @Failure 500 Internal Server Error
// @Failure 401 Not authorized
// @router  /customer [post]
func (this *FastBillController) CreateCustomer() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	fb := fastbill.New()

	customer := fastbill.Customer{
		FIRST_NAME:   this.GetString("firstname"),
		LAST_NAME:    this.GetString("lastname"),
		EMAIL:        this.GetString("email"),
		PHONE:        this.GetString("phone"),
		ADDRESS:      this.GetString("address"),
		CITY:         this.GetString("city"),
		ZIPCODE:      this.GetString("zipcode"),
		ORGANIZATION: this.GetString("organization"),
		COUNTRY_CODE: this.GetString("countrycode"),
	}

	if customer.ORGANIZATION == "" {
		customer.CUSTOMER_TYPE = fastbill.CUSTOMER_TYPE_CONSUMER
	} else {
		customer.CUSTOMER_TYPE = fastbill.CUSTOMER_TYPE_BUSINESS
	}

	customerId, err := fb.CreateCustomer(&customer)
	if err != nil {
		beego.Error("Failed to create FastBill customer:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = fastbill.CreateCustomerResponse{
		CUSTOMER_ID: customerId,
	}
	this.ServeJson()
}

// @Title UpdateCustomer
// @Description Update FastBill customer
// @Param customerid   path    int     true      "Customer ID"
// @Param firstname     query   string  false     "First name of the customer"
// @Param lastname      query   string  false     "Last name of the customer"
// @Param email         query   string  false     "Customer email"
// @Param address       query   string  false     "Customer billing address"
// @Param city          query   string  false     "Customer city"
// @Param countrycode   query   string  false     "Customer ISO 3166 ALPHA-2"
// @Param zipcode       query   string  false     "Customer zip code"
// @Param phone         query   string  false     "Customer phone number"
// @Param organization  query   string  false     "Organization of the customer"
// @Success 200 {object}
// @Failure 500 Internal Server Error
// @Failure 401 Not authorized
// @router  /customer/:customerid [put]
func (this *FastBillController) UpdateCustomer() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	fb := fastbill.New()

	customerId, err := this.GetInt64(":customerid")
	if err != nil {
		beego.Error("Failed to get :customerid")
		this.CustomAbort(500, "Internal Server Error")
	}

	// Get existing customer
	filter := fastbill.CustomerGetFilter{
		CUSTOMER_ID: strconv.FormatInt(customerId, 10),
	}

	fastBillCustomers, err := fb.GetCustomers(&filter, 0, 0)
	if err != nil {
		beego.Error("Failed to get FastBill customers:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	numCustomersGot := len(fastBillCustomers.Customers)
	if numCustomersGot <= 0 {
		beego.Error("Failed to get customer")
		this.CustomAbort(500, "Internal Server Error")
	}

	// Update only the fields that have new values
	customer := fastBillCustomers.Customers[0]

	if fn := this.GetString("firstname"); fn != "" {
		customer.FIRST_NAME = fn
	}

	if ln := this.GetString("lastname"); ln != "" {
		customer.LAST_NAME = ln
	}

	// TODO: Check email address.
	if email := this.GetString("email"); email != "" {
		customer.EMAIL = email
	}

	if phone := this.GetString("phone"); phone != "" {
		customer.PHONE = phone
	}

	if addr := this.GetString("address"); addr != "" {
		customer.ADDRESS = addr
	}

	if city := this.GetString("city"); city != "" {
		customer.CITY = city
	}

	if zip := this.GetString("zipcode"); zip != "" {
		customer.ZIPCODE = zip
	}

	// The organization can be empty
	customer.ORGANIZATION = this.GetString("organization")

	if cc := this.GetString("countrycode"); cc != "" {
		customer.COUNTRY_CODE = cc
	}

	// If there is no organization - customer can be considered a plain consumer
	if customer.ORGANIZATION == "" {
		customer.CUSTOMER_TYPE = fastbill.CUSTOMER_TYPE_CONSUMER
	} else {
		customer.CUSTOMER_TYPE = fastbill.CUSTOMER_TYPE_BUSINESS
	}

	customerIdUpd, err := fb.UpdateCustomer(&customer)
	if err != nil {
		beego.Error("Failed to update FastBill customer:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = fastbill.UpdateCustomerResponse{
		CUSTOMER_ID: customerIdUpd,
	}
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

	fb := fastbill.New()

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
