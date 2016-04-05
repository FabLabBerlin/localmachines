package controllers

import (
	"encoding/json"
	"fmt"
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
	case purchases.TYPE_CO_WORKING:
		cp := &purchases.CoWorking{
			Purchase: purchases.Purchase{
				LocationId: locId,
			},
		}
		_, err = purchases.CreateCoWorking(cp)
		if err == nil {
			purchase = cp
		}
		break
	case purchases.TYPE_SPACE:
		sp := purchases.NewSpace(locId)
		if err = sp.Save(); err == nil {
			purchase = sp
		}
		break
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
		this.Abort("401")
	}

	// This is admin and staff only
	if !this.IsStaffAt(locId) {
		beego.Error("Not authorized")
		this.Abort("401")
	}

	purchaseType := this.GetString("type")

	var ps interface{}
	var err error

	switch purchaseType {
	case purchases.TYPE_CO_WORKING:
		ps, err = purchases.GetAllCoWorkingAt(locId)
		break
	case purchases.TYPE_SPACE:
		ps, err = purchases.GetAllSpaceAt(locId)
		break
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
	case purchases.TYPE_CO_WORKING:
		var cw *purchases.CoWorking
		cw, err = purchases.GetCoWorking(id)
		locationId = cw.LocationId
		purchase = cw
		break
	case purchases.TYPE_SPACE:
		var s *purchases.Space
		s, err = purchases.GetSpace(id)
		locationId = s.LocationId
		purchase = s
		break
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
		this.Abort("500")
	}

	if !this.IsAdminAt(locationId) {
		beego.Error("Not authorized")
		this.Abort("401")
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
// @router /:rid [put]
func (this *PurchasesController) Put() {
	id, err := this.GetInt64(":id")
	if err != nil {
		this.Abort("400")
	}
	existing, err := purchases.Get(id)
	if err != nil {
		beego.Error("Cannot get purchase:", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(existing.LocationId) {
		beego.Error("Unauthorized attempt to update purchase")
		this.Abort("401")
	}

	assertSameIds := func(newId, newLocationId int64) {
		if existing.Id != newId {
			beego.Error("Id changed")
			this.Abort("403")
		}
		if existing.LocationId != newLocationId {
			beego.Error("Location Id changed")
			this.Abort("403")
		}
	}

	purchaseType := this.GetString("type")

	var response interface{}

	switch purchaseType {
	case purchases.TYPE_CO_WORKING:
		cp := &purchases.CoWorking{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(cp); err != nil {
			beego.Error("Failed to decode json:", err)
			this.CustomAbort(400, "Failed to update Co-Working purchase")
		}

		assertSameIds(cp.Id, cp.LocationId)

		if err = cp.Update(); err == nil {
			response = cp
		}
		break
	case purchases.TYPE_SPACE:
		sp := &purchases.Space{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(sp); err != nil {
			beego.Error("Failed to decode json:", err)
			this.Abort("400")
		}

		assertSameIds(sp.Id, sp.LocationId)

		if err = sp.Update(); err == nil {
			response = sp
		}
		break
	case purchases.TYPE_TUTOR:
		tp := &purchases.Tutoring{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(tp); err != nil {
			beego.Error("Failed to decode json:", err)
			this.Abort("400")
		}

		assertSameIds(tp.Id, tp.LocationId)

		if err = tp.Update(); err == nil {
			response = tp
		}
		break
	default:
		err = fmt.Errorf("unknown purchase type")
	}

	if err != nil {
		beego.Error("Failed to update purchase", err, " (purchaseType=", purchaseType, ")")
		this.Abort("500")
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
		this.Abort("400")
	}

	purchase, err := purchases.Get(purchaseId)
	if err != nil {
		beego.Error("Failed to get purchase")
		this.Abort("500")
	}

	if !this.IsAdminAt(purchase.LocationId) {
		beego.Error("Unauthorized attempt to archive purchase")
		this.Abort("401")
	}

	err = purchases.Archive(purchase)
	if err != nil {
		beego.Error("Failed to archive purchase")
		this.CustomAbort(500, "Failed to archive purchase")
	}

	this.ServeJSON()
}
