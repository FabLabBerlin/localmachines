package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type SpacePurchasesController struct {
	Controller
}

// @Title Create
// @Description Create space purchase
// @Success 200 {object} models.SpacePurchase
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [post]
func (this *SpacePurchasesController) Create() {
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	spacePurchase := &models.SpacePurchase{}

	_, err := models.CreateSpacePurchase(spacePurchase)
	if err != nil {
		beego.Error("Failed to create space purchase:", err)
		this.CustomAbort(500, "Failed to create space purchase")
	}

	this.Data["json"] = spacePurchase
	this.ServeJson()
}

// @Title GetAll
// @Description Get all space purchases
// @Success 200 {object} models.SpacePurchase
// @Failure	401 Not authorized
// @Failure	500	Failed to get all machines
// @router / [get]
func (this *SpacePurchasesController) GetAll() {

	// This is admin and staff only
	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	spacePurchases, err := models.GetAllSpacePurchases()
	if err != nil {
		beego.Error("Failed to get all space purchases:", err)
		this.CustomAbort(500, "Failed to get all space purchases")
	}

	this.Data["json"] = spacePurchases
	this.ServeJson()
}

// @Title Get
// @Description Get space purchase by ID
// @Param	id		path 	int	true		"Space Purchase ID"
// @Success 200 {object} models.SpacePurchase
// @Failure	400	Bad Request
// @Failure	401 Not authorized
// @Failure	500	Internal Server Error
// @router /:id [get]
func (this *SpacePurchasesController) Get() {
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get :id variable")
		this.CustomAbort(403, "Failed to get space purchase")
	}

	spacePurchase, err := models.GetSpacePurchase(id)
	if err != nil {
		beego.Error("Failed to get space purchase", err)
		this.CustomAbort(500, "Failed to get space purchase")
	}

	this.Data["json"] = spacePurchase
	this.ServeJson()
}

// @Title Put
// @Description Update space purchase
// @Success 200 {object} models.SpacePurchase
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:rid [put]
func (this *SpacePurchasesController) Put() {
	if !this.IsAdmin() {
		beego.Error("Unauthorized attempt to update space purchase")
		this.CustomAbort(401, "Unauthorized")
	}

	spacePurchase := &models.SpacePurchase{}

	dec := json.NewDecoder(this.Ctx.Request.Body)
	defer this.Ctx.Request.Body.Close()
	if err := dec.Decode(spacePurchase); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(400, "Failed to update Space purchase")
	}

	if err := models.UpdateSpacePurchase(spacePurchase); err != nil {
		beego.Error("Failed to update space purchase:", err)
		this.CustomAbort(500, "Failed to update Space purchase")
	}

	beego.Info("spacePurchase=", spacePurchase)

	this.Data["json"] = spacePurchase
	this.ServeJson()
}

// @Title Delete
// @Description Delete space purchase
// @Param	id	path	int	true	"Space Purchase ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure	500	Internal Server Error
// @router /:id [delete]
func (this *SpacePurchasesController) Delete() {

	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get id:", err)
		this.CustomAbort(400, "Failed to delete space purchase")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	err = models.DeleteSpacePurchase(id)
	if err != nil {
		beego.Error("Failed to delete space purchase:", err)
		this.CustomAbort(500, "Failed to delete space purchase")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
