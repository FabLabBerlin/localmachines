package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type ReservationsController struct {
	Controller
}

// @Title GetAll
// @Description Get all reservations
// @Success 200 {object} models.Reservation
// @Failure	403	Failed to get all reservations
// @Failure	401 Not authorized
// @router / [get]
func (this *ReservationsController) GetAll() {
	reservations, err := models.GetAllReservations()
	if err != nil {
		this.CustomAbort(403, "Failed to get all reservations")
	}
	this.Data["json"] = reservations
	this.ServeJson()
}

// @Title Get
// @Description Get reservation by ID
// @Param	mid		path 	int	true		"Reservation ID"
// @Success 200 {object} models.Reservation
// @Failure	403	Failed to get reservation
// @Failure	401	Not authorized
// @router /:rid [get]
func (this *ReservationsController) Get() {
	rid, err := this.GetInt64(":rid")
	if err != nil {
		beego.Error("Failed to get :rid variable")
		this.CustomAbort(403, "Failed to get reservation")
	}

	reservation, err := models.GetReservation(rid)
	if err != nil {
		beego.Error("Failed to get reservation", err)
		this.CustomAbort(403, "Failed to get reservation")
	}

	this.Data["json"] = reservation
	this.ServeJson()
}

// @Title Create
// @Description Create reservation
// @Param	model	body	models.Reservation	true	"Reservation Name"
// @Success 200 {object} models.ReservationCreatedResponse
// @Failure	403	Failed to create reservation
// @Failure	401	Not authorized
// @router / [post]
func (this *ReservationsController) Create() {
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.Reservation{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to create reservation")
	}
	beego.Info("create reservation:", req)

	id, err := models.CreateReservation(&req)
	if err != nil {
		beego.Error("Failed to create reservation", err)
		this.CustomAbort(403, "Failed to create reservation")
	}

	this.Data["json"] = models.ReservationCreatedResponse{Id: id}
	this.ServeJson()
}

// @Title Delete
// @Description Delete reservation
// @Param	rid	path	int	true	"Reservation ID"
// @Success 200 string ok
// @Failure	403	Failed to delete reservation
// @Failure	401	Not authorized
// @router /:rid [delete]
func (this *ReservationsController) Delete() {
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	rid, err := this.GetInt64(":rid")
	if err != nil {
		beego.Error("Failed to get rid:", err)
		this.CustomAbort(403, "Failed to delete reservation")
	}

	err = models.DeleteReservation(rid)
	if err != nil {
		beego.Error("Failed to delete reservation:", err)
		this.CustomAbort(403, "Failed to delete reservation")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
