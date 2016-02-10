// /api/machines
package machines

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/astaxie/beego"
	"strings"
)

type Controller struct {
	controllers.Controller
}

// GetLocIdAdmin gets the location id, if passed as URL parameter, otherwise
// it will be 0.  0 being synonym for all locations.  Also it returns whether
// the user is allowed to perform admin tasks at that location.
func (this *Controller) GetLocIdAdmin() (locId int64, authorized bool) {
	locId, err := this.GetInt64("location")
	if err == nil {
		if !this.IsAdminAt(locId) {
			return
		}
	} else {
		locId = 0
		if !this.IsSuperAdmin() {
			return
		}
	}
	return locId, true
}

// @Title GetAll
// @Description Get all machines (limited to location, if specified)
// @Param	location	query	int64	false	"Location ID"
// @Success 200 {object} models.Machine
// @Failure	401 Not authorized
// @Failure	500	Failed to get all machines
// @router / [get]
func (this *Controller) GetAll() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	allMachines, err := models.GetAllMachines()
	if err != nil {
		beego.Error("Failed to get all machines", err)
		this.CustomAbort(500, "Failed to get all machines")
	}

	machines := make([]*models.Machine, 0, len(allMachines))
	for _, m := range allMachines {
		if locId == 0 || locId == m.LocationId {
			machines = append(machines, m)
		}
	}

	this.Data["json"] = machines
	this.ServeJSON()
}

// @Title Get
// @Description Get machine by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 {object} models.Machine
// @Failure	403	Failed to get machine
// @Failure	401	Not authorized
// @router /:mid [get]
func (this *Controller) Get() {

	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(403, "Failed to get machine")
	}

	if !this.IsAdmin() && !this.IsStaff() {

		// Get user permissions to see whether user is allowed to access machine
		sessUserId, err := this.GetSessionUserId()
		if err != nil {
			beego.Error("Failed to get session user ID:", err)
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
	this.ServeJSON()
}

// @Title Create
// @Description Create machine
// @Param	mname	query	string	true	"Machine Name"
// @Success 200 {object} models.MachineCreatedResponse
// @Failure	403	Failed to create machine
// @Failure	401	Not authorized
// @router / [post]
func (this *Controller) Create() {
	machineName := this.GetString("mname")

	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized to create machine")
		this.CustomAbort(401, "Not authorized")
	}

	machineId, err := models.CreateMachine(machineName)
	if err != nil {
		beego.Error("Failed to create machine", err)
		this.CustomAbort(403, "Failed to create machine")
	}

	this.Data["json"] = &models.MachineCreatedResponse{MachineId: machineId}
	this.ServeJSON()
}

// @Title Update
// @Description Update machine
// @Param	mid	path	int	true	"Machine ID"
// @Param	model	body	models.Machine	true	"Machine model"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500	Failed to update machine
// @router /:mid [put]
func (this *Controller) Update() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.Machine{}
	if err := dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(400, "Failed to update machine")
	}

	// Get mid and check if it matches with the machine model ID
	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get :mid:", err)
		this.CustomAbort(400, "Failed to update machine")
	}
	if mid != req.Id {
		beego.Error("mid and model ID does not match:", err)
		this.CustomAbort(400, "Failed to update machine")
	}

	if err = req.Update(); err != nil {
		beego.Error("Failed updating machine:", err)
		if err == models.ErrDimensions || err == models.ErrWorkspaceDimensions {
			this.CustomAbort(400, err.Error())
		} else {
			this.CustomAbort(500, "Failed to update machine")
		}

	}

	this.Data["json"] = "ok"
	this.ServeJSON()
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
func (this *Controller) PostImage() {
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

func (this *Controller) underMaintenanceOnOrOff(onOrOff int) error {
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
func (this *Controller) ReportBroken() {
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

	uid, err := this.GetSessionUserId()
	if err != nil {
		beego.Info("Not logged in:", err)
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
	this.ServeJSON()
}

// @Title UnderMaintenanceOn
// @Description Turn On Under Maintenance
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/under_maintenance/on [post]
func (this *Controller) UnderMaintenanceOn() {
	err := this.underMaintenanceOnOrOff(ON)
	if err != nil {
		beego.Error("Failed to set UnderMaintenance: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title UnderMaintenanceOff
// @Description Turn Off Under Maintenance
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/under_maintenance/off [post]
func (this *Controller) UnderMaintenanceOff() {
	err := this.underMaintenanceOnOrOff(OFF)
	if err != nil {
		beego.Error("Failed to unset UnderMaintenance: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Data["json"] = "ok"
	this.ServeJSON()
}

func (this *Controller) switchMachine(onOrOff int) error {
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
func (this *Controller) TurnOn() {
	if err := this.switchMachine(ON); err != nil {
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title TurnOff
// @Description Turn On Machine Switch
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/turn_off [post]
func (this *Controller) TurnOff() {
	if err := this.switchMachine(OFF); err != nil {
		this.CustomAbort(500, "Internal Server Error")
	}
	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title Search
// @Description Search machine
// @Param	machine_type_id	query	int64	true	"Machine Type Id"
// @Success 200
// @Failure	400	Client Error
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /search [post]
func (this *Controller) Search() {
	machineTypeId, err := this.GetInt64("machine_type_id")
	if err != nil {
		beego.Error("machine_type_id:", err)
		this.CustomAbort(400, "Client error")
	}

	ms, err := models.GetAllMachines()
	if err != nil {
		beego.Error("get all machines", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	results := make([]*models.Machine, 0, len(ms))

	for _, m := range ms {
		if m.TypeId == machineTypeId {
			results = append(results, m)
		}
	}

	this.Data["json"] = results
	this.ServeJSON()
}
