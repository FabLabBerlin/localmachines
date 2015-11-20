package controllers

import (
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
