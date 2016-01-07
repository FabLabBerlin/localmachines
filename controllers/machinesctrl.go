package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/astaxie/beego"
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

	machines, err := models.GetAllMachines(true)
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

	machineId, err := this.GetInt64(":mid")
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

		permissions, err := models.GetUserPermissions(sessUserId)
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

	machine, err := models.GetMachine(machineId)
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

	machineId, err := this.GetInt64(":mid")
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

	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	cmList, err := models.GetConnectableMachines(machineId)
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
	machineId, err := models.CreateMachine(machineName)
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

	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.Machine{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update machine")
	}

	// Get mid and check if it matches with the machine model ID
	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get :mid:", err)
		this.CustomAbort(403, "Failed to update machine")
	}
	if mid != req.Id {
		beego.Error("mid and model ID does not match:", err)
		this.CustomAbort(403, "Failed to update machine")
	}

	if err = req.Update(); err != nil {
		beego.Error("Failed updating machine:", err)
		this.CustomAbort(403, "Failed to update machine")
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
	if err = m.Update(); err != nil {
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
	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	machine, err := models.GetMachine(machineId)
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
	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	machine, err := models.GetMachine(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.CustomAbort(403, "Failed to get machine")
	}

	uid, ok := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	if !ok {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}
	user, err := models.GetUser(uid)
	if err != nil {
		beego.Error("Failed to get machine permissions", err)
		this.CustomAbort(401, "Not authorized")
	}

	err = machine.ReportBroken(*user)
	if err != nil {
		beego.Error("Failed to report machine", err)
		this.CustomAbort(500, "Failed to report machine")
	}
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
	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	machine, err := models.GetMachine(machineId)
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
	if err := this.switchMachine(ON); err != nil {
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
	if err := this.switchMachine(OFF); err != nil {
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Data["json"] = "ok"
	this.ServeJson()
}
