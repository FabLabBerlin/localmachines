package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type PurchasesController struct {
	Controller
}

// @Title Create
// @Description Create space purchase
// @Param	type	query	string	true	"Purchase Type"
// @Success 200 {object} interface{}
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
	case models.PURCHASE_TYPE_SPACE_PURCHASE:
		spacePurchase := &models.SpacePurchase{}
		_, err = models.CreateSpacePurchase(spacePurchase)
		if err == nil {
			purchase = spacePurchase
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
// @Success 200 {object} interface{}
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

	var purchases interface{}
	var err error

	switch purchaseType {
	case models.PURCHASE_TYPE_SPACE_PURCHASE:
		purchases, err = models.GetAllSpacePurchases()
		break
	default:
		err = fmt.Errorf("unknown purchase type")
	}

	if err != nil {
		beego.Error("Failed to get purchases", err, " (purchaseType=", purchaseType, ")")
		this.CustomAbort(500, "Failed to get purchases")
	}

	this.Data["json"] = purchases
	this.ServeJson()
}

// @Title Get
// @Description Get purchase by ID
// @Param	id		path 	int	true		"Space Purchase ID"
// @Param	type	query	string	true	"Purchase Type"
// @Success 200 {object} models.SpacePurchase
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
	case models.PURCHASE_TYPE_SPACE_PURCHASE:
		purchase, err = models.GetSpacePurchase(id)
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
// @Success 200 {object} interface{}
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
	case models.PURCHASE_TYPE_SPACE_PURCHASE:
		spacePurchase := &models.SpacePurchase{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(spacePurchase); err != nil {
			beego.Error("Failed to decode json:", err)
			this.CustomAbort(400, "Failed to update Space purchase")
		}

		if err = models.UpdateSpacePurchase(spacePurchase); err == nil {
			response = spacePurchase
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
