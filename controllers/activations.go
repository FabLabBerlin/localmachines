package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
	"github.com/kr15h/fabsmith/plugins/hexaswitch"
	"time"
)

type ActivationsController struct {
	Controller
}

// @Title Get
// @Description Get activation by activation ID
// @Param	aid		path 	int	true		"Activation ID"
// @Success 200 {object} models.Activation
// @Failure	403	Failed to get activation
// @Failure	401	Not authorized
// @router /:aid [get]
func (this *ActivationsController) Get() {

}

// @Title GetAll
// @Description Get all active activations
// @Success 200 {object} models.Activation
// @Failure	403	Failed to get active activations
// @router /active [get]
func (this *ActivationsController) GetActive() {
	activations, err := models.GetActiveActivations()
	if err != nil {
		beego.Error("Failed to get active activations")
		this.CustomAbort(403, "Failed to get active activations")
	}
	this.Data["json"] = activations
	this.ServeJson()
}

// @Title Create
// @Description Create new activation
// @Param	mid		query 	int	true		"Machine ID"
// @Success 201 {object} models.ActivationCreateResponse
// @Failure	403	Failed to create activation
// @Failure 401 Not authorized
// @router / [post]
func (this *ActivationsController) Create() {
	var machineId int64
	var err error
	machineId, err = this.GetInt64("mid")
	userId := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)

	// Check if user has permissions to create activation for the machine.
	var userPermissions []*models.Permission
	userPermissions, err = models.GetUserPermissions(userId)
	if err != nil {
		beego.Error("Could not get user permissions")
		this.CustomAbort(403, "Forbidden")
	}
	var userPermitted = false
	for _, permission := range userPermissions {
		if int64(permission.MachineId) == machineId {
			userPermitted = true
			break
		}
	}
	if !userPermitted {
		beego.Error("User has no permission to activate the machine")
		this.CustomAbort(401, "Unauthorized")
	}

	// Continue with creating activation
	var activationId int64
	activationId, err = models.CreateActivation(machineId, userId)
	if err != nil {
		beego.Error("Failed to create activation")
		this.CustomAbort(403, "Failed to create activation")
	}

	// Check if there is mapping between switch and machine
	var switchMappingExists bool
	_, err = hexaswitch.GetSwitchIp(int(machineId))
	if err != nil {
		beego.Warning("Machine / switch mapping does not exist")
		switchMappingExists = false
	} else {
		switchMappingExists = true
	}

	if switchMappingExists {
		hexaswitch.Install() // TODO: remove this from here in an elegant way
		// TODO: use switch IP for the SwitchOn method
		err = hexaswitch.SwitchOn(int(machineId))
		if err != nil {
			beego.Error("Failed to turn on the switch")
			err = models.DeleteActivation(activationId)
			if err != nil {
				beego.Error("Failed to delete activation")
			}
			this.CustomAbort(403, "Failed to create activation")
		}
	}

	response := models.ActivationCreateResponse{}
	response.ActivationId = activationId
	this.Data["json"] = response
	this.ServeJson()
}

// @Title Close
// @Description Close running activation
// @Param	aid		path 	int	true		"Activation ID"
// @Success 200 {object} models.StatusResponse
// @Failure	403	Failed to close activation
// @Failure 401 Not authorized
// @router /:aid [put]
func (this *ActivationsController) Close() {
	aid, err := this.GetInt64(":aid")
	if err != nil {
		beego.Error("Failed to get :aid")
		this.CustomAbort(403, "Failed to close activation")
	}
	err = models.CloseActivation(aid)
	if err != nil {
		beego.Error("Failed to close activation")
		this.CustomAbort(403, "Failed to close activation")
	}

	// Launch go routine to switch off the machine after some time
	var deactivateTimeout int64
	deactivateTimeout, err = beego.AppConfig.Int64("deactivatetimeout")
	if err != nil {
		beego.Error("Failed to get deactivate timeout from config:", err)
		deactivateTimeout = 30
	}

	var machineId int64
	machineId, err = models.GetActivationMachineId(aid)
	if err != nil {
		beego.Error("Failed to get machine ID")
		this.CustomAbort(403, "Failed to close activation")
	}

	go deactivateMachineAfterTimeout(machineId, deactivateTimeout)

	statusResponse := models.StatusResponse{}
	statusResponse.Status = "ok"
	this.Data["json"] = statusResponse
	this.ServeJson()
}

// Deativates a machine after timeout if no activation with the machine ID
// has been made.
func deactivateMachineAfterTimeout(machineId int64, timeoutSeconds int64) {

	// Timeout
	time.Sleep(time.Duration(timeoutSeconds) * time.Second)

	// Check if activation with the machine ID exists
	o := orm.NewOrm()
	activationModel := models.Activation{Id: 0}
	beego.Info("Attempt to get an active activation with the machine ID", machineId)
	err := o.Raw("SELECT id FROM activation WHERE active=true AND machine_id=?",
		machineId).QueryRow(&activationModel)
	if err != nil {
		beego.Error("There was an error while getting activation:", err)
		// Here we assume that there is no activation
		// and switch off the machine anyway
	}
	beego.Trace("Got activation ID", activationModel.Id)
	if activationModel.Id > 0 {
		beego.Info("There is an activation for the machine, keep it on")
		return
	}

	// Attempt to switch off the machine
	hexaswitch.Install()                       // TODO: remove this from here in an elegant way
	err = hexaswitch.SwitchOff(int(machineId)) // This will take some time
	if err != nil {
		beego.Error("Failed to turn the switch off", err)
	}
}
