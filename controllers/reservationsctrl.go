package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego"
	"io/ioutil"
)

type ReservationsController struct {
	Controller
}

// @Title GetAll
// @Description Get all reservations
// @Success 200 {object}
// @Failure	403	Failed to get all reservations
// @Failure	401 Not authorized
// @router / [get]
func (this *ReservationsController) GetAll() {
	locId, authorized := this.GetLocIdMember()
	if !authorized {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	reservations, err := purchases.GetAllReservationsAt(locId)
	if err != nil {
		this.CustomAbort(403, "Failed to get all reservations")
	}
	this.Data["json"] = reservations
	this.ServeJSON()
}

// @Title Get
// @Description Get reservation by ID
// @Param	rid		path 	int	true		"Reservation ID"
// @Success 200 {object}
// @Failure	403	Failed to get reservation
// @Failure	401	Not authorized
// @router /:rid [get]
func (this *ReservationsController) Get() {
	rid, err := this.GetInt64(":rid")
	if err != nil {
		beego.Error("Failed to get :rid variable")
		this.CustomAbort(403, "Failed to get reservation")
	}

	reservation, err := purchases.GetReservation(rid)
	if err != nil {
		beego.Error("Failed to get reservation", err)
		this.CustomAbort(403, "Failed to get reservation")
	}

	this.Data["json"] = reservation
	this.ServeJSON()
}

// @Title Create
// @Description Create reservation
// @Param	model	body	string	true	"Reservation Name"
// @Success 200 {object}
// @Failure	400	Bad request
// @Failure	401	Not authorized
// @Failure	500	Failed to create reservation
// @router / [post]
func (this *ReservationsController) Create() {
	locId, authorized := this.GetLocIdMember()
	if !authorized {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := purchases.Reservation{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(400, "Bad request")
	}

	if locId == 0 || locId == req.LocationId() {
		id, err := purchases.CreateReservation(&req)
		if err != nil {
			beego.Error("Failed to create reservation", err)
			this.CustomAbort(500, "Failed to create reservation")
		}
		this.Data["json"] = purchases.ReservationCreatedResponse{Id: id}
	} else {
		this.CustomAbort(401, "Not authorized")
	}

	this.ServeJSON()
}

// @Title Put
// @Description Update reservation
// @Success 201 {object}
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:id [put]
func (this *ReservationsController) Put() {
	id, err := this.GetInt64(":id")
	if err != nil {
		this.Abort("400")
	}

	existing, err := purchases.GetReservation(id)
	if err != nil {
		beego.Error("get reservation:", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(existing.LocationId()) {
		beego.Error("Unauthorized attempt to update user")
		this.Abort("401")
	}

	reservation := &purchases.Reservation{}

	buf, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		beego.Error("Failed to read all:", err)
		this.Abort("400")
	}
	beego.Info("buf:", string(buf))
	defer this.Ctx.Request.Body.Close()

	data := bytes.NewBuffer(buf)

	dec := json.NewDecoder(data)
	if err := dec.Decode(reservation); err != nil {
		beego.Error("Failed to decode json:", err)
		this.Abort("400")
	}

	if reservation.Id() != id || reservation.LocationId() != existing.LocationId() {
		beego.Error("reservation id or location id changed")
		this.Abort("403")
	}

	if err := reservation.Update(); err != nil {
		beego.Error("Failed to update reservation:", err)
		this.Abort("500")
	}

	this.Data["json"] = reservation
	this.ServeJSON()
}
