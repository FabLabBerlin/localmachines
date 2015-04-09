package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type HexabusController struct {
	Controller
}

// @Title Get
// @Description Get hexabus mapping by by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 {object} models.HexabusMapping
// @Failure	403	Failed to get mapping
// @Failure	401	Not authorized
// @router /:mid [get]
func (this *HexabusController) Get() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error
	var mid int64

	mid, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get :mid:", err)
		this.CustomAbort(403, "Failed to get mapping")
	}

	var mapping *models.HexabusMapping
	mapping, err = models.GetHexabusMapping(mid)
	if err != nil {
		beego.Error("Failed to get hexabus maping")
		this.CustomAbort(403, "Failed to get mapping")
	}

	this.Data["json"] = mapping
	this.ServeJson()
}

// @Title Create
// @Description Create hexabus mapping with machine ID
// @Param	mid		query 	int	true		"Machine ID"
// @Success 200 int	Mapping ID
// @Failure	403	Failed to create mapping
// @Failure	401	Not authorized
// @router / [post]
func (this *HexabusController) Create() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var mid int64
	var err error

	mid, err = this.GetInt64("mid")
	if err != nil {
		beego.Error("Could not get mid:", err)
		this.CustomAbort(403, "Failed to create mapping")
	}

	var mappingId int64

	mappingId, err = models.CreateHexabusMapping(mid)
	if err != nil {
		beego.Error("Failed to create hexabus mapping:", err)
		this.CustomAbort(403, "Failed to create mapping")
	}

	this.Data["json"] = mappingId
	this.ServeJson()
}

// @Title Delete
// @Description Delete hexabus mapping by by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 string ok
// @Failure	403	Failed to delete mapping
// @Failure	401	Not authorized
// @router /:mid [delete]
func (this *HexabusController) Delete() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error
	var mid int64

	mid, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid:", err)
		this.CustomAbort(403, "Failed to delete mapping")
	}

	err = models.DeleteHexabusMapping(mid)
	if err != nil {
		beego.Error("Failed to delete hexabus mapping:", err)
		this.CustomAbort(403, "Failed to delete hexabus mapping")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title Update
// @Description Update hexabus mapping by by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Param	model	body	models.HexabusMapping	true	"Hexabus mapping model"
// @Success 200 string ok
// @Failure	403	Failed to update mapping
// @Failure	401	Not authorized
// @router /:mid [put]
func (this *HexabusController) Update() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error

	// Attempt to decode passed json
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.HexabusMapping{}
	if err = dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update mapping")
	}

	var mid int64

	mid, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid:", err)
		this.CustomAbort(403, "Failed to update mapping")
	}

	// Check if IDs match
	if mid != req.MachineId {
		beego.Error("mid and model machine ID do not match:", err)
		this.CustomAbort(403, "Failed to update mapping")
	}

	err = models.UpdateHexabusMapping(&req)
	if err != nil {
		beego.Error("Failed updating mapping:", err)
		this.CustomAbort(403, "Failed to update mapping")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
