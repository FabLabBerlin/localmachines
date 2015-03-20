package controllers

import (
	"github.com/kr15h/fabsmith/models"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/plugins/hexaswitch"
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
	var machineId int
	var err error
	machineId, err = this.GetInt("mid")
	userId := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int)
	
	var activationId int
	activationId, err = models.CreateActivation(machineId, userId)
	if err != nil {
		beego.Error("Failed to create activation")
		this.CustomAbort(403, "Failed to create activation")
	}

	// Check if there is mapping between switch and machine
	var switchMappingExists bool
	_, err = hexaswitch.GetSwitchIp(machineId)
	if err != nil {
		beego.Warning("Machine / switch mapping does not exist")
		switchMappingExists = false
	} else {
		switchMappingExists = true
	}

	if switchMappingExists {
		hexaswitch.Install() // TODO: remove this from here in an elegant way
		// TODO: use switch IP for the SwitchOn method
		err = hexaswitch.SwitchOn(machineId)
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

}