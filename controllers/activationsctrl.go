package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
	"github.com/kr15h/fabsmith/models/purchases"
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
	var activations *[]purchases.Activation
	activations, err = purchases.GetActivations(
		startTime, endTime, userId, itemsPerPage, page)
	if err != nil {
		beego.Error("Failed to get activations:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	// Get total activation count
	var numActivations int64
	numActivations, err = purchases.GetNumActivations(
		startTime, endTime, userId)
	if err != nil {
		beego.Error("Failed to get number of activations:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	r := purchases.GetActivationsResponse{}
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
	activations, err := purchases.GetActiveActivations()
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

	// Turn on the machine
	machine := models.Machine{}
	machine.Id = machineId
	err = machine.On()
	if err != nil {
		beego.Error("Failed to turn on machine:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Turn on the connected machines
	var connectedMachines *models.ConnectedMachineList
	connectedMachines, err = models.GetConnectedMachines(machineId)
	if err != nil {
		beego.Warning("Failed to get connected machines:", err)
	} else {
		if err = connectedMachines.On(); err != nil {
			beego.Error("Failed to turn on connected machines:", err)
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	// Continue with creating activation
	var activationId int64
	var startTime time.Time = time.Now()
	activationId, err = purchases.CreateActivation(machineId, userId, startTime)
	if err != nil {
		beego.Error("Failed to create activation:", err)
		this.CustomAbort(403, "Failed to create activation")
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

	var machineId int64
	machineId, err = purchases.GetActivationMachineId(aid)
	if err != nil {
		beego.Error("Failed to get machine ID")
		this.CustomAbort(403, "Failed to close activation")
	}

	// Attempt to switch off the machine first. This is a way to detect
	// network errors early as the users won't be able to end their activation
	// unless the error in the network is fixed.
	machine := models.Machine{}
	machine.Id = machineId
	if err = machine.Off(); err != nil {
		beego.Error("Failed to switch off machine")
		if !this.IsAdmin() {
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	// Get connected machines and try to swich them off as well
	var connectedMachines *models.ConnectedMachineList
	connectedMachines, err = models.GetConnectedMachines(machineId)
	if err != nil {
		beego.Warning("Failed to get connected machines:", err)
	} else {
		if err = connectedMachines.Off(); err != nil {
			beego.Error("Failed to switch off connected machines")
			if !this.IsAdmin() {
				this.CustomAbort(500, "Internal Server Error")
			}
		}
	}

	err = purchases.CloseActivation(aid, time.Now())
	if err != nil {
		beego.Error("Failed to close activation")
		this.CustomAbort(403, "Failed to close activation")
	}

	statusResponse := models.StatusResponse{}
	statusResponse.Status = "ok"
	this.Data["json"] = statusResponse
	this.ServeJson()
}

// @Title PostFeedback
// @Description Post feedback
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	400	Client Error
// @Failure	500 Internal Server Error
// @router /:aid/feedback [post]
func (this *ActivationsController) PostFeedback() {
	aid, err := this.GetInt64(":aid")
	if err != nil {
		beego.Error("Failed to get :aid from the request:", err)
		this.CustomAbort(400, "Failed to save activation feedback")
	}

	satisfaction := this.GetString("satisfaction")

	_, err = models.CreateActivationFeedback(aid, satisfaction)
	if err != nil {
		beego.Error("Failed to save activation feedback:", err)
		this.CustomAbort(403, "Failed to save activation feedback")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}
