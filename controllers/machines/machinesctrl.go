// /api/machines
package machines

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"strings"
)

type Controller struct {
	controllers.Controller
}

// @Title GetAll
// @Description Get all machines (limited to location, if specified)
// @Param	location	query	int64	false	"Location ID"
// @Success 200 machine.Machine
// @Failure	401 Not authorized
// @Failure	500	Failed to get all machines
// @router / [get]
func (this *Controller) GetAll() {
	locId, authorized := this.GetLocIdApi()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	allMachines, err := machine.GetAll()
	if err != nil {
		beego.Error("Failed to get all machines", err)
		this.CustomAbort(500, "Failed to get all machines")
	}

	machines := make([]*machine.Machine, 0, len(allMachines))
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
// @Success 200 machine.Machine
// @Failure	403	Failed to get machine
// @Failure	401	Not authorized
// @router /:mid [get]
func (this *Controller) Get() {

	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(403, "Failed to get machine")
	}

	machine, err := machine.Get(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.CustomAbort(403, "Failed to get machine")
	}

	if !this.IsStaffAt(machine.LocationId) {

		// Get user permissions to see whether user is allowed to access machine
		sessUserId, err := this.GetSessionUserId()
		if err != nil {
			beego.Error("Failed to get session user ID:", err)
			this.CustomAbort(403, "Failed to get machine")
		}

		permissions, err := user_permissions.Get(sessUserId)
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

	this.Data["json"] = machine
	this.ServeJSON()
}

// @Title Create
// @Description Create machine
// @Param	location	query	string	true	"Location Id"
// @Param	mname		query	string	true	"Machine Name"
// @Success 200 machine.MachineCreatedResponse
// @Failure	403	Failed to create machine
// @Failure	401	Not authorized
// @router / [post]
func (this *Controller) Create() {
	name := this.GetString("mname")

	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	machineId, err := machine.Create(locId, name)
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
// @Param	model	body	machine.Machine	true	"Machine model"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500	Failed to update machine
// @router /:mid [put]
func (this *Controller) Update() {
	// Get mid and check if it matches with the machine model ID
	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Could not get :mid:", err)
		this.Abort("400")
	}

	existing, err := machine.Get(mid)
	if err != nil {
		beego.Error("get machine:", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(existing.LocationId) {
		beego.Error("Not authorized")
		this.Abort("401")
	}

	dec := json.NewDecoder(this.Ctx.Request.Body)
	updated := machine.Machine{}
	if err := dec.Decode(&updated); err != nil {
		beego.Error("Failed to decode json:", err)
		this.Abort("400")
	}

	if mid != updated.Id || existing.LocationId != updated.LocationId {
		beego.Error("mid and model (location) ID does not match:", err)
		this.Abort("403")
	}

	updateGateway := existing.NetswitchConfigured() ||
		updated.NetswitchConfigured()

	if err = updated.Update(updateGateway); err != nil {
		beego.Error("Failed updating machine:", err)
		if err == machine.ErrDimensions || err == machine.ErrWorkspaceDimensions {
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
	return base64.StdEncoding.DecodeString(dataUri)
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

	m, err := machine.Get(mid)
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(403, "Failed to get machine")
	}

	if !this.IsAdminAt(m.LocationId) {
		beego.Error("Not authorized to set machine image")
		this.CustomAbort(401, "Not authorized")
	}

	m.Image = fn
	if err = m.Update(false); err != nil {
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
	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(500, "Internal Server Error")
	}

	m, err := machine.Get(mid)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.CustomAbort(403, "Failed to get machine")
	}

	if !this.IsStaffAt(m.LocationId) {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	return m.SetUnderMaintenance(onOrOff == ON)
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

	m, err := machine.Get(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.CustomAbort(403, "Failed to get machine")
	}

	uid, err := this.GetSessionUserId()
	if err != nil {
		beego.Info("Not logged in:", err)
		this.CustomAbort(401, "Not logged in")
	}
	user, err := users.GetUser(uid)
	if err != nil {
		beego.Error("Failed to get machine permissions", err)
		this.CustomAbort(401, "Not authorized")
	}

	err = m.ReportBroken(*user)
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

// @Title Search
// @Description Search machine globally
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

	ms, err := machine.GetAll()
	if err != nil {
		beego.Error("get all machines", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	results := make([]*machine.Machine, 0, len(ms))

	for _, m := range ms {
		if m.TypeId == machineTypeId {
			m.NetswitchUrlOn = ""
			m.NetswitchUrlOff = ""
			m.NetswitchHost = ""
			m.NetswitchSensorPort = 0
			results = append(results, m)
		}
	}

	this.Data["json"] = results
	this.ServeJSON()
}

// @Title ApplyConfig
// @Description Apply custom switch config
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/apply_config [post]
func (this *Controller) ApplyConfig() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	mid, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(403, "Failed to get machine")
	}

	m, err := machine.Get(mid)
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.CustomAbort(403, "Failed to get machine")
	}

	if m.LocationId != locId {
		beego.Error("Wrong location id")
		this.CustomAbort(401, "Not authorized")
	}

	if err := m.NetswitchApplyConfig(); err != nil {
		beego.Error("Error configuring:", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	this.ServeJSON()
}
