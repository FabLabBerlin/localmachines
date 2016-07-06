package controllers

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models/memberships"
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
// @Success 200 {object} memberships.Membership
// @Failure	403	Failed to get all memberships
// @router / [get]
func (this *MembershipsController) GetAll() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	ms, err := memberships.GetAllAt(locId)
	if err != nil {
		this.CustomAbort(403, "Failed to get all memberships")
	}
	this.Data["json"] = ms
	this.ServeJSON()
}

// @Title Create
// @Description Create new membership
// @Param	mname	query	string	true	"Membership Name"
// @Success	200	int	Membership ID
// @Failure	401	Not authorized
// @Failure	500	Failed to create membership
// @router / [post]
func (this *MembershipsController) Create() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	name := this.GetString("mname")

	m, err := memberships.Create(locId, name)
	if err != nil {
		beego.Error("Failed to create membership", err)
		this.Abort("500")
	}

	this.Data["json"] = m.Id
	this.ServeJSON()
}

// @Title Get
// @Description Get membership by membership ID
// @Param	mid		path 	int	true		"Membership ID"
// @Success 200 {object} models.Membership
// @Failure	401	Not authorized
// @Failure	500	Failed to get membership
// @router /:mid [get]
func (this *MembershipsController) Get() {
	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get mid")
		this.Abort("400")
	}

	m, err := memberships.Get(mid)
	if err != nil {
		beego.Error("Could not get membership")
		this.Abort("500")
	}

	if !this.IsAdminAt(m.LocationId) {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	this.Data["json"] = m
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

	existing, err := memberships.Get(mid)
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
	req := memberships.Membership{}
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

// @Title SetArchived
// @Description (Un)archive membership
// @Param	mid		path	int		true	"Membership ID"
// @Param	archive	query	bool	true	"Archive"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500	Failed to archive membership
// @router /:mid/set_archived [post]
func (this *MembershipsController) SetArchived() {
	id, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.Abort("400")
	}

	m, err := memberships.Get(id)
	if err != nil {
		beego.Error("get", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(m.LocationId) {
		beego.Error("Not authorized")
		this.Abort("401")
	}

	m.Archived, err = this.GetBool("archived")
	if err != nil {
		beego.Error("parsing archived parameter")
		this.Abort("400")
	}

	if err = m.Update(); err != nil {
		beego.Error("update:", err)
		this.Abort("500")
	}

	this.Finish()
}
