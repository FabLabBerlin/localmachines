package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type SpacesController struct {
	Controller
}

// @Title Create
// @Description Create space
// @Param	name	query	string	true	"Space Name"
// @Success 200 {object} models.Space
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [post]
func (this *SpacesController) Create() {
	name := this.GetString("name")

	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized to create space")
		this.CustomAbort(401, "Not authorized")
	}

	space, err := models.CreateSpace(name)
	if err != nil {
		beego.Error("Failed to create space", err)
		this.CustomAbort(500, "Failed to create space")
	}

	this.Data["json"] = space
	this.ServeJson()
}

// @Title Put
// @Description Update space
// @Success 201 {object} models.Space
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:rid [put]
func (this *SpacesController) Put() {
	if !this.IsAdmin() {
		beego.Error("Unauthorized attempt to update user")
		this.CustomAbort(401, "Unauthorized")
	}

	space := &models.Space{}

	dec := json.NewDecoder(this.Ctx.Request.Body)
	defer this.Ctx.Request.Body.Close()
	if err := dec.Decode(space); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(400, "Failed to update space")
	}

	if err := models.UpdateSpace(space); err != nil {
		beego.Error("Failed to update space:", err)
		this.CustomAbort(500, "Failed to update space")
	}

	this.Data["json"] = space
	this.ServeJson()
}

// @Title Get
// @Description Get space by ID
// @Param	id		path 	int	true		"Space ID"
// @Success 200 {object} models.Space
// @Failure	400	Bad Request
// @Failure	500	Failed to get space
// @router /:id [get]
func (this *SpacesController) Get() {
	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get read id")
		this.CustomAbort(400, "Failed to get space")
	}

	space, err := models.GetSpace(id)
	if err != nil {
		beego.Error("Failed to get space", err)
		this.CustomAbort(500, "Failed to get space")
	}

	this.Data["json"] = space
	this.ServeJson()
}

// @Title GetAll
// @Description Get all spaces
// @Success 200 {object} models.Space
// @Failure	401 Not authorized
// @Failure	500	Failed to get all machines
// @router / [get]
func (this *SpacesController) GetAll() {

	// This is admin and staff only
	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	spaces, err := models.GetAllSpaces()
	if err != nil {
		beego.Error("Failed to get all spaces:", err)
		this.CustomAbort(500, "Failed to get all spaces")
	}

	this.Data["json"] = spaces
	this.ServeJson()
}
