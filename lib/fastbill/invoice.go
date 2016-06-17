package fastbill

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
)

const (
	TemplateStandardId        = 1
	TemplateMakeaIndustriesId = 875511
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
	Year           int    `json:"-"`

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

func (this *InvoiceCreateResponse) Error() error {
	if len(this.Response.Errors) == 0 {
		return nil
	} else {
		return errors.New(strings.Join(this.Response.Errors, "; "))
	}
}

type InvoiceGetResponse struct {
	Request  Request `json:"REQUEST,omitempty"`
	Response struct {
		Invoices []InvoiceGetResponseInvoice `json:"INVOICES,omitempty"`
	} `json:"RESPONSE,omitempty"`
}

type InvoiceGetResponseInvoice struct {
	Id          int64   `json:"INVOICE_ID,string"`
	Type        string  `json:"TYPE,omitempty"`
	InvoiceDate string  `json:"INVOICE_DATE,omitempty"`
	PaidDate    string  `json:"PAID_DATE,omitempty"`
	IsCanceled  string  `json:"IS_CANCELED,omitempty"`
	Total       float64 `json:"TOTAL,omitempty"`
}

func (this *InvoiceGetResponseInvoice) Canceled() bool {
	return this.IsCanceled == "1"
}

type ExistingMonth struct {
	Invoices  []InvoiceGetResponseInvoice
	CanCancel bool
	CanDraft  bool
	CanSend   bool
}

type InvoiceFilter struct {
	InvoiceTitle string `json:"INVOICE_TITLE,omitempty"`
	Type         string `json:"TYPE,omitempty"`
}

func (inv *Invoice) FetchExisting() (existingMonth *ExistingMonth, err error) {
	fb := New()
	filter := InvoiceFilter{}
	if filter.InvoiceTitle, err = inv.GetTitle(); err != nil {
		return nil, fmt.Errorf("get title: %v", err)
	}
	request := Request{
		SERVICE: SERVICE_INVOICE_GET,
		FILTER:  filter,
		LIMIT:   10,
	}
	var response InvoiceGetResponse
	if err = fb.execGetRequest(&request, &response); err != nil {
		return nil, fmt.Errorf("get request: %v", err)
	}
	existingMonth = &ExistingMonth{
		Invoices: response.Response.Invoices,
	}
	numDrafts := 0
	numCanceledOutgoing := 0
	numUncanceledOutgoing := 0
	for _, inv := range existingMonth.Invoices {
		if inv.Type == INVOICE_TYPE_DRAFT {
			existingMonth.CanSend = true
			numDrafts++
		} else if inv.Type == INVOICE_TYPE_OUTGOING {
			if inv.Canceled() {
				numCanceledOutgoing++
			} else {
				numUncanceledOutgoing++
			}
		}
	}
	if numUncanceledOutgoing > 0 {
		existingMonth.CanCancel = true
	}
	if numDrafts == 0 && numUncanceledOutgoing == 0 {
		existingMonth.CanDraft = true
	}
	return
}

func (inv *Invoice) GetTitle() (title string, err error) {
	if inv.Month == "" {
		return "", fmt.Errorf("empty month")
	}
	if inv.Year <= 0 {
		return "", fmt.Errorf("need year")
	}
	if inv.CustomerNumber <= 0 {
		return "", fmt.Errorf("empty customer number")
	}
	title = fmt.Sprintf("%v %v Invoice for Customer Number %v",
		inv.Month, inv.Year, inv.CustomerNumber)
	return
}

func (inv *Invoice) Submit() (id int64, err error) {
	fb := New()
	fbInvs, err := inv.FetchExisting()
	if err != nil {
		return 0, fmt.Errorf("checking if already exported: %v", err)
	}
	beego.Info("lib.fastbill:Invoice#Submit fbInvs=", fbInvs)
	uncanceledFbInvs := 0
	for _, fbInv := range fbInvs.Invoices {
		if !fbInv.Canceled() {
			uncanceledFbInvs++
		}
	}
	if uncanceledFbInvs > 1 {
		return 0, fmt.Errorf("duplicate fastbill invoices found")
	}
	alreadyExported := uncanceledFbInvs == 1
	if alreadyExported {
		return 0, ErrInvoiceAlreadyExported
	}

	if inv.InvoiceTitle, err = inv.GetTitle(); err != nil {
		return 0, fmt.Errorf("get title: %v", err)
	}
	inv.DeliveryDate = fmt.Sprintf("%v %v", inv.Month, inv.Year)

	request := Request{
		SERVICE: SERVICE_INVOICE_CREATE,
		DATA:    *inv,
	}
	var response InvoiceCreateResponse

	if err := fb.execGetRequest(&request, &response); err != nil {
		return 0, fmt.Errorf("fb request: %v", err)
	}

	beego.Info("response response:", response.Response)

	if err := response.Error(); err != nil {
		return 0, err
	}

	inv.Id = response.Response.InvoiceId

	return inv.Id, nil
}

type Item struct {
	ArticleNumber int64   `json:"ARTICLE_NUMBER,string,omitempty"`
	Description   string  `json:"DESCRIPTION,"`
	Quantity      float64 `json:"QUANTITY,string,omitempty"`
	UnitPrice     float64 `json:"UNIT_PRICE,string"`
	VatPercent    float64 `json:"VAT_PERCENT,string"`
	IsGross       string  `json:"IS_GROSS,omitempty"`
	SortOrder     string  `json:"SORT_ORDER,"`
}
