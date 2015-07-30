package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

// API endpoint - all requests go here
const FASTBILL_API_URL = "https://my.fastbill.com/api/1.0/api.php"

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
	CUSTOMER_TYPE                  string
	TOP                            string
	NEWSLETTER_OPTIN               string
	CONTACT_ID                     string
	ORGANIZATION                   string
	POSITION                       string
	SALUTATION                     string
	FIRST_NAME                     string
	LAST_NAME                      string
	ADDRESS                        string
	ADDRESS_2                      string
	ZIPCODE                        string
	CITY                           string
	COUNTRY_CODE                   string
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
	request.SERVICE = "customer.get"
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

// Reusable helper function for the
func (this *FastBill) execGetRequest(request *FastBillRequest, response *FastBillCustomerGetResponse) error {

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
		return fmt.Errorf("Failed to get response: %v", err)
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal json: %v", err)
	}

	return nil
}
