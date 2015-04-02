package controllers

import (
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

	var err error

	uid := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)

	// Check if user is admin or staff
	userRoles, err := models.GetUserRoles(uid)
	if !userRoles.Admin && !userRoles.Staff {
		beego.Error("Not authorized to create machine")
		this.CustomAbort(401, "Not authorized")
	}

	// All clear - create machine in the database
	var machineId int64
	machineId, err = models.CreateMachine(machineName)
	if err != nil {
		beego.Error("Failed to create machine", err)
		this.CustomAbort(403, "Failed to create machine")
	}

	// Success - we even got an ID!!!
	this.Data["json"] = &models.MachineCreatedResponse{MachineId: machineId}
	this.ServeJson()
}
