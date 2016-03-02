package fastbill

import (
	"github.com/astaxie/beego"
)

const (
	API_URL                 = "https://my.fastbill.com/api/1.0/api.php"
	SERVICE_CUSTOMER_GET    = "customer.get"
	SERVICE_CUSTOMER_CREATE = "customer.create"
	SERVICE_CUSTOMER_UPDATE = "customer.update"
	SERVICE_CUSTOMER_DELETE = "customer.delete"
	SERVICE_INVOICE_CREATE  = "invoice.create"
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
