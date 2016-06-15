package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/astaxie/beego"
	"time"
)

type UserMembershipsController struct {
	Controller
}

// @Title PostUserMemberships
// @Description Post user membership
// @Param	uid				path 	int		true	"User ID"
// @Param	membershipId	query 	int		true	"Membership ID"
// @Param	startDate		query 	string	true	"Start Date"
// @Success 200 {object} models.UserMembership
// @Failure	400	Bad request
// @Failure	401	Not authorized
// @Failure	500	Failed to get user memberships
// @router /:uid/memberships [post]
func (this *UserMembershipsController) PostUserMemberships() {

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID:", err)
		this.Abort("400")
	}
	if uid <= 0 {
		beego.Error("Wrong User ID:", uid)
		this.Abort("400")
	}

	mId, err := this.GetInt64("membershipId")
	if err != nil {
		beego.Error("Failed to get membership ID")
		this.Abort("400")
	}
	m, err := memberships.GetMembership(mId)
	if err != nil {
		beego.Error("get membership:", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(m.LocationId) {
		beego.Error("not authorized for", m.LocationId)
		this.CustomAbort(401, "Not authorized")
	}

	startDate, err := time.ParseInLocation("2006-01-02",
		this.GetString("startDate"),
		time.UTC)
	if err != nil {
		beego.Error("Failed to parse startDate=", startDate)
		this.Abort("500")
	}

	userMembership, err := memberships.CreateUserMembership(uid, mId, startDate)
	if err != nil {
		beego.Error("Error creating user membership:", err)
		this.Abort("500")
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

	if locationId <= 0 {
		beego.Error("location id:", locationId)
		this.Abort("400")
	}

	inv, err := invoices.ThisMonthInvoice(locationId, uid)
	if err != nil {
		beego.Error("current invoice:", err)
		this.Abort("500")
	}

	beego.Info("current inv.Id=", inv.Id)

	all, err := memberships.GetUserMembershipsForInvoice(inv.Id)
	if err != nil {
		beego.Error("Failed to get user machine permissions")
		this.Abort("500")
	}

	list := &memberships.UserMembershipList{
		Data: make([]memberships.UserMembershipCombo, 0, len(all.Data)),
	}
	for _, um := range all.Data {
		if locationId <= 0 || locationId == um.LocationId {
			list.Data = append(list.Data, um)
		}
	}

	this.Data["json"] = list
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
	var um memberships.UserMembership

	if err := dec.Decode(&um); err != nil {
		beego.Error("Failed to decode json", err)
		this.Abort("500")
	}

	m, err := memberships.GetMembership(um.MembershipId)
	if err != nil {
		beego.Error("Get membership:", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(m.LocationId) {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	if err := m.Update(); err != nil {
		beego.Error("UpdateMembership: ", err)
		this.Abort("500")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}
