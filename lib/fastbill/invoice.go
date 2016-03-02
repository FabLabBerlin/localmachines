package fastbill

import (
	"fmt"
	"github.com/astaxie/beego"
)

const (
	TemplateStandardId = 1
	TemplateEnglishId  = 3063
)

type Invoice struct {
	Id                   int64  `json:"INVOICE_ID,string,omitempty"`
	CustomerId           int    `json:"CUSTOMER_ID,string"`
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

type InvoiceResponse struct {
	Request  Invoice `json:"REQUEST,omitempty"`
	Response struct {
		Status    string   `json:"STATUS,omitempty"`
		InvoiceId int64    `json:"INVOICE_ID,omitempty"`
		Errors    []string `json:"ERRORS,omitempty"`
	}
}

func (inv *Invoice) Submit() (id int64, err error) {
	fb := New()

	request := Request{
		SERVICE: SERVICE_INVOICE_CREATE,
		DATA:    *inv,
	}

	var response InvoiceResponse

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
