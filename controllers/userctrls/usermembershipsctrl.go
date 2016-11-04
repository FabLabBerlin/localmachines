package userctrls

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships/inv_user_memberships"
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
			msg := fmt.Sprintf("error getting this month' invoice: %v", err)
			beego.Error(msg)
			this.CustomAbort(500, msg)
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

	ums, err := user_memberships.GetAllOfDeep(locationId, uid)
	if err != nil {
		beego.Error("Failed to get user machine permissions")
		this.Abort("500")
	}

	this.Data["json"] = ums
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
	beego.Info("UserMembershipsController#PutUserMembership()")
	dec := json.NewDecoder(this.Ctx.Request.Body)
	var um user_memberships.UserMembership

	if err := dec.Decode(&um); err != nil {
		beego.Error("Failed to decode json", err)
		this.Fail("500")
	}

	if !this.IsAdminAt(um.LocationId) {
		beego.Error("Not authorized")
		this.Fail(401, "Not authorized")
	}

	o := orm.NewOrm()

	if err := o.Begin(); err != nil {
		this.Fail(500, "tx begin")
	}

	if err := um.Update(o); err != nil {
		o.Rollback()
		beego.Error("UpdateMembership:", err)
		this.Fail(500)
	}

	iums, err := inv_user_memberships.GetForUserMembershipOrm(o, um.Id)
	if err != nil {
		o.Rollback()
		beego.Error("get invoice user memberships:", err)
		this.Fail(500)
	}

	nonDraftInvoiceErrors := 0

	for _, ium := range iums {
		if err := ium.Denormalize(o, um); err == inv_user_memberships.ErrNonDraftInvoice {
			nonDraftInvoiceErrors++
			continue
		} else if err != nil {
			o.Rollback()
			beego.Error("Denormalize:", err.Error())
			this.Fail(500, err.Error())
		}
	}

	if nonDraftInvoiceErrors > 0 && nonDraftInvoiceErrors == len(iums) {
		o.Rollback()
		this.Fail(500, "changes would only affect billed invoices")
	}

	if err := o.Commit(); err != nil {
		this.Fail(500, "tx commit")
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
		this.Fail("400")
	}

	um, err := user_memberships.Get(umid)
	if err != nil {
		beego.Error("Get user membership:", err)
		this.Fail("500")
	}

	if !this.IsAdminAt(um.LocationId) {
		beego.Error("Not authorized")
		this.Fail(401, "Not authorized")
	}

	o := orm.NewOrm()

	if err := o.Begin(); err != nil {
		this.Fail(500, "tx begin")
	}

	if err = inv_user_memberships.DeleteForUserMembership(o, umid); err != nil {
		o.Rollback()
		msg := fmt.Sprintf("delete invoice user membership: %v", err)
		beego.Error(msg)
		this.Fail("500", msg)
	}

	if err = user_memberships.Delete(o, umid); err != nil {
		o.Rollback()
		msg := fmt.Sprintf("delete user membership: %v", err)
		beego.Error(msg)
		this.Fail("500", msg)
	}

	if err := o.Commit(); err != nil {
		this.Fail(500, "tx commit")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}
