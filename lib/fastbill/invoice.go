package fastbill

import (
	"fmt"
	"github.com/astaxie/beego"
)

const (
	TemplateStandardId = 1
	TemplateEnglishId  = 3063
)

const (
	INVOICE_TYPE_DRAFT    = "draft"
	INVOICE_TYPE_OUTGOING = "outgoing"
	INVOICE_TYPE_CREDIT   = "credit"
)

var ErrInvoiceAlreadyExported = fmt.Errorf("invoice has already been exported")

type Invoice struct {
	CustomerNumber int64  `json:"-"`
	Month          string `json:"-"`

	Id                   int64  `json:"INVOICE_ID,string,omitempty"`
	CustomerId           int64  `json:"CUSTOMER_ID,string"`
	CustomerCostcenterId int64  `json:"CUSTOMER_COSTCENTER_ID,string,omitempty"`
	CurrencyCode         string `json:"CURRENCY_CODE,omitempty"`
	TemplateId           int64  `json:"TEMPLATE_ID,string,omitempty"`
	Introtext            string `json:"INTROTEXT,omitempty"`
	InvoiceTitle         string `json:"INVOICE_TITLE,omitempty"`
	InvoiceDate          string `json:"INVOICE_DATE,omitempty"`
	DeliveryDate         string `json:"DELIVERY_DATE,omitempty"`
	CashDiscountPercent  string `json:"CASH_DISCOUNT_PERCENT,omitempty"`
	CashDiscountDays     string `json:"CASH_DISCOUNT_DAYS,omitempty"`
	EuDelivery           string `json:"EU_DELIVERY,omitempty"`
	Items                []Item `json:"ITEMS,"`
}

type InvoiceCreateResponse struct {
	Request  Invoice `json:"REQUEST,omitempty"`
	Response struct {
		Status    string   `json:"STATUS,omitempty"`
		InvoiceId int64    `json:"INVOICE_ID,omitempty"`
		Errors    []string `json:"ERRORS,omitempty"`
	}
}

type InvoiceGetResponse struct {
	Request  Request `json:"REQUEST,omitempty"`
	Response struct {
		Invoices []struct {
			Id int64 `json:"INVOICE_ID,string"`
		} `json:"INVOICES,omitempty"`
	} `json:"RESPONSE,omitempty"`
}

type InvoiceFilter struct {
	InvoiceTitle string `json:"INVOICE_TITLE,omitempty"`
	Type         string `json:"TYPE,omitempty"`
}

func (inv *Invoice) AlreadyExported() (yes bool, err error) {
	fb := New()
	filter := InvoiceFilter{}
	if filter.InvoiceTitle, err = inv.GetTitle(); err != nil {
		return false, fmt.Errorf("get title: %v", err)
	}
	request := Request{
		SERVICE: SERVICE_INVOICE_GET,
		FILTER:  filter,
		LIMIT:   10,
	}
	var response InvoiceGetResponse
	if err = fb.execGetRequest(&request, &response); err != nil {
		return false, fmt.Errorf("get request: %v", err)
	}
	n := len(response.Response.Invoices)
	if n > 1 {
		return true, fmt.Errorf("%v duplicate invoices found")
	}
	return n == 1, nil
}

func (inv *Invoice) GetTitle() (title string, err error) {
	if inv.Month == "" {
		return "", fmt.Errorf("empty month")
	}
	if inv.CustomerNumber <= 0 {
		return "", fmt.Errorf("empty customer number")
	}
	title = fmt.Sprintf("%v Invoice for Customer Number %v",
		inv.Month, inv.CustomerNumber)
	return
}

func (inv *Invoice) Submit() (id int64, err error) {
	fb := New()
	alreadyExported, err := inv.AlreadyExported()
	if err != nil {
		return 0, fmt.Errorf("checking if already exported: %v", err)
	}
	if alreadyExported {
		return 0, ErrInvoiceAlreadyExported
	}

	if inv.InvoiceTitle, err = inv.GetTitle(); err != nil {
		return 0, fmt.Errorf("get title: %v", err)
	}

	request := Request{
		SERVICE: SERVICE_INVOICE_CREATE,
		DATA:    *inv,
	}

	var response InvoiceCreateResponse

	if err := fb.execGetRequest(&request, &response); err != nil {
		return 0, fmt.Errorf("fb request: %v", err)
	}

	beego.Info("response response:", response.Response)

	inv.Id = response.Response.InvoiceId

	return inv.Id, nil
}

type Item struct {
	ArticleNumber int64   `json:"ARTICLE_NUMBER,string,omitempty"`
	Description   string  `json:"DESCRIPTION,"`
	Quantity      float64 `json:"QUANTITY,string,omitempty"`
	UnitPrice     float64 `json:"UNIT_PRICE,string"`
	VatPercent    int64   `json:"VAT_PERCENT,string"`
	IsGross       string  `json:"IS_GROSS,omitempty"`
	SortOrder     string  `json:"SORT_ORDER,"`
}
