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

// Main FastBill object
type FastBill struct {
	Email  string
	APIKey string
}

// Basic FastBill API get request
type FastBillGetRequest struct {
	Service string
	Filter  interface{}
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

// Customer list
type FastBillCustomerList struct {
	Customers []FastBillCustomer
}

// Filter model for customer.get request
type FastBillCustomerGetFilter struct {
	CUSTOMER_ID     string
	CUSTOMER_NUMBER string
	COUNTRY_CODE    string
	CITY            string
	TERM            string
}

type FastBillCustomerGetResponse struct {
	Request  FastBillGetRequest
	Response FastBillCustomerList
}

func (this *FastBill) GetCustomers(filter *FastBillCustomerGetFilter) (*FastBillCustomerList, error) {

	var err error
	var req *http.Request
	var resp *http.Response
	var jsonBytes []byte

	gcReq := FastBillGetRequest{}
	gcReq.Service = "customer.get"
	gcReq.Filter = filter
	jsonBytes, err = json.Marshal(gcReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal JSON: %v", err)
	}
	beego.Trace(string(jsonBytes))

	req, err = http.NewRequest("GET", FASTBILL_API_URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(this.Email, this.APIKey)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get response: %v", err)
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	customerGetResponse := FastBillCustomerGetResponse{}
	err = json.Unmarshal(body, &customerGetResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json: %v", err)
	}

	return &customerGetResponse.Response, nil
}
