package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

// Handles API /activations requests
type ActivationsController struct {
	Controller
}

type PublicActivation struct {
	Id        int
	MachineId int
	UserId    int
}

// Creates an activation for a user
func (this *ActivationsController) CreateActivation() {
	response := struct{ Status string }{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}

// Gets current activations for a user
func (this *ActivationsController) GetActivations() {
	// Check if request has user_id variable
	reqUserId, hasUserId := this.requestHasUserId()
	var rawActivations []models.Activation
	var err error
	if hasUserId {
		// Request has user id, use that to get activations
		beego.Info("Request HAS user_id")
		rawActivations, err = this.getActivationsForUserId(reqUserId)
		if err != nil {
			beego.Error(err)
			this.serveStatusResponse("error", "Could not get activations")
		}
	} else {
		// Request does not have user id, attemt to use user id from session
		beego.Info("Request does NOT have user_id")
		userId := this.GetSession("user_id")
		if userId == nil {
			// Could not find session user id, that means we're out of luck and
			// something has gone terribly wrong
			beego.Error("Could not find any user id")
			this.serveStatusResponse("error", "No usable user ID found")
		}
		// Ok, we have session user ID, use it to get activations
		rawActivations, err = this.getActivationsForUserId(userId.(int))
		if err != nil {
			beego.Error(err)
			this.serveStatusResponse("error", "Could not get activations")
		}
	}
	// Check how many activations
	if len(rawActivations) <= 0 {
		beego.Error("There are no activations")
		this.serveStatusResponse("error", "No activations found")
	}
	// Now we need to interpret them for public output
	pubActivations := this.getPublicActivations(rawActivations)
	this.Data["json"] = &pubActivations
	this.ServeJson()
}

// Closes active activation. Nothing is being deleted - the activation end time
// is being registered and total time calculated
func (this *ActivationsController) UpdateActivation() {
	response := struct{ Status string }{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}

func (this *ActivationsController) getActivationsForUserId(userId int) ([]models.Activation, error) {
	var activations []models.Activation
	o := orm.NewOrm()
	beego.Info("Attempt to get activations for user id", userId)
	num, err := o.Raw("SELECT * FROM activation WHERE user_id = ? AND active = 1",
		userId).QueryRows(&activations)
	if err != nil {
		beego.Error("Could not get activations")
		return nil, err
	}
	beego.Info("Got", num, "Activations")
	return activations, nil
}

// Interprets models.Activation to PublicActivation and returns same size array
func (this *ActivationsController) getPublicActivations(activations []models.Activation) []PublicActivation {
	var pubActivations []PublicActivation
	for i := range activations {
		act := PublicActivation{activations[i].Id,
			activations[i].MachineId,
			activations[i].UserId}
		pubActivations = append(pubActivations, act)
	}
	return pubActivations
}
