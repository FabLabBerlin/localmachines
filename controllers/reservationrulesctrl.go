package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
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
	ReservationRules, err := models.GetAllReservationRules()
	if err != nil {
		beego.Error("Failed to get all ReservationRules:", err)
		this.CustomAbort(403, "Failed to get all ReservationRules")
	}
	this.Data["json"] = ReservationRules
	this.ServeJson()
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
	this.ServeJson()
}

// @Title Create
// @Description Create ReservationRule
// @Param	model	body	models.ReservationRule	true	"ReservationRule Name"
// @Success 200 {object} models.ReservationRuleCreatedResponse
// @Failure	403	Failed to create ReservationRule
// @Failure	401	Not authorized
// @router / [post]
func (this *ReservationRulesController) Create() {
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.ReservationRule{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to create ReservationRule")
	}
	beego.Info("create ReservationRule:", req)

	id, err := models.CreateReservationRule(&req)
	if err != nil {
		beego.Error("Failed to create ReservationRule", err)
		this.CustomAbort(403, "Failed to create ReservationRule")
	}

	this.Data["json"] = map[string]int64{"Id": id}
	this.ServeJson()
}

// @router /:rid [put]
func (this *ReservationRulesController) Update() {
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.ReservationRule{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update ReservationRule")
	}
	beego.Info("create ReservationRule:", req)

	err := models.UpdateReservationRule(&req)
	if err != nil {
		beego.Error("Failed to update ReservationRule", err)
		this.CustomAbort(403, "Failed to update ReservationRule")
	}

	this.Data["json"] = req
	this.ServeJson()
}

// @Title Delete
// @Description Delete ReservationRule
// @Param	rid	path	int	true	"ReservationRule ID"
// @Success 200 string ok
// @Failure	403	Failed to delete ReservationRule
// @Failure	401	Not authorized
// @router /:rid [delete]
func (this *ReservationRulesController) Delete() {
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	rid, err := this.GetInt64(":rid")
	if err != nil {
		beego.Error("Failed to get rid:", err)
		this.CustomAbort(403, "Failed to delete ReservationRule")
	}

	err = models.DeleteReservationRule(rid)
	if err != nil {
		beego.Error("Failed to delete ReservationRule:", err)
		this.CustomAbort(403, "Failed to delete ReservationRule")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
