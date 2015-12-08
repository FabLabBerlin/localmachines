package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
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

	// Check if logged in
	uid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if uid == nil {
		beego.Info("Attempt to get all users while not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	if !this.IsAdmin(uid.(int64)) && !this.IsStaff(uid.(int64)) {
		beego.Error("Not authorized to get all memberships")
		this.CustomAbort(401, "Not authorized")
	}

	memberships, err := models.GetAllMemberships()
	if err != nil {
		this.CustomAbort(403, "Failed to get all memberships")
	}
	this.Data["json"] = memberships
	this.ServeJson()
}

// @Title Create
// @Description Create new membership
// @Param	mname	query	string	true	"Membership Name"
// @Success	200	int	Membership ID
// @Failure	403	Failed to create membership
// @Failure	401	Not authorized
// @router / [post]
func (this *MembershipsController) Create() {

	if !this.IsAdmin() {
		beego.Error("Not authorized to create membership")
		this.CustomAbort(401, "Not authorized")
	}

	membershipName := this.GetString("mname")

	id, err := models.CreateMembership(membershipName)
	if err != nil {
		beego.Error("Failed to create membership", err)
		this.CustomAbort(403, "Failed to create membership")
	}

	this.Data["json"] = id
	this.ServeJson()
}

// @Title Get
// @Description Get membership by membership ID
// @Param	mid		path 	int	true		"Membership ID"
// @Success 200 {object} models.Membership
// @Failure	403	Failed to get membership
// @Failure	401	Not authorized
// @router /:mid [get]
func (this *MembershipsController) Get() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get mid")
		this.CustomAbort(403, "Failed to get membership")
	}

	membership, err := models.GetMembership(mid)
	if err != nil {
		beego.Error("Could not get membership")
		this.CustomAbort(403, "Failed to get membership")
	}

	this.Data["json"] = membership
	this.ServeJson()
}

// @Title Update
// @Description Update membership
// @Param	mid	path	int	true	"Membership ID"
// @Param	model	body	models.Membership	true	"Membership model"
// @Success 200 string ok
// @Failure	403	Failed to update membership
// @Failure	401	Not authorized
// @router /:mid [put]
func (this *MembershipsController) Update() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	// Attempt to decode passed json
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.Membership{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update membership")
	}

	// Get mid and check if it matches with the membership model ID
	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get :mid:", err)
		this.CustomAbort(403, "Failed to update membership")
	}
	if mid != req.Id {
		beego.Error("mid and model ID do not match:", err)
		this.CustomAbort(403, "Failed to update membership")
	}

	if err = req.Update(); err != nil {
		beego.Error("Failed updating membership:", err)
		this.CustomAbort(403, "Failed to update membership")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title Delete
// @Description Delete membership
// @Param	mid	path	int	true	"Membership ID"
// @Success 200 string ok
// @Failure	403	Failed to delete membership
// @Failure	401	Not authorized
// @router /:mid [delete]
func (this *MembershipsController) Delete() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get mid:", err)
		this.CustomAbort(403, "Failed to delete membership")
	}

	err = models.DeleteMembership(mid)
	if err != nil {
		beego.Error("Failed to delete membership:", err)
		this.CustomAbort(403, "Failed to delete membership")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
