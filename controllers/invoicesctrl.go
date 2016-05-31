package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"io"
	"os"
	"time"
)

type InvoicesController struct {
	Controller
}

// @Title Get All monthly earnings
// @Description Get all monthly earnings from the database
// @Success 200 []*monthly_earning.MonthlyEarning
// @Failure	403	Failed to get all monthly earnings
// @Failure	401	Not authorized
// @router / [get]
func (this *InvoicesController) GetAll() {

	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	mes, err := monthly_earning.GetAllAt(locId)
	if err != nil {
		beego.Error("Failed to get all monthly earnings")
		this.Abort("500")
	}

	this.Data["json"] = mes
	this.ServeJSON()
}

// @Title Create monthly earning
// @Description Create monthly earning from date range
// @Param	startDate	query 	string	true	"Period start date"
// @Param	endDate		query 	string	true	"Period end date"
// @Success 200 {object}
// @Failure	403	Failed to create monthly earning
// @Failure	401	Not authorized
// @router / [post]
func (this *InvoicesController) Create() {

	// Only local admin can use this API call
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	interval, err := this.parseParams()
	if err != nil {
		beego.Error("Request parameters:", err)
		this.CustomAbort(400, "Bad request")
	}

	mes, err := monthly_earning.Create(locId, interval)
	if err != nil {
		beego.Error("Failed to create monthly earnings:", err)
		this.CustomAbort(403, "Failed to create monthly earnings")
	}

	this.Data["json"] = mes
	this.ServeJSON()
}

// @Title Create Fastbill drafts
// @Description Create fastbill drafts from date range
// @Param	startDate	query 	string	true	"Period start date"
// @Param	endDate		query 	string	true	"Period end date"
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	403	Failed to create Fastbill drafts
// @router /:iid/create_drafts [post]
func (this *InvoicesController) CreateDrafts() {

	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	if locId != 1 {
		this.CustomAbort(400, "Fastbill not activated for this location")
	}

	if !this.IsSuperAdmin() {
		beego.Error("User must be super admin")
		this.CustomAbort(401, "Not authorized")
	}

	iid, err := this.GetInt64(":iid")
	if err != nil {
		beego.Error("Failed to get iid:", err)
		this.CustomAbort(400, "Bad request")
	}

	dbMes, err := monthly_earning.Get(iid)
	if err != nil {
		beego.Error("invoices get:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	me, err := monthly_earning.New(locId, dbMes.Interval())
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	creationReport := monthly_earning.CreateFastbillDrafts(me)
	beego.Info("created invoice drafts with IDs", creationReport.Ids)

	this.Data["json"] = creationReport
	this.ServeJSON()
}

func (this *InvoicesController) parseParams() (interval lib.Interval, err error) {
	startDate := this.GetString("startDate")
	if startDate == "" {
		err = fmt.Errorf("Missing start date")
		return
	}

	endDate := this.GetString("endDate")
	if endDate == "" {
		err = fmt.Errorf("Missing end date")
		return
	}

	interval, err = lib.NewInterval(startDate, endDate)
	if err != nil {
		err = fmt.Errorf("parse interval: %v", err)
		return
	}

	return
}

// @Title DownloadExcelExport
// @Description Download existing excel export
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:id/download_excel [get]
func (this *InvoicesController) DownloadExcelExport() {

	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get id:", err)
		this.CustomAbort(400, "Bad request")
	}

	me, err := monthly_earning.Get(id)
	if err != nil {
		beego.Error("invoices get:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdminAt(me.LocationId) {
		this.CustomAbort(401, "Not authorized")
	}

	f, err := os.Open(me.FilePath)
	if err != nil {
		beego.Error("open", me.FilePath, ":", err)
		if os.IsNotExist(err) {
			this.CustomAbort(404, "Cannot find Excel Export file")
		} else {
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	defer f.Close()
	this.Ctx.Output.ContentType("xlsx")
	if _, err := io.Copy(this.Ctx.ResponseWriter, f); err != nil {
		beego.Error("copy ", me.FilePath, ":", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Finish()
}

type MonthlySummary struct {
	User          users.User
	InvoiceNumber string
	InvoiceStatus string
	Amount        float64
	CustomerId    int64
}

// @Title GetMonth
// @Description Get overview for a month
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month [get]
func (this *InvoicesController) GetMonth() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	interval := lib.Interval{
		MonthFrom: int(month),
		YearFrom:  int(year),
		MonthTo:   int(month),
		YearTo:    int(year),
	}

	me, err := monthly_earning.New(locId, interval)
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	sumsByCustomerId := make(map[int64]MonthlySummary)
	for _, inv := range me.Invoices {
		if err := inv.CalculateTotals(); err != nil {
			beego.Error("CalculateTotals:", err)
			this.Abort("500")
		}
		sum := MonthlySummary{
			User:   *inv.User,
			Amount: inv.Sums.All.PriceInclVAT,
		}
		if sum.User.ClientId != 0 {
			customerId, err := fastbill.GetCustomerId(sum.User)
			if err == nil {
				sumsByCustomerId[customerId] = sum
			} else {
				beego.Error("Failed to get customer id for customer #", sum.User.ClientId)
			}
		} else {
			beego.Error("ClientId=0 for User", sum.User.FirstName, sum.User.LastName)
		}
	}

	l, err := fastbill.ListInvoices(year, time.Month(month))
	if err != nil {
		beego.Error("Failed to get invoice list from fastbill:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	result := make([]MonthlySummary, 0, 2*len(sumsByCustomerId))
	beego.Info("Processing", len(l), "Fastbill invoices")
	for _, inv := range l {
		sum, ok := sumsByCustomerId[inv.CustomerId]
		if ok {
			tmp := sum
			beego.Info("Processing sum for CustomerId", inv.CustomerId)
			tmp.InvoiceNumber = inv.InvoiceNumber
			tmp.CustomerId = inv.CustomerId
			tmp.InvoiceStatus = inv.Type
			result = append(result, tmp)
		} else {
			beego.Info("no sum for CustomerId", inv.CustomerId)
		}
	}
	buf, _ := json.Marshal(l)
	fmt.Printf("l = %v\n", string(buf))

	this.Data["json"] = result
	this.ServeJSON()
}

// @Title GetInvoices
// @Description Get invoices for a month including drafts and canceled ones
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/invoices [get]
func (this *InvoicesController) GetInvoices() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	invs, err := invoices.GetIdsAndStatuses(locId, year, month)
	if err != nil {
		beego.Error("Failed to get invoices ids and statuses:", err)
		this.Abort("500")
	}

	this.Data["json"] = invs
	this.ServeJSON()
}

// @Title GetUser
// @Description Get monthly overview for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid [get]
func (this *InvoicesController) GetUser() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt64(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.CustomAbort(400, "Bad request")
	}

	interval := lib.Interval{
		MonthFrom: int(month),
		YearFrom:  int(year),
		MonthTo:   int(month),
		YearTo:    int(year),
	}

	me, err := monthly_earning.New(locId, interval)
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	sums := make([]MonthlySummary, 0, len(me.Invoices))
	for _, inv := range me.Invoices {
		sum := MonthlySummary{
			User: *inv.User,
		}
		for _, p := range inv.Purchases {
			sum.Amount += p.DiscountedTotal
		}
		sums = append(sums, sum)
	}

	var userInv *invoices.Invoice

	for _, inv := range me.Invoices {
		if inv.User.Id == uid {
			userInv = inv
		}
	}

	if err := userInv.CalculateTotals(); err != nil {
		beego.Error("CalculateTotals:", err)
		this.Abort("500")
	}

	this.Data["json"] = userInv
	this.ServeJSON()
}

// @Title GetStatuses
// @Description Get statuses for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid/statuses [get]
func (this *InvoicesController) GetStatuses() {
	_, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt64(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.CustomAbort(400, "Bad request")
	}

	user, err := users.GetUser(uid)
	if err != nil {
		beego.Error("Failed to get user:", err)
		this.Abort("500")
	}

	inv := fastbill.Invoice{
		Month:          time.Month(month).String(),
		Year:           int(year),
		CustomerNumber: user.ClientId,
	}

	var existingMonth fastbill.ExistingMonth

	key := fmt.Sprintf("/months/%v/%v/users/%v/statuses", year, month, user.Id)
	redis.Cached(key, 3600, &existingMonth, func() interface{} {
		ivs, err := inv.FetchExisting()
		if err != nil {
			beego.Error("Failed to fetch existing fastbill invoice:", err)
			this.Abort("500")
		}
		return *ivs
	})

	this.Data["json"] = existingMonth
	this.ServeJSON()
}

// @Title Create draft
// @Description Create draft for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid/draft [post]
func (this *InvoicesController) CreateDraft() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt64(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.CustomAbort(400, "Bad request")
	}

	interval := lib.Interval{
		MonthFrom: int(month),
		YearFrom:  int(year),
		MonthTo:   int(month),
		YearTo:    int(year),
	}

	me, err := monthly_earning.New(locId, interval)
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	sums := make([]MonthlySummary, 0, len(me.Invoices))
	for _, inv := range me.Invoices {
		sum := MonthlySummary{
			User: *inv.User,
		}
		for _, p := range inv.Purchases {
			sum.Amount += p.DiscountedTotal
		}
		sums = append(sums, sum)
	}

	var userInv *invoices.Invoice

	for _, inv := range me.Invoices {
		if inv.User.Id == uid {
			userInv = inv
		}
	}

	if err := userInv.CalculateTotals(); err != nil {
		beego.Error("CalculateTotals:", err)
		this.Abort("500")
	}

	fbDraft, empty, err := monthly_earning.CreateFastbillDraft(me, userInv)
	if err != nil {
		beego.Error("Create fastbill draft:", err)
		this.Abort("500")
	}
	beego.Info("empty=", empty)
	beego.Info("fbDraft=", fbDraft)

	this.Data["json"] = fbDraft
	this.ServeJSON()
}

// @Title Update user invoicing data
// @Description Update invoicing data for a user
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/users/:uid/update [post]
func (this *InvoicesController) Update() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt64(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get uid:", err)
		this.CustomAbort(400, "Bad request")
	}

	interval := lib.Interval{
		MonthFrom: int(month),
		YearFrom:  int(year),
		MonthTo:   int(month),
		YearTo:    int(year),
	}

	me, err := monthly_earning.New(locId, interval)
	if err != nil {
		beego.Error("Failed to make new invoices:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	sums := make([]MonthlySummary, 0, len(me.Invoices))
	for _, inv := range me.Invoices {
		sum := MonthlySummary{
			User: *inv.User,
		}
		for _, p := range inv.Purchases {
			sum.Amount += p.DiscountedTotal
		}
		sums = append(sums, sum)
	}

	var userInv *invoices.Invoice

	for _, inv := range me.Invoices {
		if inv.User.Id == uid {
			userInv = inv
		}
	}

	if err := userInv.CalculateTotals(); err != nil {
		beego.Error("CalculateTotals:", err)
		this.Abort("500")
	}

	beego.Error("Not implemented")
	this.CustomAbort(500, "Not implemented")
}

// @Title SyncFastbillInvoices
// @Description SyncFastbillInvoices
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /months/:year/:month/sync [get]
func (this *InvoicesController) SyncFastbillInvoices() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	year, err := this.GetInt(":year")
	if err != nil {
		beego.Error("Failed to get year:", err)
		this.CustomAbort(400, "Bad request")
	}

	month, err := this.GetInt64(":month")
	if err != nil {
		beego.Error("Failed to get month:", err)
		this.CustomAbort(400, "Bad request")
	}

	l, err := fastbill.ListInvoices(year, time.Month(month))
	if err != nil {
		beego.Error("Failed to get invoice list from fastbill:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	usrs, err := users.GetAllUsersAt(locId)
	if err != nil {
		beego.Error("Failed to get user list:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	for _, fbInv := range l {
		inv := invoices.Invoice{
			LocationId: locId,
			FastbillId: fbInv.Id,
			FastbillNo: fbInv.InvoiceNumber,
			CustomerId: fbInv.CustomerId,
			Status:     fbInv.Type,
		}
		inv.Month, inv.Year, inv.CustomerNo, err = fbInv.ParseTitle()
		if err != nil {
			beego.Error("Cannot parse", fbInv.InvoiceTitle)
			continue
		}
		for _, u := range usrs {
			if u.ClientId == inv.CustomerNo {
				inv.UserId = u.Id
				break
			}
		}
		if inv.UserId == 0 {
			beego.Error("Cannot find user for customer number", inv.CustomerNo)
			continue
		}
		beego.Info("Adding invoice of user", inv.UserId)
		if _, err := invoices.CreateOrUpdate(&inv); err != nil {
			beego.Error("Failed to create or update inv:", err)
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	this.ServeJSON()
}
