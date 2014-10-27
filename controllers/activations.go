package controllers

import (
	"database/sql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
	"time"
)

// Handles API /activations requests
type ActivationsController struct {
	Controller
}

// JSON data that is returned via JSON serve
type PublicActivation struct {
	Id        int
	MachineId int
	UserId    int
}

// Creates an activation of a machine for a user. If user not set, session
// user ID is used instead.
func (this *ActivationsController) CreateActivation() {
	reqMachineId, err := this.GetInt(REQUEST_FIELD_NAME_MACHINE_ID)
	if err != nil {
		this.serveErrorResponse("Missing machine ID")
	}
	reqUserId, err := this.GetInt(REQUEST_FIELD_NAME_USER_ID)
	var userId int
	if err != nil {
		userId, err = this.getSessionUserId()
		if err != nil {
			this.serveErrorResponse("Session error occured")
		}
	} else {
		if !this.isAdmin() && !this.isStaff() {
			this.serveErrorResponse("You do not have permissions to create activations for other users")
		}
		userId = int(reqUserId)
	}
	var id int
	id, err = this.createActivation(userId, int(reqMachineId))
	if err != nil {
		this.serveErrorResponse("Could not create activation")
	} else {
		this.serveCreatedResponse(id)
	}
}

// Gets current activations for a user. If user ID not set, use current session
// user ID
func (this *ActivationsController) GetActivations() {
	reqUserId, err := this.GetInt(REQUEST_FIELD_NAME_USER_ID)
	var userId int
	if err != nil {
		// User ID not set, attempt to use session user ID
		userId, err = this.getSessionUserId()
		if err != nil {
			// Can't get session user id, exit
			this.serveErrorResponse("Could not get activations")
		}
	} else {
		// Use request user ID
		userId = int(reqUserId)
	}
	var rawActivations *[]models.Activation
	rawActivations, err = this.getActivationsForUserId(userId)
	if err != nil {
		this.serveErrorResponse("Could not get activations")
	}
	if len(*rawActivations) <= 0 {
		this.serveErrorResponse("No activations found")
	}
	// Now we need to interpret raw activation models for public output
	pubActivations := this.getPublicActivations(rawActivations)
	this.Data["json"] = pubActivations
	this.ServeJson()
}

// Closes active activation. Nothing is being deleted - the activation end time
// is being registered and total time calculated
func (this *ActivationsController) CloseActivation() {
	reqActivationId, err := this.GetInt(REQUEST_FIELD_NAME_ACTIVATION_ID)
	if err != nil {
		this.serveErrorResponse("Missing activation ID")
	}
	err = this.finalizeActivation(int(reqActivationId))
	if err != nil {
		this.serveErrorResponse("Could not close activation")
	}
	this.serveOkResponse()
}

func (this *ActivationsController) getActivationsForUserId(userId int) (*[]models.Activation, error) {
	activations := new([]models.Activation)
	o := orm.NewOrm()
	beego.Trace("Attempt to get activations for user ID", userId)
	num, err := o.Raw("SELECT * FROM activation WHERE user_id = ? AND active = 1",
		userId).QueryRows(activations)
	if err != nil {
		beego.Error("Could not get activations")
		return nil, err
	}
	beego.Trace("Got", num, "activations")
	return activations, nil
}

// Interprets *[]models.Activation to *[]PublicActivation and returns same size array
func (this *ActivationsController) getPublicActivations(activations *[]models.Activation) *[]PublicActivation {
	pubActivations := []PublicActivation{}
	for i := range *activations {
		item := PublicActivation{(*activations)[i].Id,
			(*activations)[i].MachineId,
			(*activations)[i].UserId}
		pubActivations = append(pubActivations, item)
	}
	return &pubActivations
}

// Create new activation with user ID and machine ID.
// Return created activation ID.
func (this *ActivationsController) createActivation(userId int, machineId int) (int, error) {
	// Beego model time stuff is bad, use workaround that works
	beego.Trace("Attempt to create activation")
	o := orm.NewOrm()
	res, err := o.Raw("INSERT INTO activation VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		nil, nil, userId, machineId, true,
		time.Now().Format("2006-01-02 15:04:05"),
		nil, 0, 0, 0, 0, 0, "", false, false).Exec()
	if err != nil {
		beego.Error("Failed to insert activation in to DB:", err)
		return 0, err
	}
	var activationId int64
	activationId, err = res.LastInsertId()
	if err != nil {
		beego.Critical("Last inserted ID not set, something terribly wrong")
		panic("Failing")
	}
	beego.Trace("Activation with ID", activationId, "created")
	err = this.setMachineUnavailable(machineId)
	return int(activationId), err
}

// Finalize activation, save end time
func (this *ActivationsController) finalizeActivation(id int) error {
	activationModel := new(models.Activation)
	activationModel.Id = id
	// Get activation start time
	beego.Trace("Attempt to get activation with ID", id)
	var tempModel struct {
		TimeStart string
		MachineId int
	}
	o := orm.NewOrm()
	err := o.Raw("SELECT time_start, machine_id FROM activation WHERE active = true AND id = ?",
		id).QueryRow(&tempModel)
	if err != nil {
		beego.Error("Could not get activation:", err)
		return err
	}
	beego.Trace("Success")
	// Convert start time into Unix timestamp, workaround here as Beego models
	// do not handle time properly
	const timeForm = "2006-01-02 15:04:05"
	timeStart, _ := time.ParseInLocation(timeForm, tempModel.TimeStart, time.Now().Location())
	timeNow := time.Now() // time.Now().Format("2006-01-02 15:04:05")
	totalDuration := timeNow.Sub(timeStart)
	// Set model vars
	activationModel.TimeEnd = timeNow
	activationModel.TimeTotal = int(totalDuration.Seconds())
	activationModel.Active = false
	// Update database
	beego.Trace("Attempt to update activation")
	var res sql.Result
	res, err = o.Raw("UPDATE activation SET active=false, time_end=?, time_total=? WHERE id=?",
		timeNow.Format("2006-01-02 15:04:05"),
		totalDuration.Seconds(), id).Exec()
	if err != nil {
		beego.Error("Failed to update activation:", err)
		if beego.AppConfig.String("runmode") == "dev" {
			panic("Failed to update activation")
		} else {
			return err
		}
	}
	var rowsAffected int64
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		beego.Critical("Could not get num affected rows")
		if beego.AppConfig.String("runmode") == "dev" {
			panic("Could not get num affected rows")
		} else {
			return err
		}
	}
	beego.Trace("Success, rows affected: ", rowsAffected)
	err = this.setMachineAvailable(tempModel.MachineId)
	return err
}

func (this *ActivationsController) setMachineUnavailable(id int) error {
	machineModel := new(models.Machine)
	machineModel.Id = id
	machineModel.Available = false
	o := orm.NewOrm()
	beego.Trace("Attempt to set machine with ID", id, "unavailable")
	num, err := o.Update(machineModel, "Available")
	if err != nil {
		beego.Error("Failed:", err)
		return err
	}
	beego.Trace("Success, rows affected: ", num)
	return nil
}

func (this *ActivationsController) setMachineAvailable(id int) error {
	machineModel := new(models.Machine)
	machineModel.Id = id
	machineModel.Available = true
	o := orm.NewOrm()
	beego.Trace("Attempt to set machine with ID", id, "available")
	num, err := o.Update(machineModel, "Available")
	if err != nil {
		beego.Error("Failed:", err)
		return err
	}
	beego.Trace("Success, rows affected: ", num)
	return nil
}
