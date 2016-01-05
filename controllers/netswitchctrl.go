package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type NetSwitchController struct {
	Controller
}

// @Title GetAll
// @Description Get all netswitch mapings
// @Success 200 {object} models.NetSwitchMapping
// @Failure	403	Failed to get all netswitch mappings
// @router / [get]
func (this *NetSwitchController) GetAll() {

	// This is admin and staff only
	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	ms, err := models.GetAllNetSwitchMapping()
	if err != nil {
		this.CustomAbort(403, "Failed to get all netswitch mappings")
	}
	this.Data["json"] = ms
	this.ServeJson()
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

	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get :mid:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	mapping, err := models.GetNetSwitchMapping(mid)
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

	mid, err := this.GetInt64("mid")
	if err != nil {
		beego.Error("Could not get mid:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	mappingId, err := models.CreateNetSwitchMapping(mid)
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

	mid, err := this.GetInt64(":mid")
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

// @Title Update
// @Description Update NetSwitch mapping by by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Param	model	body	models.NetSwitchMapping	true	"NetSwitch mapping model"
// @Success 200 string ok
// @Failure	500	Internal Server Error
// @Failure	401	Not authorized
// @router /:mid [put]
func (this *NetSwitchController) Update() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error

	// Attempt to decode passed json
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.NetSwitchMapping{}
	if err = dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	var mid int64

	mid, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Check if IDs match
	if mid != req.MachineId {
		beego.Error("mid and model machine ID do not match:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	if err = req.Update(); err != nil {
		beego.Error("Failed updating mapping:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
