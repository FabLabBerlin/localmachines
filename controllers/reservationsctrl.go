package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/lib/icalendar"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
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
// @router /:rid([0-9]+) [get]
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

	inv, err := invoices.ThisMonthInvoice(locId, req.UserId())
	if err != nil {
		beego.Error("getting this month' invoice:", err)
		this.Abort("500")
	}
	req.Purchase.InvoiceId = inv.Id

	if locId == 0 || locId == req.LocationId() {
		_, err := purchases.CreateReservation(&req)
		if err != nil {
			beego.Error("Failed to create reservation", err)
			this.Abort("500")
		}
		this.Data["json"] = req.Purchase
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
		beego.Error("Unauthorized attempt to update reservation")
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

// @Title Cancel
// @Description Cancel reservation
// @Success 201 {object}
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:id/cancel [post]
func (this *ReservationsController) Cancel() {
	id, err := this.GetInt64(":id")
	if err != nil {
		this.Abort("400")
	}

	r, err := purchases.GetReservation(id)
	if err != nil {
		beego.Error("get reservation:", err)
		this.Abort("500")
	}

	uid, err := this.GetSessionUserId()
	if err != nil {
		beego.Error("cannot get session user id")
		this.Abort("500")
	}
	if r.UserId() != uid && !this.IsAdminAt(r.LocationId()) {
		beego.Error("Unauthorized attempt to update user")
		this.Abort("401")
	}

	r.Purchase.Cancelled = true

	if err := r.Update(); err != nil {
		beego.Error("update:", err)
		this.Abort("500")
	}

	this.Data["json"] = r
	this.ServeJSON()
}

// @Title ICalendar
// @Description Get iCalendar export
// @Param	location	query	int	true	"Location ID"
// @Param	machine		query	int	false	"Machine ID"
// @Success 200		{object}
// @Failure	403		Failed to get reservation
// @Failure	401		Not authorized
// @router /icalendar [get]
func (this *ReservationsController) ICalendar() {
	locationId, err := this.GetInt64("location")
	if err != nil {
		this.Abort("400")
	}

	machineId, err := this.GetInt64("machine")
	if err != nil {
		machineId = 0
	}

	rs, err := purchases.GetAllReservationsAt(locationId)
	if err != nil {
		beego.Error("get all reservations:", err)
		this.Abort("500")
	}

	ms, err := machine.GetAllAt(locationId)
	if err != nil {
		beego.Error("get all users:", err)
		this.Abort("500")
	}
	msById := make(map[int64]*machine.Machine)
	for _, m := range ms {
		msById[m.Id] = m
	}

	uls, err := user_locations.GetAllForLocation(locationId)
	if err != nil {
		beego.Error("get all user locations:", err)
		this.Abort("500")
	}
	rolesByUserId := make(map[int64]user_roles.Role)
	for _, ul := range uls {
		rolesByUserId[ul.UserId] = ul.GetRole()
	}

	events := make([]icalendar.Event, 0, len(rs))
	for _, r := range rs {
		e := icalendar.Event{
			Reservation: r,
			Machine:     msById[r.Purchase.MachineId],
			UserRole:    rolesByUserId[r.Purchase.UserId],
		}
		if machineId == 0 || r.Purchase.MachineId == machineId {
			events = append(events, e)
		}
	}

	this.Ctx.Output.ContentType("ics")
	this.Ctx.WriteString(icalendar.ToIcal(events))
	this.Finish()
}
