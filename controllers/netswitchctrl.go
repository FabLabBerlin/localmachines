package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type NetSwitchController struct {
	Controller
}

// @Title Get
// @Description Get NetSwitch mapping by by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 {object} models.NetSwitchMapping
// @Failure	500	Internal Server Error
// @Failure	401	Not authorized
// @router /:mid [get]
func (this *NetSwitchController) Get() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error
	var mid int64

	mid, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get :mid:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	var mapping *models.NetSwitchMapping
	mapping, err = models.GetNetSwitchMapping(mid)
	if err != nil {
		beego.Error("Failed to get NetSwitch maping")
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = mapping
	this.ServeJson()
}

// @Title Create
// @Description Create UrlSwitch mapping with machine ID
// @Param	mid		query 	int	true		"Machine ID"
// @Success 200 int	Mapping ID
// @Failure	500	Internal Server Error
// @Failure	401	Not authorized
// @router / [post]
func (this *NetSwitchController) Create() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var mid int64
	var err error

	mid, err = this.GetInt64("mid")
	if err != nil {
		beego.Error("Could not get mid:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	var mappingId int64

	mappingId, err = models.CreateNetSwitchMapping(mid)
	if err != nil {
		beego.Error("Failed to create NetSwitch mapping:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = mappingId
	this.ServeJson()
}

// @Title Delete
// @Description Delete NetSwitch mapping by by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 string ok
// @Failure	500	Internal Server Error
// @Failure	401	Not authorized
// @router /:mid [delete]
func (this *NetSwitchController) Delete() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error
	var mid int64

	mid, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid:", err)
		this.CustomAbort(403, "Internal Server Error")
	}

	err = models.DeleteNetSwitchMapping(mid)
	if err != nil {
		beego.Error("Failed to delete NetSwitch mapping:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
