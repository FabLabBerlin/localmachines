package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/astaxie/beego"
	"time"
)

type UserMembershipsController struct {
	Controller
}

// @Title PostUserMemberships
// @Description Post user membership
// @Param	uid							path 		int			true		"User ID"
// @Param	membershipId		query 	int			true		"Membership ID"
// @Param	startDate				query 	string	true		"Membership ID"
// @Success 200 {object} models.UserMembership
// @Failure	401	Not authorized
// @Failure	500	Failed to get user memberships
// @router /:uid/memberships [post]
func (this *UserMembershipsController) PostUserMemberships() {

	if !this.IsAdmin() {
		this.CustomAbort(401, "Not authorized")
	}

	// Get request user ID
	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID:", err)
		this.CustomAbort(403, "Failed to update user memberships")
	}
	if uid <= 0 {
		beego.Error("Wrong User ID:", uid)
		this.CustomAbort(403, "Failed to update user memberships")
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
	userMembershipId, err := models.CreateUserMembership(uid, membershipId, startDate)
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
// @Param	uid			path 	int	true		"User ID"
// @Param	location	query	int	false		"Location ID"
// @Success 200 {object} models.UserMembershipList
// @Failure	403	Failed to get user memberships
// @Failure	401	Not authorized
// @router /:uid/memberships [get]
func (this *UserMembershipsController) GetUserMemberships() {

	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(400, "Wrong uid in url or not authorized")
	}
	locationId, _ := this.GetInt64("location")

	all, err := models.GetUserMemberships(uid)
	if err != nil {
		beego.Error("Failed to get user machine permissions")
		this.CustomAbort(403, "Failed to get user machines")
	}

	userMemberships := &models.UserMembershipList{
		Data: make([]models.UserMembershipCombo, 0, len(all.Data)),
	}
	for _, um := range all.Data {
		if locationId <= 0 || locationId == um.LocationId {
			userMemberships.Data = append(userMemberships.Data, um)
		}
	}

	this.Data["json"] = userMemberships
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
