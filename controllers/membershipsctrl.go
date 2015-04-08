package controllers

import (
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
	beego.Trace(membershipName)

	// All clear - create membership in the database
	var membershipId int64
	var err error
	membershipId, err = models.CreateMembership(membershipName)
	if err != nil {
		beego.Error("Failed to create membership", err)
		this.CustomAbort(403, "Failed to create membership")
	}

	this.Data["json"] = membershipId
	this.ServeJson()
}
