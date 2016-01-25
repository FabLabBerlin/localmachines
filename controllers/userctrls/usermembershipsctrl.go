package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/astaxie/beego"
	"time"
)

type UserMembershipsController struct {
	controllers.Controller
}

// @Title PostUserMemberships
// @Description Post user membership
// @Param	uid							path 		int			true		"User ID"
// @Param	membershipId		query 	int			true		"Membership ID"
// @Param	startDate				query 	string	true		"Membership ID"
// @Success 200 {object} models.UserMembership
// @Failure	403	Failed to get user memberships
// @Failure	401	Not authorized
// @router /:uid/memberships [post]
func (this *UserMembershipsController) PostUserMemberships() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	// Get requested user ID
	ruid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get :uid")
	}
	if ruid <= 0 {
		beego.Error("Bad :uid")
		this.CustomAbort(500, "Bad uid")
	}

	// Get requested user membership ID
	membershipId, err := this.GetInt64("membershipId")
	if err != nil {
		beego.Error("Failed to get membership ID")
		this.CustomAbort(403, "Failed to get membership ID")
	}

	// Get requested start date
	startDate, err := time.ParseInLocation("2006-01-02",
		this.GetString("startDate"),
		time.UTC)
	if err != nil {
		beego.Error("Failed to parse startDate=%v", startDate)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Create user membership by using the model function
	userMembershipId, err := models.CreateUserMembership(ruid, membershipId, startDate)
	if err != nil {
		beego.Error("Error creating user membership:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Read the user membership back
	var userMembership *models.UserMembership
	userMembership, err = models.GetUserMembership(userMembershipId)
	if err != nil {
		beego.Error("Failed to get user membership:", err)
		this.CustomAbort(500, "Failed to get user membership")
	}

	this.Data["json"] = userMembership
	this.ServeJSON()
}

// @Title GetUserMemberships
// @Description Get user memberships
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.UserMembershipList
// @Failure	403	Failed to get user memberships
// @Failure	401	Not authorized
// @router /:uid/memberships [get]
func (this *UserMembershipsController) GetUserMemberships() {

	// Check if logged in
	suid, err := this.GetSessionUserId()
	if err != nil {
		beego.Info("Not logged in:", err)
		this.CustomAbort(401, "Not logged in")
	}

	// Get requested user ID
	var ruid int64
	ruid, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get :uid")
	}

	// We need the user roles in order to understand
	// whether we are allowed to access other user machines

	if suid != ruid {
		if !this.IsAdmin() {

			// The currently logged in user is not allowed to access
			// other user machines
			beego.Error("Not authorized")
			this.CustomAbort(401, "Not authorized")
		}
	}

	// If the requested user roles is not admin and staff
	// we need to get machine permissions first and then the machines
	var ums *models.UserMembershipList
	ums, err = models.GetUserMemberships(ruid)
	if err != nil {
		beego.Error("Failed to get user machine permissions")
		this.CustomAbort(403, "Failed to get user machines")
	}

	// Serve machines
	this.Data["json"] = ums
	this.ServeJSON()
}

// @Title DeleteUserMembership
// @Description Delete user membership
// @Param	uid		path 	int	true		"User ID"
// @Param	umid	path	int	true		"User Membership ID"
// @Success 200
// @Failure	403	Failed to get user memberships
// @Failure	401	Not authorized
// @router /:uid/memberships/:umid [delete]
func (this *UserMembershipsController) DeleteUserMembership() {
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	umid, err := this.GetInt64(":umid")
	if err != nil {
		beego.Error("Failed to get :umid")
		this.CustomAbort(403, "Failed to get :umid")
	}
	beego.Trace("User membership ID:", umid)

	err = models.DeleteUserMembership(umid)
	if err != nil {
		beego.Error("Failed to delete user membership")
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title Put
// @Description Update UserMembership
// @Param	uid		path 	int	true						"User Membership Id"
// @Param	body	body	models.UserMembership	true	"User Membership model"
// @Success	200	ok
// @Failure	400	Variable message
// @Failure	401	Unauthorized
// @Failure	403	Variable message
// @router /:uid/memberships/:umid [put]
func (this *UserMembershipsController) PutUserMembership() {
	dec := json.NewDecoder(this.Ctx.Request.Body)
	var userMembership models.UserMembership
	if err := dec.Decode(&userMembership); err == nil {
		beego.Info("userMembership: ", userMembership)
	} else {
		beego.Error("Failed to decode json", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	if err := userMembership.Update(); err != nil {
		beego.Error("UpdateMembership: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}
