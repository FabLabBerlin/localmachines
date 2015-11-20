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
// @Param	model	body	models.Space	true	"Reservation Name"
// @Success 200 {object} models.Space
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [post]
func (this *SpacesController) Create() {
	dec := json.NewDecoder(this.Ctx.Request.Body)
	space := models.Space{}
	if err := dec.Decode(&space); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(400, "Failed to create reservation")
	}

	_, err := models.CreateSpace(&space)
	if err != nil {
		beego.Error("Failed to create space", err)
		this.CustomAbort(500, "Failed to create space")
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
