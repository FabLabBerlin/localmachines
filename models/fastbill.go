package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// API endpoint - all requests go here
const (
	FASTBILL_API_URL                 = "https://my.fastbill.com/api/1.0/api.php"
	FASTBILL_SERVICE_CUSTOMER_GET    = "customer.get"
	FASTBILL_SERVICE_CUSTOMER_CREATE = "customer.create"
	FASTBILL_SERVICE_CUSTOMER_UPDATE = "customer.update"
	FASTBILL_SERVICE_CUSTOMER_DELETE = "customer.delete"
	FASTBILL_CUSTOMER_TYPE_BUSINESS  = "business"
	FASTBILL_CUSTOMER_TYPE_CONSUMER  = "consumer"
)

// Main FastBill object. All the functionality goes through this object.
type FastBill struct {
	email  string
	apiKey string
}

func NewFastBill() (fb *FastBill) {
	fb = &FastBill{
		email:  beego.AppConfig.String("fastbillemail"),
		apiKey: beego.AppConfig.String("fastbillapikey"),
	}
	return
}

// Base FastBill API request model
type FastBillRequest struct {
	LIMIT   int64 `json:",omitempty"`
	OFFSET  int64 `json:",omitempty"`
	SERVICE string
	FILTER  *FastBillCustomerGetFilter `json:",omitempty"`
	DATA    interface{}
}

// customer.get response model
// For now define a response per request. The interface{} thing
// does not work well with the JSON unmarshal thing for some reason.
// But this is more clear in a way.
type FastBillCustomerGetResponse struct {
	REQUEST  FastBillRequest
	RESPONSE FastBillCustomerList
}

// Response model that we expect from the FastBill API on customer.create
// request
type FastBillCustomerCreateResponse struct {
	REQUEST  FastBillRequest
	RESPONSE struct {
		ERRORS      []string `json:",omitempty"`
		STATUS      string   `json:",omitempty"`
		CUSTOMER_ID int64    `json:",omitempty"`
	}
}

// Model received from FastBill on customer.update
type FastBillCustomerUpdateResponse struct {
	REQUEST  FastBillRequest
	RESPONSE struct {
		STATUS      string
		CUSTOMER_ID string
	}
}

// FastBill customer.delete response model
type FastBillCustomerDeleteResponse struct {
	REQUEST  FastBillRequest
	RESPONSE struct {
		STATUS string
	}
}

// Create customer response model returned to THIS API clients
type FastBillCreateCustomerResponse struct {
	CUSTOMER_ID int64
}

// Update customer response model returned to THIS API clients
type FastBillUpdateCustomerResponse struct {
	CUSTOMER_ID int64
}

// FastBill customer model
type FastBillCustomer struct {
	CUSTOMER_ID                    string `json:",omitempty"`
	CUSTOMER_NUMBER                string `json:",omitempty"`
	DAYS_FOR_PAYMENT               string `json:",omitempty"`
	CREATED                        string `json:",omitempty"`
	PAYMENT_TYPE                   string `json:",omitempty"`
	BANK_NAME                      string `json:",omitempty"`
	BANK_ACCOUNT_NUMBER            string `json:",omitempty"`
	BANK_CODE                      string `json:",omitempty"`
	BANK_ACCOUNT_OWNER             string `json:",omitempty"`
	BANK_IBAN                      string `json:",omitempty"`
	BANK_BIC                       string `json:",omitempty"`
	BANK_ACCOUNT_MANDATE_REFERENCE string `json:",omitempty"`
	SHOW_PAYMENT_NOTICE            string `json:",omitempty"`
	ACCOUNT_RECEIVABLE             string `json:",omitempty"`
	CUSTOMER_TYPE                  string // Required. Customer type: business | consumer
	TOP                            string `json:",omitempty"`
	NEWSLETTER_OPTIN               string `json:",omitempty"`
	CONTACT_ID                     string `json:",omitempty"`
	ORGANIZATION                   string `json:",omitempty"` // Company name [REQUIRED] when CUSTOMER_TYPE = business
	POSITION                       string `json:",omitempty"`
	SALUTATION                     string `json:",omitempty"`
	FIRST_NAME                     string `json:",omitempty"`
	LAST_NAME                      string // Last name [REQUIRED] when CUSTOMER_TYPE = consumer
	ADDRESS                        string `json:",omitempty"`
	ADDRESS_2                      string `json:",omitempty"`
	ZIPCODE                        string `json:",omitempty"`
	CITY                           string `json:",omitempty"`
	COUNTRY_CODE                   string `json:",omitempty"`
	SECONDARY_ADDRESS              string `json:",omitempty"`
	PHONE                          string `json:",omitempty"`
	PHONE_2                        string `json:",omitempty"`
	FAX                            string `json:",omitempty"`
	MOBILE                         string `json:",omitempty"`
	EMAIL                          string `json:",omitempty"`
	VAT_ID                         string `json:",omitempty"`
	CURRENCY_CODE                  string `json:",omitempty"`
	LASTUPDATE                     string `json:",omitempty"`
	TAGS                           string `json:",omitempty"`
}

// Customer list model
type FastBillCustomerList struct {
	Customers []FastBillCustomer
}

// Filter model for customer.get request
type FastBillCustomerGetFilter struct {
	CUSTOMER_ID     string
	CUSTOMER_NUMBER string
	COUNTRY_CODE    string
	CITY            string
	TERM            string // Search term in one of the given fields: ORGANIZATION, FIRST_NAME, LAST_NAME, ADDRESS, ADDRESS_2, ZIPCODE, EMAIL, TAGS.
}

// Get customers with support for limit and offset
func (this *FastBill) GetCustomers(filter *FastBillCustomerGetFilter,
	limit int64, offset int64) (*FastBillCustomerList, error) {

	request := FastBillRequest{
		SERVICE: FASTBILL_SERVICE_CUSTOMER_GET,
		LIMIT:   limit,
		OFFSET:  offset,
		FILTER:  filter,
	}

	response := FastBillCustomerGetResponse{}
	err := this.execGetRequest(&request, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute get request: %v", err)
	}

	return &response.RESPONSE, nil
}

// Create FastBill customer, returns Customer ID
func (this *FastBill) CreateCustomer(customer *FastBillCustomer) (int64, error) {

	request := FastBillRequest{
		SERVICE: FASTBILL_SERVICE_CUSTOMER_CREATE,
		DATA:    customer,
	}

	response := FastBillCustomerCreateResponse{}
	err := this.execGetRequest(&request, &response)
	if err != nil {
		return 0, fmt.Errorf("Failed to execute get request: %v", err)
	}

	if response.RESPONSE.STATUS == "success" {
		return response.RESPONSE.CUSTOMER_ID, nil
	} else {
		errInfo := fmt.Sprintf("%v", response.RESPONSE.STATUS)
		if response.RESPONSE.ERRORS != nil && len(response.RESPONSE.ERRORS) > 0 {
			errInfo += " " + strings.Join(response.RESPONSE.ERRORS, ",")
		}
		beego.Error(errInfo)
		return 0, errors.New("There was an error while creating a customer: " + errInfo)
	}
}

// Update FastBill customer
func (this *FastBill) UpdateCustomer(customer *FastBillCustomer) (int64, error) {

	request := FastBillRequest{
		SERVICE: FASTBILL_SERVICE_CUSTOMER_UPDATE,
		DATA:    customer,
	}

	response := FastBillCustomerUpdateResponse{}
	err := this.execGetRequest(&request, &response)
	if err != nil {
		beego.Error("Failed to execute get request:", err)
		return 0, fmt.Errorf("Failed to execute get request")
	}

	if response.RESPONSE.STATUS == "success" {
		var idInt int64
		idInt, err = strconv.ParseInt(response.RESPONSE.CUSTOMER_ID, 10, 64)
		if err != nil {
			beego.Error("Failed to parse int:", err)
			return 0, fmt.Errorf("Failed to parse int")
		}
		return idInt, nil
	} else {
		beego.Error(response.RESPONSE.STATUS)
		return 0, errors.New("There was an error while updating customer")
	}
}

func (this *FastBill) DeleteCustomer(customerId int64) error {

	request := FastBillRequest{
		SERVICE: FASTBILL_SERVICE_CUSTOMER_DELETE,
		DATA: FastBillCustomer{
			CUSTOMER_ID: strconv.FormatInt(customerId, 10),
		},
	}

	response := FastBillCustomerDeleteResponse{}
	err := this.execGetRequest(&request, &response)
	if err != nil {
		return fmt.Errorf("Failed to execute get request: %v", err)
	}

	if response.RESPONSE.STATUS == "success" {
		return nil
	} else {
		return errors.New("There was an error while deleting a customer")
	}

}

// Reusable helper function for the
func (this *FastBill) execGetRequest(request *FastBillRequest, response interface{}) error {

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("GET", FASTBILL_API_URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(this.email, this.apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("Failed to get response:", err)
		return fmt.Errorf("Failed to get response")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("Failed to read response body:", err)
		return fmt.Errorf("Failed to read response body")
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		beego.Error("Failed to unmarshal JSON:", err)
		beego.Error("JSON was: ", string(body))
		return fmt.Errorf("Failed to unmarshal json")
	}

	return nil
}
