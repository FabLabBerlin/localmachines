package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
	"time"
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
	// Get request machine ID
	reqHasMachineId, reqMachineId := this.requestHasMachineId()
	if !reqHasMachineId {
		// If request machine ID not set
		// Serve error message
		this.serveStatusResponse("error", "Missing machine ID")
	}
	// Check if we have user ID passed with the request
	reqHasUserId, reqUserId := this.requestHasUserId()
	if reqHasUserId {
		// If we have request user ID
		// Check if session user is authorized to use it
		if this.isAdmin() || this.isStaff() {
			// If session user is admin or staff
			beego.Info("User is admin or staff")
			// Create activation with request user ID and request's machine ID
			err := this.createActivation(reqUserId, reqMachineId)
			if err != nil {
				this.serveStatusResponse("error", "Could not activate machine")
			}
			// Serve created response with activation id
			this.serveSuccessResponse()
		} else {
			// If session user IS NOT admin or staff
			beego.Info("User not admin or staff")
			// Serve error message
			this.serveStatusResponse("error", "Not authorized")
		}
	} else {
		// If we DO NOT have request user ID
		// Use session user ID and machine ID to create activation
		sessUserId, err := this.getSessionUserId()
		if err != nil {
			// If there IS NO session user ID
			// Serve error message
			this.serveStatusResponse("error", "Could not get session user ID")
		}
		// Attempt to create activation
		err = this.createActivation(sessUserId, reqMachineId)
		if err == nil {
			// If successful
			// Serve success message
			this.serveSuccessResponse()
		} else {
			// If could not create activation
			// Serve error message
			this.serveStatusResponse("error", "Could not create activation")
		}
	}
}

// Gets current activations for a user
func (this *ActivationsController) GetActivations() {
	// Check if request has user_id variable
	hasUserId, reqUserId := this.requestHasUserId()
	var rawActivations []models.Activation
	var err error
	if hasUserId {
		// Request has user id, use that to get activations
		beego.Info("Request HAS user_id")
		// TODO: check if session user is authorized to get activations for other
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
	// Get activation id
	activationId, err := this.getRequestActivationId()
	// If no activation id serve error
	if err != nil {
		this.serveStatusResponse("error", "No activation ID")
	}
	// Else continue and finalize activation
	err = this.finalizeActivation(activationId)
	if err != nil {
		// If there is an error, serve error message
		this.serveStatusResponse("error", "Could not finalize activation")
	}
	// Else activation finalized serve success json
	this.serveSuccessResponse()
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

// Create new activation with user ID and machine ID
func (this *ActivationsController) createActivation(userId int, machineId int) error {
	o := orm.NewOrm()
	// Beego model time stuff is bad, use workaround that works
	beego.Info("Attempt to create activation")
	_, err := o.Raw("INSERT INTO activation VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		nil, nil, userId, machineId, true,
		time.Now().Format("2006-01-02 15:04:05"),
		nil, 0, 0, 0, 0, 0, "", false, false).Exec()
	if err != nil {
		beego.Error("Failed to create activation for user ID",
			userId, "and machine ID", machineId, ":", err)
		return err
	}
	beego.Info("Created activation")
	return nil
}

// Get activation_id variable from the request
func (this *ActivationsController) getRequestActivationId() (int, error) {
	beego.Info("Attempt to get request activation ID")
	id, err := this.GetInt("activation_id")
	if err != nil {
		beego.Error("Could not get request activation ID:", err)
		return int(0), err
	}
	beego.Info("Got request activation ID:", id)
	return int(id), nil
}

// Finalize activation, save end time
func (this *ActivationsController) finalizeActivation(id int) error {
	o := orm.NewOrm()
	activationModel := new(models.Activation)
	activationModel.Id = id
	// Get activation start time
	beego.Info("Attempt to get activation with ID", id)
	var tempModel struct {
		TimeStart string
		MachineId int
	}
	err := o.Raw("SELECT time_start, machine_id FROM activation WHERE active = true AND id = ?",
		id).QueryRow(&tempModel)
	if err != nil {
		beego.Error("Could not get activation:", err)
		return err
	}
	beego.Info("Success")
	// Convert start time into Unix timestamp
	const timeForm = "2006-01-02 15:04:05"
	timeStart, _ := time.ParseInLocation(timeForm, tempModel.TimeStart, time.Now().Location())
	timeNow := time.Now() // time.Now().Format("2006-01-02 15:04:05")
	totalDuration := timeNow.Sub(timeStart)
	// Set model vars
	activationModel.TimeEnd = timeNow
	activationModel.TimeTotal = int(totalDuration.Seconds())
	activationModel.Active = false
	// Update database
	beego.Info("Attempt to update activation")
	_, err = o.Raw("UPDATE activation SET active=false, time_end=?, time_total=? WHERE id=?",
		timeNow.Format("2006-01-02 15:04:05"),
		totalDuration.Seconds(), id).Exec()
	//num, err := o.Update(activationModel, "TimeEnd", "TimeTotal", "Active")
	if err != nil {
		beego.Error("Failed to update activation", err)
		return err
	}
	beego.Info("Success")
	return nil
}
