package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models/purchases"
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
	case purchases.PURCHASE_TYPE_CO_WORKING:
		cp := &purchases.CoWorkingPurchase{}
		_, err = purchases.CreateCoWorkingPurchase(cp)
		if err == nil {
			purchase = cp
		}
		break
	case purchases.PURCHASE_TYPE_SPACE_PURCHASE:
		spacePurchase := &purchases.SpacePurchase{}
		_, err = purchases.CreateSpacePurchase(spacePurchase)
		if err == nil {
			purchase = spacePurchase
		}
		break
	case purchases.PURCHASE_TYPE_TUTOR:
		tutoringPurchase := &purchases.TutoringPurchase{}
		_, err = purchases.CreateTutoringPurchase(tutoringPurchase)
		if err == nil {
			purchase = tutoringPurchase
		}
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
	case purchases.PURCHASE_TYPE_CO_WORKING:
		ps, err = purchases.GetAllCoWorkingPurchases()
		break
	case purchases.PURCHASE_TYPE_SPACE_PURCHASE:
		ps, err = purchases.GetAllSpacePurchases()
		break
	case purchases.PURCHASE_TYPE_TUTOR:
		ps, err = purchases.GetAllTutoringPurchases()
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
// @Success 200 {object} purchases.SpacePurchase
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
	case purchases.PURCHASE_TYPE_CO_WORKING:
		purchase, err = purchases.GetCoWorkingPurchase(id)
		break
	case purchases.PURCHASE_TYPE_SPACE_PURCHASE:
		purchase, err = purchases.GetSpacePurchase(id)
		break
	case purchases.PURCHASE_TYPE_TUTOR:
		purchase, err = purchases.GetTutoringPurchase(id)
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
// @Description Update space purchase
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
	case purchases.PURCHASE_TYPE_CO_WORKING:
		cp := &purchases.CoWorkingPurchase{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(cp); err != nil {
			beego.Error("Failed to decode json:", err)
			this.CustomAbort(400, "Failed to update Co-Working purchase")
		}

		if err = purchases.UpdateCoWorkingPurchase(cp); err == nil {
			response = cp
		}
		break
	case purchases.PURCHASE_TYPE_SPACE_PURCHASE:
		spacePurchase := &purchases.SpacePurchase{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(spacePurchase); err != nil {
			beego.Error("Failed to decode json:", err)
			this.CustomAbort(400, "Failed to update Space purchase")
		}

		if err = purchases.UpdateSpacePurchase(spacePurchase); err == nil {
			response = spacePurchase
		}
		break
	case purchases.PURCHASE_TYPE_TUTOR:
		tutoringPurchase := &purchases.TutoringPurchase{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(tutoringPurchase); err != nil {
			beego.Error("Failed to decode json:", err)
			this.CustomAbort(400, "Failed to update Tutoring purchase")
		}
		beego.Info("tp: time end planned:", tutoringPurchase.TimeEndPlanned)
		if err = purchases.UpdateTutoringPurchase(tutoringPurchase); err == nil {
			response = tutoringPurchase
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

	var purchase *purchases.Purchase
	purchase, err = purchases.GetPurchase(purchaseId)
	if err != nil {
		beego.Error("Failed to get purchase")
		this.CustomAbort(500, "Failed to get purchase")
	}

	err = purchases.ArchivePurchase(purchase)
	if err != nil {
		beego.Error("Failed to archive purchase")
		this.CustomAbort(500, "Failed to archive purchase")
	}

	this.ServeJson()
}
