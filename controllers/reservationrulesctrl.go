package controllers

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/astaxie/beego"
)

type ReservationRulesController struct {
	Controller
}

// @Title GetAll
// @Description Get all ReservationRules
// @Success 200 {object} models.ReservationRule
// @Failure	403	Failed to get all ReservationRules
// @Failure	401 Not authorized
// @router / [get]
func (this *ReservationRulesController) GetAll() {
	locId, authorized := this.GetLocIdMember()
	if !authorized {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	ReservationRules, err := models.GetAllReservationRulesAt(locId)
	if err != nil {
		beego.Error("Failed to get all ReservationRules:", err)
		this.CustomAbort(403, "Failed to get all ReservationRules")
	}
	this.Data["json"] = ReservationRules
	this.ServeJSON()
}

// @Title Get
// @Description Get ReservationRule by ID
// @Param	mid		path 	int	true		"ReservationRule ID"
// @Success 200 {object} models.ReservationRule
// @Failure	403	Failed to get ReservationRule
// @Failure	401	Not authorized
// @router /:rid [get]
func (this *ReservationRulesController) Get() {
	rid, err := this.GetInt64(":rid")
	if err != nil {
		beego.Error("Failed to get :rid variable")
		this.CustomAbort(403, "Failed to get ReservationRule")
	}

	ReservationRule, err := models.GetReservationRule(rid)
	if err != nil {
		beego.Error("Failed to get ReservationRule", err)
		this.CustomAbort(403, "Failed to get ReservationRule")
	}

	this.Data["json"] = ReservationRule
	this.ServeJSON()
}

// @Title Create
// @Description Create ReservationRule
// @Param	model	body	models.ReservationRule	true	"ReservationRule Name"
// @Success 200 {object} models.ReservationRuleCreatedResponse
// @Failure	403	Failed to create ReservationRule
// @Failure	401	Not authorized
// @router / [post]
func (this *ReservationRulesController) Create() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.ReservationRule{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to create ReservationRule")
	}
	beego.Info("create ReservationRule:", req)
	req.LocationId = locId

	id, err := models.CreateReservationRule(&req)
	if err != nil {
		beego.Error("Failed to create ReservationRule", err)
		this.CustomAbort(403, "Failed to create ReservationRule")
	}

	this.Data["json"] = models.ReservationRuleCreatedResponse{Id: id}
	this.ServeJSON()
}

// @Title Update
// @Description Update ReservationRule
// @Param	model	body	models.ReservationRule	true	"ReservationRule"
// @Success 200 {object} models.ReservationRule
// @Failure	403	Failed to update ReservationRule
// @Failure	401	Not authorized
// @router /:rid [put]
func (this *ReservationRulesController) Update() {
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.ReservationRule{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update ReservationRule")
	}

	if err := req.Update(); err != nil {
		beego.Error("Failed to update ReservationRule", err)
		this.CustomAbort(403, "Failed to update ReservationRule")
	}

	this.Data["json"] = req
	this.ServeJSON()
}

// @Title Delete
// @Description Delete ReservationRule
// @Param	rid	path	int	true	"ReservationRule ID"
// @Success 200 string ok
// @Failure	403	Failed to delete ReservationRule
// @Failure	401	Not authorized
// @router /:rid [delete]
func (this *ReservationRulesController) Delete() {
	rid, err := this.GetInt64(":rid")
	if err != nil {
		beego.Error("Failed to get rid:", err)
		this.Abort("400")
	}

	rule, err := models.GetReservationRule(rid)
	if err != nil {
		beego.Error("Failed to get rule:", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(rule.LocationId) {
		beego.Error("Not authorized")
		this.Abort("401")
	}

	err = models.DeleteReservationRule(rid)
	if err != nil {
		beego.Error("Failed to delete ReservationRule:", err)
		this.Abort("500")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}
