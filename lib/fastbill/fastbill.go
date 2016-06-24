package fastbill

import (
	"bytes"

	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

const (
	SERVICE_CUSTOMER_GET     = "customer.get"
	SERVICE_CUSTOMER_CREATE  = "customer.create"
	SERVICE_CUSTOMER_UPDATE  = "customer.update"
	SERVICE_CUSTOMER_DELETE  = "customer.delete"
	SERVICE_INVOICE_CREATE   = "invoice.create"
	SERVICE_INVOICE_GET      = "invoice.get"
	SERVICE_INVOICE_DELETE   = "invoice.delete"
	SERVICE_INVOICE_COMPLETE = "invoice.complete"
	SERVICE_INVOICE_CANCEL   = "invoice.cancel"
)

// For unit tests we need to change this here actually
var (
	API_URL = "https://my.fastbill.com/api/1.0/api.php"
)

// Main FastBill object. All the functionality goes through this object.
type FastBill struct {
	email  string
	apiKey string
}

func New() (fb *FastBill) {
	fb = &FastBill{
		email:  beego.AppConfig.String("fastbillemail"),
		apiKey: beego.AppConfig.String("fastbillapikey"),
	}
	return
}

// Base FastBill API request model
type Request struct {
	LIMIT   int64 `json:",omitempty"`
	OFFSET  int64 `json:",omitempty"`
	SERVICE string
	FILTER  interface{} `json:",omitempty"`
	DATA    interface{}
}

// Reusable helper function for the
func (this *FastBill) execGetRequest(request *Request, response interface{}) error {

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON: %v", err)
	}

	beego.Info("jsonBytes:", string(jsonBytes))

	req, err := http.NewRequest("GET", API_URL, bytes.NewBuffer(jsonBytes))
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

	if code := resp.StatusCode; code < 200 || code >= 300 {
		return fmt.Errorf("status code %v (%v)", code, resp.Status)
	}

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
