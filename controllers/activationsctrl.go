package controllers

import (
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego"
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
// @Success 200 {object}
// @Failure	403	Failed to get activations
// @Failure	401	Not authorized
// @router / [get]
func (this *ActivationsController) GetAll() {

	// Only admin can use this API call
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	// Get variables
	startDate := this.GetString("startDate")
	if startDate == "" {
		beego.Error("Missing start date")
		this.CustomAbort(403, "Failed to get activations")
	}

	endDate := this.GetString("endDate")
	if endDate == "" {
		beego.Error("Missing end date")
		this.CustomAbort(403, "Failed to get activations")
	}

	userId, err := this.GetInt64("userId")
	if err != nil {
		beego.Error("Could not get userId request variable:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	itemsPerPage, err := this.GetInt64("itemsPerPage")
	if err != nil {
		beego.Error("Could not get itemsPerPage request variable:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	page, err := this.GetInt64("page")
	if err != nil {
		beego.Error("Could not get page request variable:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	// Convert / parse string time values as time.Time
	var timeForm = "2006-01-02"

	startTime, err := time.ParseInLocation(
		timeForm, startDate, time.Now().Location())
	if err != nil {
		beego.Error("Failed to parse startDate:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	endTime, err := time.ParseInLocation(
		timeForm, endDate, time.Now().Location())
	if err != nil {
		beego.Error("Failed to parse endDate:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	// Get activations
	activations, err := purchases.GetActivations(
		startTime, endTime, userId, itemsPerPage, page)
	if err != nil {
		beego.Error("Failed to get activations:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	// Get total activation count
	numActivations, err := purchases.GetNumActivations(
		startTime, endTime, userId)
	if err != nil {
		beego.Error("Failed to get number of activations:", err)
		this.CustomAbort(403, "Failed to get activations")
	}

	this.Data["json"] = purchases.GetActivationsResponse{
		NumActivations:  numActivations,
		ActivationsPage: activations,
	}
	this.ServeJSON()
}

// @Title Get
// @Description Get activation by activation ID
// @Param	aid		path 	int	true		"Activation ID"
// @Success 200 {object}
// @Failure	403	Failed to get activation
// @Failure	401	Not authorized
// @router /:aid [get]
func (this *ActivationsController) Get() {

}

// @Title Get Active
// @Description Get all active activations
// @Success 200 {object}
// @Failure	403	Failed to get active activations
// @router /active [get]
func (this *ActivationsController) GetActive() {
	activations, err := purchases.GetActiveActivations()
	if err != nil {
		beego.Error("Failed to get active activations")
		this.CustomAbort(403, "Failed to get active activations")
	}
	this.Data["json"] = activations
	this.ServeJSON()
}

// @Title Create
// @Description Create new activation
// @Param	mid		query 	int	true		"Machine ID"
// @Success 201 {object} models.ActivationCreateResponse
// @Failure	403	Failed to create activation
// @Failure 401 Not authorized
// @router / [post]
func (this *ActivationsController) Create() {
	locId, isStaff := this.GetLocIdMember()
	if !isStaff {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	if xff := this.Ctx.Request.Header.Get("X-Forwarded-For"); xff != "" {
		if loc, err := locations.Get(locId); err == nil {
			if loc.LocalIp != "" {
				if xff != loc.LocalIp {
					beego.Error("remote user detected, no activation allowed")
					this.CustomAbort(403, "No remote activation")
				}
			}
		} else {
			beego.Error("Failed to get location:", err)
			this.CustomAbort(500, "Failed to create activation")
		}
	}

	userId, err := this.GetSessionUserId()
	if err != nil {
		beego.Error("Failed to get session user ID:", err)
		this.CustomAbort(403, "Failed to create activation")
	}

	machineId, err := this.GetInt64("mid")
	if err != nil {
		beego.Error("Failed to get mid:", err)
		this.CustomAbort(403, "Failed to create activation")
	}

	// Admins can activate any machine (except broken ones).
	// Regular users have to refer to their permissions.
	if !isStaff {

		// Check if user has permissions to create activation for the machine.
		userPermissions, err := models.GetUserPermissions(userId)
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

	machine, err := machine.GetMachine(machineId)
	if err != nil {
		beego.Error("Unable to get machine:", err)
		this.CustomAbort(500, "Unable to get machine")
	}

	if err = machine.On(); err != nil {
		beego.Error("Failed to turn on machine:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Continue with creating activation
	var startTime time.Time = time.Now()
	activationId, err := purchases.CreateActivation(machineId, userId, startTime)
	if err != nil {
		beego.Error("Failed to create activation:", err)
		this.CustomAbort(403, "Failed to create activation")
	}

	this.Data["json"] = models.ActivationCreateResponse{
		ActivationId: activationId,
	}
	this.ServeJSON()
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

	machineId, err := purchases.GetActivationMachineId(aid)
	if err != nil {
		beego.Error("Failed to get machine ID")
		this.CustomAbort(403, "Failed to close activation")
	}

	// Attempt to switch off the machine first. This is a way to detect
	// network errors early as the users won't be able to end their activation
	// unless the error in the network is fixed.
	machine, err := machine.GetMachine(machineId)
	if err != nil {
		beego.Error("Unable to get machine:", err)
		this.CustomAbort(500, "Unable to get machine")
	}
	if err = machine.Off(); err != nil {
		beego.Error("Failed to switch off machine")
		if !this.IsAdmin() {
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	err = purchases.CloseActivation(aid, time.Now())
	if err != nil {
		beego.Error("Failed to close activation")
		this.CustomAbort(403, "Failed to close activation")
	}

	this.Data["json"] = models.StatusResponse{
		Status: "ok",
	}
	this.ServeJSON()
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
	this.ServeJSON()
}
