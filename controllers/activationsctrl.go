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

// @Title Get All
// @Description Get all activations
// @Param	startDate		query 	string	true		"Period start date"
// @Param	endDate		query 	string	true		"Period end date"
// @Param	userId		query 	int	true		"User ID"
// @Param	includeInvoiced		query 	bool	true		"Whether to include already invoiced activations"
// @Param	itemsPerPage		query 	int	true		"Items per page or max number of items to return"
// @Param	page		query 	int	true		"Current page to show"
// @Success 200 {object} models.GetActivationsResponse
// @Failure	403	Failed to get activations
// @Failure	401	Not authorized
// @router / [get]
func (this *ActivationsController) GetAll() {

	// Only admin can use this API call
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error
	var startDate string
	var endDate string
	var userId int64
	var includeInvoiced bool
	var itemsPerPage int64
	var page int64

	// Get variables
	startDate = this.GetString("startDate")
	if startDate == "" {
		beego.Error("Missing start date")
		this.CustomAbort(403, "Failed to get activations")
	}

	endDate = this.GetString("endDate")
	if endDate == "" {
		beego.Error("Missing end date")
		this.CustomAbort(403, "Failed to get activations")
	}

	userId, err = this.GetInt64("userId")
	if err != nil {
		beego.Error("Could not get userId request variable:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	includeInvoiced, err = this.GetBool("includeInvoiced")
	if err != nil {
		beego.Error("Could not get includeInvoiced request variable:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	itemsPerPage, err = this.GetInt64("itemsPerPage")
	if err != nil {
		beego.Error("Could not get itemsPerPage request variable:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	page, err = this.GetInt64("page")
	if err != nil {
		beego.Error("Could not get page request variable:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	beego.Trace("startDate:", startDate)
	beego.Trace("endDate:", endDate)
	beego.Trace("userId:", userId)
	beego.Trace("includeInvoiced:", includeInvoiced)
	beego.Trace("itemsPerPage:", itemsPerPage)
	beego.Trace("page:", page)

	// Convert / parse string time values as time.Time
	var timeForm = "2006-01-02"
	var startTime time.Time

	startTime, err = time.ParseInLocation(
		timeForm, startDate, time.Now().Location())
	if err != nil {
		beego.Error("Failed to parse startDate:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	var endTime time.Time

	endTime, err = time.ParseInLocation(
		timeForm, endDate, time.Now().Location())
	if err != nil {
		beego.Error("Failed to parse endDate:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	beego.Trace(startTime)
	beego.Trace(endTime)

	// Get activations
	var activations *[]models.Activation
	activations, err = models.GetActivations(
		startTime, endTime, userId, includeInvoiced, itemsPerPage, page)
	if err != nil {
		beego.Error("Failed to get activations:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	// Get total activation count
	var numActivations int64
	numActivations, err = models.GetNumActivations(
		startTime, endTime, userId, includeInvoiced)
	if err != nil {
		beego.Error("Failed to get number of activations:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	r := models.GetActivationsResponse{}
	r.NumActivations = numActivations
	r.ActivationsPage = activations

	this.Data["json"] = r
	this.ServeJson()
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

// @Title Get Active
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

	userId, ok := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	if ok == false {
		beego.Error("Failed to get session user ID")
		this.CustomAbort(403, "Failed to create activation")
	}

	machineId, err = this.GetInt64("mid")
	if err != nil {
		beego.Error("Failed to get mid:", err)
		this.CustomAbort(403, "Failed to create activation")
	}

	// Admins can activate any machine (except broken ones).
	// Regular users have to refer to their permissions.
	if !this.IsAdmin() {

		// Check if user has permissions to create activation for the machine.
		var userPermissions *[]models.Permission
		userPermissions, err = models.GetUserPermissions(userId)
		if err != nil {
			beego.Error("Could not get user permissions")
			this.CustomAbort(403, "Failed to create activation")
		}
		var userPermitted = false
		for _, permission := range *userPermissions {
			if int64(permission.MachineId) == machineId {
				userPermitted = true
				break
			}
		}
		if !userPermitted {
			beego.Error("User has no permission to activate the machine")
			this.CustomAbort(401, "Not authorized")
		}
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

// @Title Delete Activation
// @Description Delete an activation
// @Param	aid		path 	int	true		"Activation ID"
// @Success 200 string ok
// @Failure	403	Failed to delete activation
// @Failure 401 Not authorized
// @router /:aid [delete]
func (this *ActivationsController) Delete() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	aid, err := this.GetInt64(":aid")
	if err != nil {
		beego.Error("Failed to get :aid from the request:", err)
		this.CustomAbort(403, "Failed to delete activation")
	}

	err = models.DeleteActivation(aid)
	if err != nil {
		beego.Error("Failed to delete activation:", err)
		this.CustomAbort(403, "Failed to delete activation")
	}

	this.Data["json"] = "ok"
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
