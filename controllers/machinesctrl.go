package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
	"strings"
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

		var permissions *[]models.Permission
		permissions, err = models.GetUserPermissions(sessUserId)
		if err != nil {
			beego.Error("Failed to get machine permissions", err)
			this.CustomAbort(401, "Not authorized")
		}

		permissionFound := false
		for _, value := range *permissions {
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

// @Title GetConnections
// @Description Get connected machines
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 {object} models.ConnectedMachineList
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:mid/connections [get]
func (this *MachinesController) GetConnections() {

	var machineId int64
	var err error
	machineId, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(500, "Internal Server Error")
	}

	machineList, err := models.GetConnectedMachines(machineId)
	if err != nil {
		beego.Error("Failed to get connected machines")
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = machineList
	this.ServeJson()
}

// @Title GetConnected
// @Description Get connectable machines
// @Param	mid		path 	int	true		Machine ID
// @Success 200 {object} models.ConnectableMachineList
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:mid/connectable [get]
func (this *MachinesController) GetConnectable() {

	var machineId int64
	var err error
	machineId, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var cmList *models.ConnectableMachineList
	cmList, err = models.GetConnectableMachines(machineId)
	if err != nil {
		beego.Error("Could not get connectable machine list:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = cmList
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

func decodeDataUri(dataUri string) ([]byte, error) {
	sep := "base64,"
	i := strings.Index(dataUri, sep)
	if i < 0 {
		return nil, fmt.Errorf("cannot remove prefix from data uri")
	}
	dataUri = dataUri[i+len(sep):]
	data, err := base64.StdEncoding.DecodeString(dataUri)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// @Title PostImage
// @Description Post machine image
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/image [post]
func (this *MachinesController) PostImage() {
	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized to create machine")
		this.CustomAbort(401, "Not authorized")
	}

	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get mid:", err)
		this.CustomAbort(403, "Failed to delete machine")
	}
	img, err := decodeDataUri(this.GetString("Image"))
	if err != nil {
		beego.Error("decode data uri:", err)
		this.CustomAbort(400, "Bad Request")
	}
	i := strings.LastIndex(this.GetString("Filename"), ".")
	var fileExt string
	if i >= 0 {
		fileExt = this.GetString("Filename")[i:]
	} else {
		this.CustomAbort(500, "File name has no proper extension")
	}
	fn := imageFilename(this.GetString(":mid"), fileExt)
	if err = models.UploadImage(fn, img); err != nil {
		this.CustomAbort(500, "Internal Server Error")
	}

	m, err := models.GetMachine(mid)
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(403, "Failed to get machine")
	}
	m.Image = fn
	if err = models.UpdateMachine(m); err != nil {
		beego.Error("Failed updating machine:", err)
		this.CustomAbort(403, "Failed to update machine")
	}
}

func imageFilename(machineId string, fileExt string) string {
	return "machine-" + machineId + fileExt
}

const (
	ON = iota
	OFF
)

func (this *MachinesController) underMaintenanceOnOrOff(onOrOff int) error {
	var machineId int64
	var err error
	machineId, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var machine *models.Machine
	machine, err = models.GetMachine(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.CustomAbort(403, "Failed to get machine")
	}

	return machine.SetUnderMaintenance(onOrOff == ON)
}

// @Title ReportBroken
// @Description Report machine as broken
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/report_broken [post]
func (this *MachinesController) ReportBroken() {
	var machineId int64
	var err error
	machineId, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	var machine *models.Machine
	machine, err = models.GetMachine(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.CustomAbort(403, "Failed to get machine")
	}

	beego.Info("Now calling machine#ReportBroken:")
	err = machine.ReportBroken()
	if err != nil {
		beego.Error("Failed to report machine", err)
		this.CustomAbort(500, "Failed to report machine")
	}
	beego.Info("machine#ReportBroken called!")
	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title UnderMaintenanceOn
// @Description Turn On Under Maintenance
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/under_maintenance/on [post]
func (this *MachinesController) UnderMaintenanceOn() {
	err := this.underMaintenanceOnOrOff(ON)
	if err != nil {
		beego.Error("Failed to set UnderMaintenance: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title UnderMaintenanceOff
// @Description Turn Off Under Maintenance
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/under_maintenance/off [post]
func (this *MachinesController) UnderMaintenanceOff() {
	err := this.underMaintenanceOnOrOff(OFF)
	if err != nil {
		beego.Error("Failed to unset UnderMaintenance: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Data["json"] = "ok"
	this.ServeJson()
}

func (this *MachinesController) switchMachine(onOrOff int) error {
	var machineId int64
	var err error
	machineId, err = this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var machine *models.Machine
	machine, err = models.GetMachine(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.CustomAbort(403, "Failed to get machine")
	}

	if onOrOff == ON {
		return machine.On()
	} else {
		return machine.Off()
	}

	return nil
}

// @Title TurnOn
// @Description Turn On Machine Switch
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/turn_on [post]
func (this *MachinesController) TurnOn() {
	err := this.switchMachine(ON)
	if err != nil {
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title TurnOff
// @Description Turn On Machine Switch
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/turn_off [post]
func (this *MachinesController) TurnOff() {
	err := this.switchMachine(OFF)
	if err != nil {
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Data["json"] = "ok"
	this.ServeJson()
}
