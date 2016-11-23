// /api/machines
package machines

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/metrics/machine_earnings"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"os/exec"
	"strings"
	"time"
)

type Controller struct {
	controllers.Controller
}

// @Title GetAll
// @Description Get all machines (limited to location, if specified)
// @Param	location	query	int64	false	"Location ID"
// @Param	archived	query	bool	false	"Include archived items"
// @Success 200 machine.Machine
// @Failure	401 Not authorized
// @Failure	500	Failed to get all machines
// @router / [get]
func (this *Controller) GetAll() {
	locId, authorized := this.GetLocIdApi()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	includeArchived := true

	if a, err := this.GetBool("archived"); err != nil {
		if !a && err == nil {
			includeArchived = false
		}
	}

	allMachines, err := machine.GetAll()
	if err != nil {
		beego.Error("Failed to get all machines", err)
		this.CustomAbort(500, "Failed to get all machines")
	}

	machines := make([]*machine.Machine, 0, len(allMachines))
	for _, m := range allMachines {
		if !includeArchived {
			continue
		}

		if m.NetswitchLastPing.Unix() > 0 {
			if time.Since(m.NetswitchLastPing).Minutes() > 30 {
				m.Status = "Offline"
			} else {
				m.Status = "Online"
			}
		}
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
// @Failure	400	Wrong Input
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:mid [get]
func (this *Controller) Get() {

	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.Abort("400")
	}

	machine, err := machine.Get(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.Abort("500")
	}

	if !this.IsStaffAt(machine.LocationId) {

		// Get user permissions to see whether user is allowed to access machine
		sessUserId, err := this.GetSessionUserId()
		if err != nil {
			beego.Error("Failed to get session user ID:", err)
			this.Abort("400")
		}

		permissions, err := user_permissions.Get(sessUserId)
		if err != nil {
			beego.Error("Failed to get machine permissions", err)
			this.Abort("401")
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
			this.Abort("401")
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
// @Failure	500	Failed to create machine
// @Failure	401	Not authorized
// @router / [post]
func (this *Controller) Create() {
	name := this.GetString("mname")

	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.Abort("401")
	}

	m, err := machine.Create(locId, name)
	if err != nil {
		beego.Error("Failed to create machine", err)
		this.Abort("500")
	}

	this.Data["json"] = &models.MachineCreatedResponse{MachineId: m.Id}
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

	updated.NetswitchSensorPort = 1

	updateGateway := existing.NetswitchConfigured() ||
		updated.NetswitchConfigured()

	if err = updated.Update(updateGateway); err != nil {
		beego.Error("Failed updating machine:", err)
		if err == machine.ErrDimensions ||
			err == machine.ErrWorkspaceDimensions ||
			err == machine.ErrDuplicateNetswitchHost {
			this.CustomAbort(400, err.Error())
		} else {
			this.CustomAbort(500, "Failed to update machine")
		}

	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title GetEarnings
// @Description Get earnings by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 machine.Machine
// @Failure	400	Wrong Input
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:mid/earnings [get]
func (this *Controller) GetEarnings() {
	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.Abort("400")
	}

	machine, err := machine.Get(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.Abort("500")
	}

	invs, err := invutil.GetAllAt(machine.LocationId)
	if err != nil {
		this.Fail(500, "Failed to get invoices")
	}

	if !this.IsAdminAt(machine.LocationId) {
		this.Fail(403)
	}

	me := machine_earnings.New(
		machine,
		month.New(1, 2015),
		month.New(12, 2017),
		invs,
	)

	resp := make(map[string]interface{})

	resp["Memberships"] = me.MembershipsCached()
	resp["PayAsYouGo"] = me.PayAsYouGoCached()
	this.Data["json"] = resp
	this.ServeJSON()
}

// @Title SetArchived
// @Description (Un)archive machine
// @Param	mid		path	int		true	"Machine ID"
// @Param	archive	query	bool	true	"Archive"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500	Failed to archive machine
// @router /:mid/set_archived [post]
func (this *Controller) SetArchived() {
	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.Abort("400")
	}

	machine, err := machine.Get(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(machine.LocationId) {
		beego.Error("Not authorized")
		this.Abort("401")
	}

	machine.Archived, err = this.GetBool("archived")
	if err != nil {
		beego.Error("parsing archived parameter")
		this.Abort("400")
	}

	if err = machine.Update(true); err != nil {
		beego.Error("update:", err)
		this.Abort("500")
	}

	this.Finish()
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
	dataUri := this.GetString("Image")

	i := strings.LastIndex(this.GetString("Filename"), ".")
	var fileExt string
	if i >= 0 {
		fileExt = this.GetString("Filename")[i:]
	} else {
		this.CustomAbort(500, "File name has no proper extension")
	}
	fn, fnSmall := imageFilename(this.GetString(":mid"), fileExt)
	if err = models.UploadImage(fn, dataUri); err != nil {
		this.CustomAbort(500, "Internal Server Error")
	}

	m, err := machine.Get(mid)
	if err != nil {
		beego.Error("Failed to get machine:", err)
		this.CustomAbort(500, "Failed to get machine")
	}

	if ext := strings.ToLower(fileExt); ext != ".svg" {
		args := make([]string, 0, 5)

		args = append(args, imageMagickCompressArgs(ext)...)

		args = append(args,
			"-resize", "x400",
			"files/"+fn,
			"files/"+fnSmall,
		)
		cmd := exec.Command("convert", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			beego.Error("convert:", err, string(out))
			this.CustomAbort(500, "Internal Server Error")
		}
		m.ImageSmall = fnSmall
	} else {
		m.ImageSmall = fn
	}

	m.Image = fn
	if ext := strings.ToLower(fileExt); ext != ".svg" {
		args := make([]string, 0, 5)
		args = append(args, imageMagickCompressArgs(ext)...)
		args = append(args,
			"-resize", "2000x2000>",
			"files/"+fn,
		)
		cmd := exec.Command("mogrify", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			beego.Error("mogrify:", err, string(out), "(skip extra compression)")
		}
	}

	if !this.IsAdminAt(m.LocationId) {
		beego.Error("Not authorized to set machine image")
		this.CustomAbort(401, "Not authorized")
	}

	if err = m.Update(false); err != nil {
		beego.Error("Failed updating machine:", err)
		this.CustomAbort(500, "Failed to update machine")
	}
}

func imageFilename(machineId string, fileExt string) (normal, small string) {
	normal = "machine-" + machineId + fileExt
	small = "machine-" + machineId + "-small" + fileExt
	return
}

func imageMagickCompressArgs(fileExt string) (args []string) {
	args = make([]string, 0, 5)
	ext := strings.ToLower(fileExt)
	switch ext {
	case ".jpg":
		args = append(args, "-quality", "75%")
		// Browsers have problems with EXIF data based rotations,
		// let's get rid of them.
		args = append(args, "-auto-orient")
		break
	case ".png":
		args = append(args,
			"-define",
			"png:compression-filter=5",
			"-define",
			"png:compression-level=9",
			"-define",
			"png:compression-strategy=1",
		)
		break
	}
	return
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
// @Param	mid		path	int		true	"Machine ID"
// @Param	text	body	string	true	"Text"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:mid/report_broken [post]
func (this *Controller) ReportBroken() {
	beego.Info("/api/machines/:id/report_broken")
	machineId, err := this.GetInt64(":mid")
	if err != nil {
		beego.Error("Failed to get :mid variable")
		this.Abort("500")
	}

	text := this.GetString("text")

	m, err := machine.Get(machineId)
	if err != nil {
		beego.Error("Failed to get machine", err)
		this.Abort("500")
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

	err = m.ReportBroken(*user, text)
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

	uid, err := this.GetSessionUserId()
	if err != nil {
		beego.Info("Not logged in:", err)
		this.CustomAbort(401, "Not logged in")
	}

	if err := m.NetswitchApplyConfig(uid); err != nil {
		beego.Error("Error configuring:", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	this.ServeJSON()
}
