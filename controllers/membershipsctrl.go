package controllers

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/astaxie/beego"
)

type MembershipsController struct {
	Controller
}

// Override our custom root controller's Prepare method as it is checking
// if we are logged in and we don't want that here at this point
func (this *MembershipsController) Prepare() {
	beego.Info("Skipping global login check")
}

// @Title GetAll
// @Description Get all memberships
// @Success 200 {object} models.Membership
// @Failure	403	Failed to get all memberships
// @router / [get]
func (this *MembershipsController) GetAll() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	memberships, err := models.GetAllMembershipsAt(locId)
	if err != nil {
		this.CustomAbort(403, "Failed to get all memberships")
	}
	this.Data["json"] = memberships
	this.ServeJSON()
}

// @Title Create
// @Description Create new membership
// @Param	mname	query	string	true	"Membership Name"
// @Success	200	int	Membership ID
// @Failure	403	Failed to create membership
// @Failure	401	Not authorized
// @router / [post]
func (this *MembershipsController) Create() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	membershipName := this.GetString("mname")

	id, err := models.CreateMembership(locId, membershipName)
	if err != nil {
		beego.Error("Failed to create membership", err)
		this.CustomAbort(403, "Failed to create membership")
	}

	this.Data["json"] = id
	this.ServeJSON()
}

// @Title Get
// @Description Get membership by membership ID
// @Param	mid		path 	int	true		"Membership ID"
// @Success 200 {object} models.Membership
// @Failure	403	Failed to get membership
// @Failure	401	Not authorized
// @router /:mid [get]
func (this *MembershipsController) Get() {
	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get mid")
		this.CustomAbort(500, "Failed to get membership")
	}

	membership, err := models.GetMembership(mid)
	if err != nil {
		beego.Error("Could not get membership")
		this.CustomAbort(500, "Failed to get membership")
	}

	if !this.IsAdminAt(membership.LocationId) {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	this.Data["json"] = membership
	this.ServeJSON()
}

// @Title Update
// @Description Update membership
// @Param	mid	path	int	true	"Membership ID"
// @Param	model	body	models.Membership	true	"Membership model"
// @Success 200 string ok
// @Failure	401	Not authorized
// @Failure	403	Forbidden
// @Failure	500	Internal Server Error
// @router /:mid [put]
func (this *MembershipsController) Update() {

	// Get mid and check if it matches with the membership model ID
	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get :mid:", err)
		this.CustomAbort(500, "Failed to update membership")
	}

	existing, err := models.GetMembership(mid)
	if err != nil {
		beego.Error("Get membership:", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	if !this.IsAdminAt(existing.LocationId) {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	// Attempt to decode passed json
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.Membership{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(400, "Failed to update membership")
	}

	if mid != req.Id {
		beego.Error("mid and model ID do not match:", err)
		this.CustomAbort(403, "Forbidden")
	}
	if req.LocationId != existing.LocationId {
		beego.Error("old and new location id do not match")
		this.CustomAbort(403, "Forbidden")
	}

	if err = req.Update(); err != nil {
		beego.Error("Failed updating membership:", err)
		this.CustomAbort(500, "Failed to update membership")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}
