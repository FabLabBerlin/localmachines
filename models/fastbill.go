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
	Email  string
	APIKey string
}

// Base FastBill API request model
type FastBillRequest struct {
	LIMIT    int64
	OFFSET   int64
	SERVICE  string
	REQUEST  interface{}
	FILTER   interface{}
	DATA     interface{}
	RESPONSE interface{}
	ERRORS   interface{}
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
		STATUS      string
		CUSTOMER_ID int64
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
	CUSTOMER_ID                    string
	CUSTOMER_NUMBER                string
	DAYS_FOR_PAYMENT               string
	CREATED                        string
	PAYMENT_TYPE                   string
	BANK_NAME                      string
	BANK_ACCOUNT_NUMBER            string
	BANK_CODE                      string
	BANK_ACCOUNT_OWNER             string
	BANK_IBAN                      string
	BANK_BIC                       string
	BANK_ACCOUNT_MANDATE_REFERENCE string
	SHOW_PAYMENT_NOTICE            string
	ACCOUNT_RECEIVABLE             string
	CUSTOMER_TYPE                  string // Required. Customer type: business | consumer
	TOP                            string
	NEWSLETTER_OPTIN               string
	CONTACT_ID                     string
	ORGANIZATION                   string // Company name [REQUIRED] when CUSTOMER_TYPE = business
	POSITION                       string
	SALUTATION                     string
	FIRST_NAME                     string
	LAST_NAME                      string // Last name [REQUIRED] when CUSTOMER_TYPE = consumer
	ADDRESS                        string
	ADDRESS_2                      string
	ZIPCODE                        string
	CITY                           string
	COUNTRY_CODE                   interface{}
	SECONDARY_ADDRESS              string
	PHONE                          string
	PHONE_2                        string
	FAX                            string
	MOBILE                         string
	EMAIL                          string
	VAT_ID                         string
	CURRENCY_CODE                  string
	LASTUPDATE                     string
	TAGS                           string
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

	request := FastBillRequest{}
	request.SERVICE = FASTBILL_SERVICE_CUSTOMER_GET
	request.LIMIT = limit
	request.OFFSET = offset
	request.FILTER = filter

	response := FastBillCustomerGetResponse{}
	err := this.execGetRequest(&request, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute get request: %v", err)
	}

	return &response.RESPONSE, nil
}

// Create FastBill customer, returns Customer ID
func (this *FastBill) CreateCustomer(customer *FastBillCustomer) (int64, error) {

	request := FastBillRequest{}
	request.SERVICE = FASTBILL_SERVICE_CUSTOMER_CREATE
	request.DATA = customer

	response := FastBillCustomerCreateResponse{}
	err := this.execGetRequest(&request, &response)
	if err != nil {
		return 0, fmt.Errorf("Failed to execute get request: %v", err)
	}

	if response.RESPONSE.STATUS == "success" {
		return response.RESPONSE.CUSTOMER_ID, nil
	} else {
		beego.Error(response.RESPONSE.STATUS)
		return 0, errors.New("There was an error while creating a customer")
	}
}

// Update FastBill customer
func (this *FastBill) UpdateCustomer(customer *FastBillCustomer) (int64, error) {

	request := FastBillRequest{}
	request.SERVICE = FASTBILL_SERVICE_CUSTOMER_UPDATE
	request.DATA = customer

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

	request := FastBillRequest{}
	request.SERVICE = FASTBILL_SERVICE_CUSTOMER_DELETE
	request.DATA = FastBillCustomer{CUSTOMER_ID: strconv.FormatInt(customerId, 10)}

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

	var err error
	var req *http.Request
	var resp *http.Response
	var jsonBytes []byte

	jsonBytes, err = json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON: %v", err)
	}
	beego.Trace(string(jsonBytes))

	req, err = http.NewRequest("GET", FASTBILL_API_URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(this.Email, this.APIKey)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		beego.Error("Failed to get response:", err)
		return fmt.Errorf("Failed to get response")
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("Failed to read response body:", err)
		return fmt.Errorf("Failed to read response body")
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		beego.Error("Failed to unmarshal JSON:", err)
		return fmt.Errorf("Failed to unmarshal json")
	}

	return nil
}
