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
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	purchaseType := this.GetString("type")

	var purchase interface{}
	var err error

	switch purchaseType {
	case purchases.TYPE_CO_WORKING:
		cp := &purchases.CoWorking{}
		_, err = purchases.CreateCoWorking(cp)
		if err == nil {
			purchase = cp
		}
		break
	case purchases.TYPE_SPACE:
		sp := purchases.NewSpace()
		if err = sp.Save(); err == nil {
			purchase = sp
		}
		break
	case purchases.TYPE_TUTOR:
		tp := &purchases.Tutoring{}
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
	this.ServeJson()
}

// @Title GetAll
// @Description Get all purchases
// @Param	type	query	string	true	"Purchase Type"
// @Success 200 string
// @Failure	401 Not authorized
// @Failure	500	Failed to get all machines
// @router / [get]
func (this *PurchasesController) GetAll() {

	// This is admin and staff only
	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	purchaseType := this.GetString("type")

	var ps interface{}
	var err error

	switch purchaseType {
	case purchases.TYPE_CO_WORKING:
		ps, err = purchases.GetAllCoWorking()
		break
	case purchases.TYPE_SPACE:
		ps, err = purchases.GetAllSpace()
		break
	case purchases.TYPE_TUTOR:
		ps, err = purchases.GetAllTutorings()
		break
	default:
		err = fmt.Errorf("unknown purchase type")
	}

	if err != nil {
		beego.Error("Failed to get purchases", err, " (purchaseType=", purchaseType, ")")
		this.CustomAbort(500, "Failed to get purchases")
	}

	this.Data["json"] = ps
	this.ServeJson()
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
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get :id variable")
		this.CustomAbort(403, "Failed to get space purchase")
	}

	purchaseType := this.GetString("type")

	var purchase interface{}

	switch purchaseType {
	case purchases.TYPE_CO_WORKING:
		purchase, err = purchases.GetCoWorking(id)
		break
	case purchases.TYPE_SPACE:
		purchase, err = purchases.GetSpace(id)
		break
	case purchases.TYPE_TUTOR:
		purchase, err = purchases.GetTutoring(id)
		break
	default:
		err = fmt.Errorf("unknown purchase type")
	}

	if err != nil {
		beego.Error("Failed to get purchase", err, " (purchaseType=", purchaseType, ")")
		this.CustomAbort(500, "Failed to get purchase")
	}

	this.Data["json"] = purchase
	this.ServeJson()
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
	if !this.IsAdmin() {
		beego.Error("Unauthorized attempt to update space purchase")
		this.CustomAbort(401, "Unauthorized")
	}

	purchaseType := this.GetString("type")

	var response interface{}
	var err error

	switch purchaseType {
	case purchases.TYPE_CO_WORKING:
		cp := &purchases.CoWorking{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(cp); err != nil {
			beego.Error("Failed to decode json:", err)
			this.CustomAbort(400, "Failed to update Co-Working purchase")
		}

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
			this.CustomAbort(400, "Failed to update Space purchase")
		}

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
			this.CustomAbort(400, "Failed to update Tutoring purchase")
		}
		beego.Info("tp: time end planned:", tp.TimeEndPlanned)
		if err = tp.Update(); err == nil {
			response = tp
		}
		break
	default:
		err = fmt.Errorf("unknown purchase type")
	}

	if err != nil {
		beego.Error("Failed to update purchase", err, " (purchaseType=", purchaseType, ")")
		this.CustomAbort(500, "Failed to update purchase")
	}

	this.Data["json"] = response
	this.ServeJson()
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
	if !this.IsAdmin() {
		beego.Error("Unauthorized attempt to archive purchase")
		this.CustomAbort(401, "Unauthorized")
	}

	purchaseId, err := this.GetInt64(":purchaseId")
	if err != nil {
		beego.Error("Failed to get :purchaseId variable")
		this.CustomAbort(400, "Incorrect purchaseId")
	}

	purchase, err := purchases.Get(purchaseId)
	if err != nil {
		beego.Error("Failed to get purchase")
		this.CustomAbort(500, "Failed to get purchase")
	}

	err = purchases.Archive(purchase)
	if err != nil {
		beego.Error("Failed to archive purchase")
		this.CustomAbort(500, "Failed to archive purchase")
	}

	this.ServeJson()
}
