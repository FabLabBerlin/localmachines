package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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
	m, err := memberships.Get(mId)
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

	t := time.Now()

	invoiceIds := make([]int64, 0, 1)

	for i := 0; ; i++ {
		inv, err := invoices.GetDraft(m.LocationId, uid, t)
		if err != nil {
			beego.Error("error getting this month' invoice:", err)
			this.Abort("500")
		}

		invoiceIds = append(invoiceIds, inv.Id)

		if t.Month() == startDate.Month() && t.Year() == startDate.Year() {
			break
		} else {
			t = t.AddDate(0, -1, 0)
		}

		if i > 100 {
			beego.Error("loop executed more than 100x")
			this.Abort("500")
		}
	}

	if len(invoiceIds) > 2 {
		beego.Error("more than 2 invoice months would be affected")
		this.Abort("500")
	}

	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		beego.Error("begin tx:", err)
		this.Abort("500")
	}

	for _, invId := range invoiceIds {
		_, err = user_memberships.Create(o, uid, mId, invId, startDate)
		if err != nil {
			beego.Error("Error creating user membership:", err)
			this.Abort("500")
		}
	}

	if err := o.Commit(); err != nil {
		beego.Error("commit tx:", err)
		this.Abort("500")
	}

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

	inv, err := invoices.GetDraft(locationId, uid, time.Now())
	if err != nil {
		beego.Error("current invoice:", err)
		this.Abort("500")
	}

	beego.Info("current inv.Id=", inv.Id)

	all, err := user_memberships.GetForInvoice(inv.Id)
	if err != nil {
		beego.Error("Failed to get user machine permissions")
		this.Abort("500")
	}

	list := &user_memberships.List{
		Data: make([]*user_memberships.Combo, 0, len(all.Data)),
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
	var um user_memberships.UserMembership

	if err := dec.Decode(&um); err != nil {
		beego.Error("Failed to decode json", err)
		this.Abort("500")
	}

	inv, err := invoices.Get(um.InvoiceId)
	if err != nil {
		beego.Error("get invoice:", err)
		this.Abort("500")
	}

	if inv.Status != "draft" {
		beego.Error("cannot delete user membership because it's bound to non-draft invoice")
		this.Abort("403")
	}

	if !this.IsAdminAt(inv.LocationId) {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	if err := um.Update(); err != nil {
		beego.Error("UpdateMembership: ", err)
		this.Abort("500")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title Delete
// @Description Delete UserMembership
// @Param	uid		path 	int	true						"User Membership Id"
// @Success	200	ok
// @Failure	400	Variable message
// @Failure	401	Unauthorized
// @Failure	403	Variable message
// @router /:uid/memberships/:umid [delete]
func (this *UserMembershipsController) DeleteUserMembership() {
	umid, err := this.GetInt64(":umid")
	if err != nil {
		this.Abort("400")
	}

	um, err := user_memberships.Get(umid)
	if err != nil {
		beego.Error("Get user membership:", err)
		this.Abort("500")
	}

	inv, err := invoices.Get(um.InvoiceId)
	if err != nil {
		beego.Error("get invoice:", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(inv.LocationId) {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	if inv.Status != "draft" {
		beego.Error("cannot delete user membership because it's bound to non-draft invoice")
		this.Abort("403")
	}
	err = user_memberships.Delete(umid)
	if err != nil {
		beego.Error("delete user membership:", err)
		this.Abort("500")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}
