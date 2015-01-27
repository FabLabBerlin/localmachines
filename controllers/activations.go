package controllers

import (
	"database/sql"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
	"github.com/kr15h/fabsmith/plugins/hexaswitch"
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

	// Get machine ID to be activated from the request variables
	reqMachineId, err := this.GetInt(REQUEST_FIELD_NAME_MACHINE_ID)
	if err != nil {
		this.serveErrorResponse("Missing machine ID")
	}

	// Get the user ID for the activation
	reqUserId, err := this.GetInt(REQUEST_FIELD_NAME_USER_ID)
	var userId int
	if err != nil {

		// It there is no user ID in the request, get it from session -
		// logged in user ID
		userId, err = this.getSessionUserId()
		if err != nil {

			// It was impossible to get the user ID, fail now
			if beego.AppConfig.String("runmode") == "dev" {
				panic("Could not get session user ID")
			}
			this.serveErrorResponse("Session error occured")
		}
	} else {

		// There is an user ID in the request, check if we are authorized
		// to use it (only if we are admin or staff)
		if !this.isAdmin() && !this.isStaff() {

			// We don't have admin or staff privileges, so fail now
			this.serveErrorResponse("You do not have permissions to create activations for other users")
		}
		userId = int(reqUserId)
	}

	// Create activation and get it's ID
	var activationId int
	activationId, err = this.createActivation(userId, int(reqMachineId))
	if err != nil {

		// Could not create valid activation, fail now
		this.serveErrorResponse("Could not create activation")
	} else {

		// Activation successfully created, serve "Activation Created" response
		this.serveCreatedResponse(activationId)
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
			if beego.AppConfig.String("runmode") == "dev" {
				panic("Could not get session user ID")
			}
			this.serveErrorResponse("Could not get activations")
		}
	} else {
		// Use request user ID
		userId = int(reqUserId)
	}
	var rawActivations []models.Activation
	rawActivations, err = this.getActivationsForUserId(userId)
	if err != nil {
		if beego.AppConfig.String("runmode") == "dev" {
			panic("Could not get activations")
		}
		this.serveErrorResponse("Could not get activations")
	}
	if len(rawActivations) <= 0 {
		this.serveErrorResponse("No activations found")
	}
	// Now we need to interpret raw activation models for public output
	pubActivations := this.getPublicActivations(&rawActivations)
	this.Data["json"] = &pubActivations
	this.ServeJson()
}

// Closes active activation. Nothing is being deleted - the activation end time
// is being registered and total time calculated
func (this *ActivationsController) CloseActivation() {
	reqActivationId, err := this.GetInt(REQUEST_FIELD_NAME_ACTIVATION_ID)
	if err != nil {
		this.serveErrorResponse("Missing activation ID")
	} else {
		err = this.finalizeActivation(int(reqActivationId))
		if err != nil {

			// Could not finalize activation, serve error
			this.serveErrorResponse("Could not close activation")
		}
		this.serveOkResponse()
	}
}

// Returns activations for user ID specified
func (this *ActivationsController) getActivationsForUserId(userId int) ([]models.Activation, error) {
	activations := []models.Activation{}
	o := orm.NewOrm()
	beego.Trace("Attempt to get activations for user ID", userId)
	num, err := o.Raw("SELECT * FROM activation WHERE user_id = ? AND active = 1",
		userId).QueryRows(&activations)
	if err != nil {
		beego.Error("Could not get activations")
		return nil, err
	}
	beego.Trace("Got", num, "activations")
	return activations, nil
}

// Interprets *[]models.Activation to *[]PublicActivation and returns same size array
func (this *ActivationsController) getPublicActivations(activations *[]models.Activation) []PublicActivation {
	pubActivations := []PublicActivation{}
	for i := range *activations {
		item := PublicActivation{(*activations)[i].Id,
			(*activations)[i].MachineId,
			(*activations)[i].UserId}
		pubActivations = append(pubActivations, item)
	}
	return pubActivations
}

// Create new activation with user ID and machine ID.
// Return created activation ID.
func (this *ActivationsController) createActivation(userId int, machineId int) (int, error) {

	// Check for duplicate activations
	o := orm.NewOrm()
	var dupActivations []models.Activation
	beego.Info("Checking for duplicate activations")
	num, err := o.Raw("SELECT id FROM activation WHERE machine_id = ? AND user_id = ? AND active = 1",
		machineId, userId).QueryRows(&dupActivations)
	if err != nil {
		beego.Error(err)
		return 0, err
	}

	// Got duplicates, that's bad
	if num > 0 {
		beego.Error("Duplicate activations found:", num)
		return 0, errors.New("Duplicate activations found")
	}

	// No duplicates - good
	beego.Info("No duplicate activations found")

	// Beego model time stuff is bad, here we use workaround that works.
	// TODO: explore the beego ORM time management,
	// try to fix or use as it should be used.

	beego.Trace("Attempt to create activation")
	var res sql.Result
	res, err = o.Raw("INSERT INTO activation VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		nil, nil, userId, machineId, true,
		time.Now().Format("2006-01-02 15:04:05"),
		nil, 0, 0, 0, 0, 0, "", false, false).Exec()
	if err != nil {
		// Failed to create database record, fail now
		beego.Error("Failed to insert activation in to DB:", err)
		return 0, err
	}

	// Get the ID of the record we just inserted
	var activationId int64
	activationId, err = res.LastInsertId()
	if err != nil {
		beego.Error("Failed to get insterted activation ID")
		return int(0), err
	}
	beego.Trace("Activation with ID", activationId, "created")

	// Set the machine or item unavailable as we just activated it
	err = this.setMachineUnavailable(machineId)
	if err != nil {

		beego.Error("Failed to set machine unavailable, machine ID", machineId)

		_ = this.setMachineAvailable(machineId)

		// Clear the activation initiated earlier
		beego.Info("Attempt to clear activation", activationId)
		err = this.clearActivation(activationId)
		if err != nil {
			// This should never happen
			beego.Error("Failed to clear activation", activationId)
			return 0, errors.New("Failed to clear activation")
		}

		return 0, errors.New("Failed to set machine unavailable")
	}

	// Check if there is mapping between switch and machine
	_, err = hexaswitch.GetSwitchIp(machineId)
	if err != nil {
		beego.Warning("Machine / switch mapping does not exist")

		_ = this.setMachineAvailable(machineId)

		// Clear the activation initiated earlier
		beego.Info("Attempt to clear activation", activationId)
		err = this.clearActivation(activationId)
		if err != nil {
			beego.Error("Failed to clear activation", activationId)
			return 0, errors.New("Failed to clear activation")
		}

		// Nothing bad with not having a mapping
		// maybe it's a screwdriver - so we keep the activation
		return int(activationId), nil
	}

	// TODO: use switch IP for the SwitchOn method

	// Try to turn on the switch as we have the mapping and all
	hexaswitch.Install()                 // TODO: remove this from here in an elegant way
	err = hexaswitch.SwitchOn(machineId) // This will take some time
	if err != nil {

		tempErr := err
		beego.Error("Failed to turn on the switch", tempErr)

		_ = this.setMachineAvailable(machineId)

		beego.Info("Attempt to remove activation with ID", activationId)
		err = this.clearActivation(activationId)
		if err != nil {
			beego.Error("Failed to remove activation with ID", activationId)
		}
		beego.Info("Successfuly removed activation with ID", activationId)

		return 0, tempErr
	}

	return int(activationId), err
}

// Removes activation in case something else related to that is wrong
func (this *ActivationsController) clearActivation(activationId int64) error {

	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM activation WHERE id=?", activationId).Exec()
	if err != nil {
		return err
	}

	return nil
}

// Finalize activation, save activation end time
func (this *ActivationsController) finalizeActivation(id int) error {

	// Here we can use beego ORM models in a safe way
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

		// Failed to get activation from the database
		beego.Error("Could not get activation:", err)
		return err
	}
	beego.Trace("Success")

	// Launch go routine to switch off the machine after some time
	var deactivateTimeout int64
	deactivateTimeout, err = beego.AppConfig.Int64("deactivatetimeout")
	if err != nil {
		beego.Error("Failed to get deactivate timeout from config:", err)
		deactivateTimeout = 30
	}
	go deactivateMachineAfterTimeout(tempModel.MachineId, deactivateTimeout)

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

		// Failed to update database
		beego.Error("Failed to update activation:", err)
		return err
	}
	var rowsAffected int64
	rowsAffected, err = res.RowsAffected()
	if err != nil {

		// Failed to get affected rows from the database
		beego.Critical("Could not get num affected rows")
		return err
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

// Deativates a machine after timeout if no activation with the machine ID
// has been made.
func deactivateMachineAfterTimeout(machineId int, timeoutSeconds int64) {

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
	hexaswitch.Install()                  // TODO: remove this from here in an elegant way
	err = hexaswitch.SwitchOff(machineId) // This will take some time
	if err != nil {
		beego.Error("Failed to turn the switch off", err)
	}
}
