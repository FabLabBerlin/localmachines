package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
	"github.com/astaxie/beego"
	"io/ioutil"
	"time"
)

type ActivationsController struct {
	Controller
}

// @Title Get All
// @Description Get all activations
// @Param	startDate		query 	string	true	"Period start date"
// @Param	endDate			query 	string	true	"Period end date"
// @Param	search			query	string	false	"Search term"
// @Param	userId			query 	int		true	"User ID"
// @Param	itemsPerPage	query 	int		true	"Items per page or max number of items to return"
// @Param	page			query 	int		true	"Current page to show"
// @Success 200 {object}
// @Failure	400	Bad request
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [get]
func (this *ActivationsController) GetAll() {

	locId, isAdmin := this.GetLocIdAdmin()
	if !isAdmin {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	startDate := this.GetString("startDate")
	if startDate == "" {
		beego.Error("Missing start date")
		this.CustomAbort(400, "Failed to get activations")
	}

	endDate := this.GetString("endDate")
	if endDate == "" {
		beego.Error("Missing end date")
		this.CustomAbort(400, "Failed to get activations")
	}

	search := this.GetString("search")

	interval, err := lib.NewInterval(startDate, endDate)
	if err != nil {
		this.CustomAbort(400, "Cannot parse interval")
	}

	// Get activations
	activations, err := purchases.GetActivations(locId, interval, search)
	if err != nil {
		beego.Error("Failed to get activations:", err)
		this.CustomAbort(500, "Failed to get activations")
	}

	this.Data["json"] = activations
	this.ServeJSON()
}

// @Title Create
// @Description Create activation manually
// @Param	location_id	query	int	true	"Location ID"
// @Success 200 {object}
// @Failure 400 Bad Request
// @Failure 401 Not authorized
// @Failure 500 Internal Server Error
// @router / [post]
func (this *ActivationsController) Post() {
	locId, isAdmin := this.GetLocIdAdmin()
	if !isAdmin {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	uid, err := this.GetInt64("user")
	beego.Info("user=", uid)
	if err != nil {
		uid, err = this.GetSessionUserId()
		if err != nil {
			beego.Info("Not logged in:", err)
			this.CustomAbort(401, "Not logged in")
		}
	}
	beego.Info("uid=", uid)

	var inv *invoices.Invoice
	if invId, err := this.GetInt64("invoice"); err == nil {
		if inv, err = invoices.Get(invId); err != nil {
			beego.Info("get invoice:", err)
			this.CustomAbort(400, "Get invoice")
		}
		if inv.Status != "draft" {
			beego.Error("invoice status is", inv.Status)
			this.CustomAbort(400, "Invoice status must be draft")
		}
	} else {
		if inv, err = invoices.GetDraft(locId, uid, time.Now()); err != nil {
			beego.Info("get invoice draft:", err)
			this.CustomAbort(400, "Get invoice draft")
		}
	}

	p := &purchases.Purchase{
		LocationId: locId,
		TimeStart:  time.Now(),
		TimeEnd:    time.Now().Add(time.Hour),
		InvoiceId:  inv.Id,
		UserId:     uid,
		Type:       purchases.TYPE_ACTIVATION,
	}

	if err := purchases.Create(p); err != nil {
		beego.Error("Create activation:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = p
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
	id, err := this.GetInt64(":aid")
	if err != nil {
		this.CustomAbort(400, "Bad request")
	}

	a, err := purchases.GetActivation(id)
	if err != nil {
		beego.Error("get activation:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdminAt(a.Purchase.LocationId) {
		beego.Error("Unauthorized attempt to get activation")
		this.CustomAbort(401, "Unauthorized")
	}

	this.Data["json"] = a
	this.ServeJSON()
}

// @Title Put
// @Description Update activation
// @Success 201 {object}
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:rid [put]
func (this *ActivationsController) Put() {
	locId, isAdmin := this.GetLocIdAdmin()
	if !isAdmin {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	activation := &purchases.Activation{}

	buf, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		beego.Error("Failed to read all:", err)
		this.CustomAbort(400, "Failed to read all")
	}
	defer this.Ctx.Request.Body.Close()

	data := bytes.NewBuffer(buf)

	dec := json.NewDecoder(data)
	if err := dec.Decode(&activation.Purchase); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(400, "Failed to update Activation")
	}

	inv, err := invoices.GetDraft(locId, activation.Purchase.UserId, activation.Purchase.TimeStart)
	if err != nil {
		beego.Error("Get draft:", err)
		this.Abort("500")
	}

	if inv.Status != "draft" {
		beego.Error("cannot edit because invoice in status", inv.Status)
		this.Abort("500")
	}

	activation.Purchase.InvoiceId = inv.Id

	if activation.Purchase.Type == purchases.TYPE_ACTIVATION {
		m, err := machine.Get(activation.Purchase.MachineId)
		if err != nil {
			beego.Error("Failed to get machine:", err)
			this.CustomAbort(500, "Internal Server Error")
		}
		if !this.IsAdminAt(m.LocationId) {
			beego.Error("Unauthorized attempt to update activation")
			this.CustomAbort(401, "Unauthorized")
		}

		if err := activation.Update(); err != nil {
			beego.Error("Failed to update activation:", err)
			this.CustomAbort(500, "Failed to update Activation")
		}
	} else {
		beego.Error("Only activations updates here")
		this.CustomAbort(500, "Only activations updates here")
	}

	this.Data["json"] = activation
	this.ServeJSON()
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

// @Title Start
// @Description Start new activation
// @Param	mid		query 	int	true		"Machine ID"
// @Success 201 {object} models.ActivationCreateResponse
// @Failure	403	Failed to start activation
// @Failure 401 Not authorized
// @router /start [post]
func (this *ActivationsController) Start() {
	locId, isStaff := this.GetLocIdMember()
	if !isStaff {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	if loc, err := locations.Get(locId); err == nil {
		if loc.LocalIp != "" {
			if this.ClientIp() != loc.LocalIp {
				beego.Error("remote user detected, no activation allowed")
				this.CustomAbort(403, "No remote activation")
			}
		}
	} else {
		beego.Error("Failed to get location:", err)
		this.CustomAbort(500, "Failed to create activation")
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
		userPermissions, err := user_permissions.Get(userId)
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

	machine, err := machine.Get(machineId)
	if err != nil {
		beego.Error("Unable to get machine:", err)
		this.CustomAbort(500, "Unable to get machine")
	}

	if err = machine.On(userId); err != nil {
		beego.Error("Failed to turn on machine:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.ServeJSON()
}

// @Title Close
// @Description Close running activation
// @Param	aid		path 	int	true		"Activation ID"
// @Success 200 {object} models.StatusResponse
// @Failure	403	Failed to close activation
// @Failure 401 Not authorized
// @router /:aid/close [post]
func (this *ActivationsController) Close() {
	aid, err := this.GetInt64(":aid")
	if err != nil {
		beego.Error("Failed to get :aid")
		this.CustomAbort(403, "Failed to close activation")
	}

	a, err := purchases.GetActivation(aid)
	if err != nil {
		beego.Error("Unable to get activation:", err)
		this.CustomAbort(500, "Unable to get activation")
	}

	userId, err := this.GetSessionUserId()
	if err != nil {
		beego.Error("Failed to get session user ID:", err)
		this.CustomAbort(403, "Failed to create activation")
	}

	// Attempt to switch off the machine first. This is a way to detect
	// network errors early as the users won't be able to end their activation
	// unless the error in the network is fixed.
	machine, err := machine.Get(a.Purchase.MachineId)
	if err != nil {
		beego.Error("Unable to get machine:", err)
		this.CustomAbort(500, "Unable to get machine")
	}
	if err = machine.Off(userId); err != nil {
		beego.Error("Failed to switch off machine")
		if !this.IsAdminAt(machine.LocationId) {
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	this.ServeJSON()
}
