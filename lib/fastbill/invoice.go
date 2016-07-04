package fastbill

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
)

const (
	FB_DATE_LONG    = "2006-01-02 15:04:05"
	FB_DATE_SHORT   = "2006-01-02"
	FB_DATE_PAID    = FB_DATE_LONG
	FB_DATE_INVOICE = FB_DATE_SHORT
	FB_DATE_DUE     = FB_DATE_LONG
)

const (
	TemplateStandardId        = 1
	TemplateMakeaIndustriesId = 889700
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
	Id                int64   `json:"INVOICE_ID,string"`
	Type              string  `json:"TYPE,omitempty"`
	InvoiceDateString string  `json:"INVOICE_DATE,omitempty"`
	PaidDateString    string  `json:"PAID_DATE,omitempty"`
	DueDateString     string  `json:"DUE_DATE,omitempty"`
	IsCanceled        string  `json:"IS_CANCELED,omitempty"`
	Total             float64 `json:"TOTAL,omitempty"`
	CustomerId        int64   `json:"CUSTOMER_ID,string,omitempty"`
	InvoiceNumber     string  `json:"INVOICE_NUMBER,omitempty"`
	InvoiceTitle      string  `json:"INVOICE_TITLE,omitempty"`
	VatPercent        float64 `json:"VAT_PERCENT,omitempty,string"`
}

func (this *InvoiceGetResponseInvoice) Canceled() bool {
	return this.IsCanceled == "1"
}

func (this *InvoiceGetResponseInvoice) DueDate() time.Time {
	return this.parseDate(FB_DATE_DUE, this.DueDateString)
}

func (this *InvoiceGetResponseInvoice) InvoiceDate() time.Time {
	return this.parseDate(FB_DATE_INVOICE, this.InvoiceDateString)
}

func (this *InvoiceGetResponseInvoice) PaidDate() time.Time {
	return this.parseDate(FB_DATE_PAID, this.PaidDateString)
}

func (this *InvoiceGetResponseInvoice) parseDate(layout, s string) time.Time {
	t, err := time.Parse(layout, s)
	if err == nil {
		return t
	} else {
		return time.Time{}
	}
}

// ParseTitle like "March 2016 Invoice for Customer Number 696"
func (this *InvoiceGetResponseInvoice) ParseTitle() (month, year int, customerNo int64, err error) {
	parts := strings.Split(this.InvoiceTitle, " ")

	if len(parts) != 7 {
		return 0, 0, 0, fmt.Errorf("cannot parse '%v'", this.InvoiceTitle)
	}

	for i := 1; i <= 12; i++ {
		if time.Month(i).String() == parts[0] {
			month = i
			break
		}
	}
	if month == 0 {
		return 0, 0, 0, fmt.Errorf("cannot parse month '%v'", parts[0])
	}

	if year, err = strconv.Atoi(parts[1]); err != nil {
		return 0, 0, 0, fmt.Errorf("cannot parse year '%v'", parts[1])
	}

	if customerNo, err = strconv.ParseInt(parts[6], 10, 64); err != nil {
		return 0, 0, 0, fmt.Errorf("cannot parse customer no '%v'", parts[6])
	}

	return
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
	CustomerId   int64  `json:"CUSTOMER_ID,string,omitempty"`
	Month        int    `json:"MONTH,string,omitempty"`
	Year         int    `json:"YEAR,string,omitempty"`
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

func (inv *Invoice) Submit(overwriteExisting bool) (id int64, err error) {
	fb := New()
	fbInvs, err := inv.FetchExisting()
	if err != nil {
		return 0, fmt.Errorf("checking if already exported: %v", err)
	}
	beego.Info("lib.fastbill:Invoice#Submit fbInvs=", fbInvs)
	uncanceledFbInvIds := make([]int64, 0, 1)
	for _, fbInv := range fbInvs.Invoices {
		if !fbInv.Canceled() {
			uncanceledFbInvIds = append(uncanceledFbInvIds, fbInv.Id)
		}
	}
	if len(uncanceledFbInvIds) > 1 {
		return 0, fmt.Errorf("duplicate fastbill invoices found")
	}
	alreadyExported := len(uncanceledFbInvIds) == 1
	if alreadyExported {
		if overwriteExisting {
			if err := deleteInvoice(uncanceledFbInvIds[0]); err != nil {
				return 0, fmt.Errorf("error while deleting old draft: %v", err)
			}
		} else {
			return 0, ErrInvoiceAlreadyExported
		}
	}

	if inv.InvoiceTitle, err = inv.GetTitle(); err != nil {
		return 0, fmt.Errorf("get title: %v", err)
	}
	inv.DeliveryDate = fmt.Sprintf("%v %v", inv.Month, inv.Year)

	for _, item := range inv.Items {
		if item.VatPercent < 0.01 {
			return 0, fmt.Errorf("VAT seems to be zero for an item")
		}
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

	if err := response.Error(); err != nil {
		return 0, err
	}

	inv.Id = response.Response.InvoiceId

	return inv.Id, nil
}

type InvoiceDeleteResponse struct {
	Request struct {
		Id int64 `json:"INVOICE_ID,string,omitempty"`
	} `json:"REQUEST,omitempty"`
	Response struct {
		Status string   `json:"STATUS,omitempty"`
		Errors []string `json:"ERRORS,omitempty"`
	}
}

func (this *InvoiceDeleteResponse) Error() error {
	if len(this.Response.Errors) == 0 {
		return nil
	} else {
		return errors.New(strings.Join(this.Response.Errors, "; "))
	}
}

func deleteInvoice(id int64) (err error) {
	fb := New()

	if id <= 0 {
		return fmt.Errorf("invalid id: %v", err)
	}

	request := Request{
		SERVICE: SERVICE_INVOICE_DELETE,
		DATA: map[string]string{
			"INVOICE_ID": strconv.FormatInt(id, 10),
		},
	}

	var response InvoiceDeleteResponse

	if err := fb.execGetRequest(&request, &response); err != nil {
		return fmt.Errorf("fb request: %v", err)
	}

	beego.Info("response response:", response.Response)

	if err := response.Error(); err != nil {
		return err
	}

	return nil
}

type InvoiceCompleteResponse struct {
	Request struct {
		Id int64 `json:"INVOICE_ID,string,omitempty"`
	} `json:"REQUEST,omitempty"`
	Response struct {
		Status        string   `json:"STATUS,omitempty"`
		InvoiceNumber string   `json:"INVOICE_NUMBER,omitempty"`
		Errors        []string `json:"ERRORS,omitempty"`
	}
}

func (this *InvoiceCompleteResponse) Error() error {
	if len(this.Response.Errors) == 0 {
		return nil
	} else {
		return errors.New(strings.Join(this.Response.Errors, "; "))
	}
}

func CompleteInvoice(id int64) (fastbillNo string, err error) {
	fb := New()

	if id <= 0 {
		return "", fmt.Errorf("invalid id: %v", err)
	}

	request := Request{
		SERVICE: SERVICE_INVOICE_COMPLETE,
		DATA: map[string]string{
			"INVOICE_ID": strconv.FormatInt(id, 10),
		},
	}

	var response InvoiceCompleteResponse

	if err := fb.execGetRequest(&request, &response); err != nil {
		return "", fmt.Errorf("fb request: %v", err)
	}

	beego.Info("response response:", response.Response)

	if err := response.Error(); err != nil {
		return "", err
	}

	return response.Response.InvoiceNumber, nil
}

type InvoiceGenericStatusResponse struct {
	Request struct {
		Id int64 `json:"INVOICE_ID,string,omitempty"`
	} `json:"REQUEST,omitempty"`
	Response struct {
		Status string   `json:"STATUS,omitempty"`
		Errors []string `json:"ERRORS,omitempty"`
	}
}

func (this *InvoiceGenericStatusResponse) Error() error {
	if len(this.Response.Errors) == 0 {
		return nil
	} else {
		return errors.New(strings.Join(this.Response.Errors, "; "))
	}
}

func CancelInvoice(id int64) (err error) {
	fb := New()

	if id <= 0 {
		return fmt.Errorf("invalid id: %v", err)
	}

	request := Request{
		SERVICE: SERVICE_INVOICE_CANCEL,
		DATA: map[string]string{
			"INVOICE_ID": strconv.FormatInt(id, 10),
		},
	}

	var response InvoiceGenericStatusResponse

	if err := fb.execGetRequest(&request, &response); err != nil {
		return fmt.Errorf("fb request: %v", err)
	}

	beego.Info("response response:", response.Response)

	if err := response.Error(); err != nil {
		return err
	}

	return nil
}

func SendInvoiceByEmail(id int64, user *users.User) (err error) {
	fb := New()

	if id <= 0 {
		return fmt.Errorf("invalid id: %v", id)
	}

	beego.Info("user.Email=", user.Email)

	request := Request{
		SERVICE: SERVICE_INVOICE_SEND_BY_EMAIL,
		DATA: map[string]interface{}{
			"INVOICE_ID": strconv.FormatInt(id, 10),
			"RECIPIENT": map[string]string{
				"TO": user.Email,
			},
			"RECEIPT_CONFIRMATION": "0",
		},
	}

	var response InvoiceGenericStatusResponse

	if err := fb.execGetRequest(&request, &response); err != nil {
		return fmt.Errorf("fb request: %v", err)
	}

	beego.Info("response response:", response.Response)

	if err := response.Error(); err != nil {
		return err
	}

	return nil
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

func ListInvoices(customerId int64) ([]InvoiceGetResponseInvoice, error) {
	all := make([]InvoiceGetResponseInvoice, 0, 100)
	limit := 100
	for offset := 0; ; offset += limit {
		l, err := listInvoices(customerId, offset, limit)
		if err != nil {
			return nil, fmt.Errorf("@offset=%v: err", offset, err)
		}
		for _, inv := range l {
			all = append(all, inv)
		}
		if len(l) < limit {
			break
		}
	}
	return all, nil
}

func listInvoices(customerId int64, offset, limit int) ([]InvoiceGetResponseInvoice, error) {
	var err error
	fb := New()
	filter := InvoiceFilter{
		CustomerId: customerId,
	}
	request := Request{
		SERVICE: SERVICE_INVOICE_GET,
		FILTER:  filter,
		LIMIT:   int64(limit),
		OFFSET:  int64(offset),
	}
	var response InvoiceGetResponse
	if err = fb.execGetRequest(&request, &response); err != nil {
		return nil, fmt.Errorf("get request: %v", err)
	}
	return response.Response.Invoices, err
}
