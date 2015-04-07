package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type MachinesController struct {
	Controller
}

// @Title GetAll
// @Description Get all machines
// @Success 200 {object} models.Machine
// @Failure	403	Failed to get all machines
// @Failure	401 Not authorized
// @router / [get]
func (this *MachinesController) GetAll() {

	// This is admin and staff only
	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var machines []*models.Machine
	var err error
	machines, err = models.GetAllMachines()
	if err != nil {
		beego.Error("Failed to get all machines", err)
		this.CustomAbort(403, "Failed to get all machines")
	}

	this.Data["json"] = machines
	this.ServeJson()
}

// @Title Get
// @Description Get machine by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 {object} models.Machine
// @Failure	403	Failed to get machine
// @Failure	401	Not authorized
// @router /:mid [get]
func (this *MachinesController) Get() {

	var machineId int64
	var err error
	machineId, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(403, "Failed to get machine")
	}

	if !this.IsAdmin() && !this.IsStaff() {

		// Get user permissions to see whether user is allowed to access machine
		sessUserId, ok := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
		if !ok {
			beego.Error("Failed to get session user ID")
			this.CustomAbort(403, "Failed to get machine")
		}

		var permissions []*models.Permission
		permissions, err = models.GetUserPermissions(sessUserId)
		if err != nil {
			beego.Error("Failed to get machine permissions", err)
			this.CustomAbort(401, "Not authorized")
		}

		permissionFound := false
		for _, value := range permissions {
			if value.MachineId == machineId {
				permissionFound = true
				break
			}
		}

		if !permissionFound {
			beego.Error("User not authorized to view this machine")
			this.CustomAbort(401, "Not authorized")
		}
	}

	var machine *models.Machine
	machine, err = models.GetMachine(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.CustomAbort(403, "Failed to get machine")
	}

	this.Data["json"] = machine
	this.ServeJson()
}

// @Title Create
// @Description Create machine
// @Param	mname	query	string	true	"Machine Name"
// @Success 200 {object} models.MachineCreatedResponse
// @Failure	403	Failed to create machine
// @Failure	401	Not authorized
// @router / [post]
func (this *MachinesController) Create() {
	machineName := this.GetString("mname")
	beego.Trace(machineName)

	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized to create machine")
		this.CustomAbort(401, "Not authorized")
	}

	// All clear - create machine in the database
	var machineId int64
	var err error
	machineId, err = models.CreateMachine(machineName)
	if err != nil {
		beego.Error("Failed to create machine", err)
		this.CustomAbort(403, "Failed to create machine")
	}

	// Success - we even got an ID!!!
	this.Data["json"] = &models.MachineCreatedResponse{MachineId: machineId}
	this.ServeJson()
}

// @Title Update
// @Description Update machine
// @Param	mid	path	int	true	"Machine ID"
// @Param	model	body	models.Machine	true	"Machine model"
// @Success 200 string ok
// @Failure	403	Failed to update machine
// @Failure	401	Not authorized
// @router /:mid [put]
func (this *MachinesController) Update() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	// expecting the following JSON
	/*
		{
			"Available": true,
			"Description": "My machine description",
			"Id": 1,
			"Image": "/image",
			"Name": "My Machine Name",
			"Price": 2.35,
			"PriceUnit": "minute",
			"Shortname": "MMN",
			"UnavailMsg": "Can be empty",
			"UnavailTill": "2002-10-02T10:00:00-05:00"
		}
	*/
	// The UnavailTill field must be RFC 3339 formatted
	// https://golang.org/src/time/format.go

	var err error

	// Attempt to decode passed json
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.Machine{}
	if err = dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update machine")
	}

	var mid int64

	// Get mid and check if it matches with the machine model ID
	mid, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get :mid:", err)
		this.CustomAbort(403, "Failed to update machine")
	}
	if mid != req.Id {
		beego.Error("mid and model ID does not match:", err)
		this.CustomAbort(403, "Failed to update machine")
	}

	// Update the database
	err = models.UpdateMachine(&req)
	if err != nil {
		beego.Error("Failed updating machine:", err)
		this.CustomAbort(403, "Failed to update machine")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title Delete
// @Description Delete machine
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	403	Failed to delete machine
// @Failure	401	Not authorized
// @router /:mid [delete]
func (this *MachinesController) Delete() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error
	var mid int64

	mid, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get mid:", err)
		this.CustomAbort(403, "Failed to delete machine")
	}

	err = models.DeleteMachine(mid)
	if err != nil {
		beego.Error("Failed to delete machine:", err)
		this.CustomAbort(403, "Failed to delete machine")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
