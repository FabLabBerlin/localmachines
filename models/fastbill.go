package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const FASTBILL_API_URL = "https://my.fastbill.com/api/1.0/api.php"

type FastBill struct {
	Email  string
	APIKey string
}

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

type FastBillCustomerList struct {
	Customers []FastBillCustomer
}

func (this *FastBill) GetCustomers() (*FastBillCustomerList, error) {

	var err error
	var req *http.Request
	var resp *http.Response
	var jsonBytes []byte

	type GetCustomerRequest struct {
		Service string
		Filter  []int
	}

	type GetCustomerResponseBody struct {
		Request  GetCustomerRequest
		Response FastBillCustomerList
	}

	gcReq := GetCustomerRequest{}
	gcReq.Service = "customer.get"
	jsonBytes, err = json.Marshal(gcReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal JSON: %v", err)
	}

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

	responseBody := GetCustomerResponseBody{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json: %v", err)
	}

	return &responseBody.Response, nil
}
