// FastBill API Wrapper
package fastbill

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

// API endpoint - all requests go here
const (
	CUSTOMER_TYPE_BUSINESS = "business"
	CUSTOMER_TYPE_CONSUMER = "consumer"
)

// customer.get response model
// For now define a response per request. The interface{} thing
// does not work well with the JSON unmarshal thing for some reason.
// But this is more clear in a way.
type CustomerGetResponse struct {
	REQUEST  Request
	RESPONSE CustomerList
}

// Response model that we expect from the FastBill API on customer.create
// request
type CustomerCreateResponse struct {
	REQUEST  Request
	RESPONSE struct {
		ERRORS      []string `json:",omitempty"`
		STATUS      string   `json:",omitempty"`
		CUSTOMER_ID int64    `json:",omitempty"`
	}
}

// Model received from FastBill on customer.update
type CustomerUpdateResponse struct {
	REQUEST  Request
	RESPONSE struct {
		STATUS      string
		CUSTOMER_ID string
	}
}

// FastBill customer.delete response model
type CustomerDeleteResponse struct {
	REQUEST  Request
	RESPONSE struct {
		STATUS string
	}
}

// Create customer response model returned to THIS API clients
type CreateCustomerResponse struct {
	CUSTOMER_ID int64
}

// Update customer response model returned to THIS API clients
type UpdateCustomerResponse struct {
	CUSTOMER_ID int64
}

// FastBill customer model
type Customer struct {
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
type CustomerList struct {
	Customers []Customer
}

// Filter model for customer.get request
type CustomerGetFilter struct {
	CUSTOMER_ID     string
	CUSTOMER_NUMBER string
	COUNTRY_CODE    string
	CITY            string
	TERM            string // Search term in one of the given fields: ORGANIZATION, FIRST_NAME, LAST_NAME, ADDRESS, ADDRESS_2, ZIPCODE, EMAIL, TAGS.
}

// Get customers with support for limit and offset
func (this *FastBill) GetCustomers(filter *CustomerGetFilter,
	limit int64, offset int64) (*CustomerList, error) {

	request := Request{
		SERVICE: SERVICE_CUSTOMER_GET,
		LIMIT:   limit,
		OFFSET:  offset,
		FILTER:  filter,
	}

	response := CustomerGetResponse{}
	err := this.execGetRequest(&request, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute get request: %v", err)
	}

	return &response.RESPONSE, nil
}

func GetCustomerId(user users.User) (customerId int64, err error) {
	fb := New()

	customerNumber := user.ClientId
	if customerNumber <= 0 {
		return 0, fmt.Errorf("wrong customer number (a.k.a. Fastbill ID a.k.a. Client Id) for user ID %v: %v",
			user.Id, user.ClientId)
	}
	filter := CustomerGetFilter{
		CUSTOMER_NUMBER: strconv.FormatInt(customerNumber, 10),
	}
	list, err := fb.GetCustomers(&filter, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("get customers: %v", err)
	}
	if len(list.Customers) == 0 {
		return 0, fmt.Errorf("no customer found for customer number %v",
			customerNumber)
	}
	if n := len(list.Customers); n > 1 {
		return 0, fmt.Errorf("%v matches found for customer number %v",
			n, customerNumber)
	}
	return strconv.ParseInt(list.Customers[0].CUSTOMER_ID, 10, 64)
}

// Create FastBill customer, returns Customer ID
func (this *FastBill) CreateCustomer(customer *Customer) (int64, error) {

	request := Request{
		SERVICE: SERVICE_CUSTOMER_CREATE,
		DATA:    customer,
	}

	response := CustomerCreateResponse{}
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
func (this *FastBill) UpdateCustomer(customer *Customer) (int64, error) {

	request := Request{
		SERVICE: SERVICE_CUSTOMER_UPDATE,
		DATA:    customer,
	}

	response := CustomerUpdateResponse{}
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

	request := Request{
		SERVICE: SERVICE_CUSTOMER_DELETE,
		DATA: Customer{
			CUSTOMER_ID: strconv.FormatInt(customerId, 10),
		},
	}

	response := CustomerDeleteResponse{}
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
