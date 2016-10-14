package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego"
)

type PurchasesController struct {
	Controller
}

// @Title Create
// @Description Create space purchase
// @Param	type	query	string	true	"Purchase Type"
// @Success 200 string
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [post]
func (this *PurchasesController) Create() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	purchaseType := this.GetString("type")

	var purchase interface{}
	var err error

	switch purchaseType {
	case purchases.TYPE_TUTOR:
		tp := &purchases.Tutoring{
			Purchase: purchases.Purchase{
				LocationId: locId,
			},
		}
		_, err = purchases.CreateTutoring(tp)
		if err == nil {
			purchase = tp
		}
		break
	default:
		err = fmt.Errorf("unknown purchase type")
	}

	if err != nil {
		beego.Error("Failed to create purchase", err, " (purchaseType=", purchaseType, ")")
		this.CustomAbort(500, "Failed to create purchase")
	}

	this.Data["json"] = purchase
	this.ServeJSON()
}

// @Title GetAll
// @Description Get all purchases
// @Param	type	query	string	true	"Purchase Type"
// @Success 200 string
// @Failure	401 Not authorized
// @Failure	500	Failed to get all machines
// @router / [get]
func (this *PurchasesController) GetAll() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		beego.Error("Not authorized")
		this.Fail("401")
	}

	// This is admin and staff only
	if !this.IsStaffAt(locId) {
		beego.Error("Not authorized")
		this.Fail("401")
	}

	purchaseType := this.GetString("type")

	var ps interface{}
	var err error

	switch purchaseType {
	case purchases.TYPE_TUTOR:
		ps, err = purchases.GetAllTutoringsAt(locId)
		break
	default:
		err = fmt.Errorf("unknown purchase type")
	}

	if err != nil {
		beego.Error("Failed to get purchases", err, " (purchaseType=", purchaseType, ")")
		this.CustomAbort(500, "Failed to get purchases")
	}

	this.Data["json"] = ps
	this.ServeJSON()
}

// @Title Get
// @Description Get purchase by ID
// @Param	id		path 	int	true		"Space Purchase ID"
// @Param	type	query	string	true	"Purchase Type"
// @Success 200 {object}
// @Failure	400	Bad Request
// @Failure	401 Not authorized
// @Failure	500	Internal Server Error
// @router /:id [get]
func (this *PurchasesController) Get() {
	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get :id variable")
		this.CustomAbort(403, "Failed to get space purchase")
	}

	purchaseType := this.GetString("type")

	var locationId int64
	var purchase interface{}

	switch purchaseType {
	case purchases.TYPE_TUTOR:
		var t *purchases.Tutoring
		t, err = purchases.GetTutoring(id)
		locationId = t.LocationId
		purchase = t
		break
	default:
		err = fmt.Errorf("unknown purchase type")
	}

	if err != nil {
		beego.Error("Failed to get purchase", err, " (purchaseType=", purchaseType, ")")
		this.Fail("500")
	}

	if !this.IsAdminAt(locationId) {
		beego.Error("Not authorized")
		this.Fail("401")
	}

	this.Data["json"] = purchase
	this.ServeJSON()
}

// @Title Put
// @Description Update purchase
// @Param	type	query	string	true	"Purchase Type"
// @Success 200 string
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:id [put]
func (this *PurchasesController) Put() {
	var existing *purchases.Purchase
	if this.GetString(":id") != "create" {
		id, err := this.GetInt64(":id")
		if err != nil {
			this.Fail("400")
		}
		existing, err = purchases.Get(id)
		if err != nil {
			beego.Error("Cannot get purchase:", err)
			this.Fail("500")
		}

		if !this.IsAdminAt(existing.LocationId) {
			beego.Error("Unauthorized attempt to update purchase")
			this.Fail("401")
		}

		existingInv, err := invoices.Get(existing.InvoiceId)
		if err != nil {
			beego.Error("get invoice associated to existing purchase:", err)
			this.Fail("500")
		}
		if existingInv.Status != "draft" {
			beego.Error("existing invoice has status", existingInv.Status)
			this.Fail("403")
		}
	}

	assertSameIds := func(newId, newLocationId int64) {
		if existing.Id != newId {
			beego.Error("Id changed")
			this.Fail("403")
		}
		if existing.LocationId != newLocationId {
			beego.Error("Location Id changed")
			this.Fail("403")
		}
	}

	purchaseType := this.GetString("type")

	var response interface{}
	var err error

	switch purchaseType {
	case purchases.TYPE_OTHER:
		p := &purchases.Purchase{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		if err := dec.Decode(p); err != nil {
			beego.Error("Failed to decode json:", err)
			this.Fail("400")
		}

		inv, err := invoices.Get(p.InvoiceId)
		if err != nil {
			beego.Error("getting invoice", inv.Id, ":", err)
			this.Fail("500")
		}

		if inv.Status != "draft" {
			this.Fail(400, "Expected status draft")
		}

		assertSameIds(p.Id, p.LocationId)
		if err = purchases.Update(p); err == nil {
			response = p
		}

		break
	case purchases.TYPE_TUTOR:
		tp := &purchases.Tutoring{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(tp); err != nil {
			beego.Error("Failed to decode json:", err)
			this.Fail("400")
		}

		t := tp.TimeStart

		var inv *invoices.Invoice
		if tp.InvoiceId == 0 {
			inv, err = invoices.GetDraft(tp.LocationId, tp.UserId, t)
			if err != nil {
				beego.Error("getting invoice of", t.Format("01-2006"), ":", err)
				this.Fail("500")
			}
			tp.InvoiceId = inv.Id
		} else {
			inv, err = invoices.Get(tp.InvoiceId)
			if err != nil {
				beego.Error("getting invoice", inv.Id, ":", err)
				this.Fail("500")
			}
		}

		if inv.Status != "draft" {
			beego.Error("the invoice for that month is in status", inv.Status)
			this.Fail("500")
		}

		if this.GetString(":id") == "create" {
			tp.Purchase.Type = purchases.TYPE_TUTOR
			if err = purchases.Create(&tp.Purchase); err == nil {
				response = tp.Purchase
			}
		} else {
			assertSameIds(tp.Id, tp.LocationId)
			if err = tp.Update(); err == nil {
				response = tp
			}
		}

		break
	default:
		err = fmt.Errorf("unknown purchase type")
	}

	if err != nil {
		beego.Error("Failed to save purchase", err, " (purchaseType=", purchaseType, ")")
		this.Fail("500")
	}

	this.Data["json"] = response
	this.ServeJSON()
}

// @Title Archive Purchase
// @Description Archive a purchase
// @Param	purchaseId	query	int	true	"Purchase ID"
// @Success 200 string
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:purchaseId/archive [put]
func (this *PurchasesController) ArchivePurchase() {
	purchaseId, err := this.GetInt64(":purchaseId")
	if err != nil {
		beego.Error("Failed to get :purchaseId variable")
		this.Fail("400")
	}

	purchase, err := purchases.Get(purchaseId)
	if err != nil {
		beego.Error("Failed to get purchase")
		this.Fail("500")
	}

	if !this.IsAdminAt(purchase.LocationId) {
		beego.Error("Unauthorized attempt to archive purchase")
		this.Fail("401")
	}

	err = purchases.Archive(purchase)
	if err != nil {
		beego.Error("Failed to archive purchase")
		this.CustomAbort(500, "Failed to archive purchase")
	}

	this.ServeJSON()
}
