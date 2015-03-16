package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
	"time"
)

type MachinesController struct {
	Controller
}

type PublicMachine struct {
	Id                       int
	Name                     string
	Description              string
	Price                    float32
	PriceUnit                string
	PriceCurrency            string
	Status                   string
	OccupiedByUserId         int
	OccupiedByUserFullName   string
	ActivationId             int
	ActivationSecondsElapsed int64
	UnavailableMessage       string
}

// Status constants
const MACHINE_STATUS_AVAILABLE = "free"
const MACHINE_STATUS_OCCUPIED = "occupied"
const MACHINE_STATUS_USED = "used"
const MACHINE_STATUS_UNAVAILABLE = "unavailable"

// Output JSON list with available machines for a user
func (this *MachinesController) GetMachines() {

	// Prepare output array for the machines
	var response struct{ Machines []PublicMachine }

	// Get request user ID
	var reqUserId int = 0
	var err error
	reqUserId, err = this.GetInt(REQUEST_FIELD_NAME_USER_ID)
	if err != nil {
		beego.Error("Could not get request user ID")
	}

	// Get session user ID
	var sessionUserId int = 0
	sessionUserId, err = this.getSessionUserId()
	if err != nil {
		beego.Error("Could not get session user ID")
		this.serveErrorResponse("Could not get session user ID")
	}

	// If not Admin or Staff, cant get machines of another user
	// and that means if sessionUserId does not match with reqUserId
	// we throw an error
	var machinesUserId int
	if sessionUserId != reqUserId && !this.isAdmin() && !this.isStaff() {
		
		// Just return machines of the current user
		machinesUserId = sessionUserId
	} else {

		// Request made by admin or staff, 
		// allow to get machines of another user
		machinesUserId = reqUserId
	}

	// Get array of machines for current user
	machines, err := this.getUserMachines(machinesUserId)

	if err != nil {
		beego.Error(fmt.Sprintf("Error getting machines for user ID: %v", machinesUserId))

		// Create empty machine array for the response and serve
		emptyArray := []PublicMachine{}
		response.Machines = emptyArray
		this.Data["json"] = &response
		this.ServeJson()
	} else {

		// else respond with array full of machines
		response.Machines = machines
		this.Data["json"] = &response
		this.ServeJson()
	}

}

func (this *MachinesController) getUserMachines(userId int) ([]PublicMachine, error) {

	// Prepare empty array for machines stored in the database
	machines := []models.Machine{}

	// Check if we are admin or staff
	if this.isAdmin() || this.isStaff() {

		// If user is admin or staff, get all machines
		beego.Trace("Attemt to get all machines")

		// Attempt to get all machines
		o := orm.NewOrm()
		num, err := o.Raw("SELECT * FROM machine").QueryRows(&machines)

		if err != nil {
			beego.Error(err)
			return nil, err
		}

		beego.Trace("Got", num, "machines")

	} else {

		// If user is just a member or nothing, get user machines
		beego.Trace("Attempt to get machines for user ID", userId)

		// Get machines for user ID
		o := orm.NewOrm()
		num, err := o.Raw("SELECT * FROM machine INNER JOIN permission ON machine.id = permission.machine_id WHERE permission.user_id = ?",
			userId).QueryRows(&machines)

		if err != nil {
			beego.Error(err)
			return nil, err
		}

		beego.Trace("Got", num, "machines")

	}

	// Interpret machines as PublicMachine
	var pubMachines = []PublicMachine{}
	for i := range machines {

		var price float32 = 0.0
		var priceUnit string = "unit"

		if machines[i].CalcByEnergy {
			price = machines[i].CostsPerKwh
			priceUnit = "KWh"
		} else if machines[i].CalcByTime {
			price = machines[i].CostsPerMin
			priceUnit = "minute"
		}

		var status string = MACHINE_STATUS_AVAILABLE
		var occupiedByUserId int
		var occupiedByUsesFullName string

		var activationId int = 0
		var activationSecondsElapsed int64 = 0
		var unavailableMessage string = ""

		if !machines[i].Available {

			status = MACHINE_STATUS_OCCUPIED

			// Machine is not available, check if there is an activation with it
			activation, err := this.getActivation(machines[i].Id)
			if err != nil {
				// TODO: output unavail message
				// TODO: status = "unavailable"
				//occupiedBy = "Unavailable"
				status = MACHINE_STATUS_UNAVAILABLE
				//return nil, err
				occupiedByUserId = 0
				unavailableMessage = machines[i].UnavailMsg
			} else {
				occupiedByUserId = activation.UserId

				// Get user full name right away
				var userModel models.User
				userModel.Id = activation.UserId
				o := orm.NewOrm()
				err = o.Read(&userModel)
				if err != nil {
					beego.Error("Could not get user responsible for the activation:", err)
					return []PublicMachine{}, err
				} else {
					occupiedByUsesFullName = fmt.Sprintf("%s %s", userModel.FirstName, userModel.LastName)
				}

				// Change status to "USED" if user id is the same as logged user ID
				if (activation.UserId == userId) || (this.isAdmin() || this.isStaff()) {

					if activation.UserId == userId {
						status = MACHINE_STATUS_USED
					}
					activationId = activation.Id

					// We need to get raw time as beego does something strange with it
					// TODO: has to be fixed
					beego.Trace("Attempt to get start time as string for activation ID ", activation.Id)
					var tempModel struct {
						TimeStart string
					}
					err = o.Raw("SELECT time_start FROM activation WHERE id = ?",
						activation.Id).QueryRow(&tempModel)
					if err != nil {
						beego.Error("Could not get activation:", err)
					}
					beego.Trace("Successfuly got activation start time")

					// Pass that string to our output machine array
					const timeForm = "2006-01-02 15:04:05"
					timeStart, _ := time.ParseInLocation(timeForm, tempModel.TimeStart, time.Now().Location())
					//activationStartTime = timeStart.Format(timeForm)

					// Set current time
					//currentTime = time.Now().Format(timeForm)

					activationSecondsElapsed = int64(time.Now().Sub(timeStart).Seconds())
				}
			}
		}

		// Fill public machine struct for output
		machine := PublicMachine{
			Id:                       machines[i].Id,
			Name:                     machines[i].Name,
			Description:              machines[i].Description,
			Price:                    price,
			PriceUnit:                priceUnit,
			PriceCurrency:            "€", // TODO: add price currency table
			Status:                   status,
			OccupiedByUserId:         occupiedByUserId,
			OccupiedByUserFullName:   occupiedByUsesFullName,
			ActivationId:             activationId,
			ActivationSecondsElapsed: activationSecondsElapsed,
			UnavailableMessage:       unavailableMessage}
		// Append to array
		pubMachines = append(pubMachines, machine)
	} // for
	return pubMachines, nil
}

func (this *MachinesController) getActivation(machineId int) (models.Activation, error) {
	o := orm.NewOrm()
	activationModel := models.Activation{}
	beego.Trace("Attempt to get activation for machine ID", machineId)

	err := o.Raw("SELECT * FROM activation WHERE machine_id = ? AND active = 1",
		machineId).QueryRow(&activationModel)
	if err != nil {
		beego.Error(err)
		return models.Activation{}, err
	}

	beego.Trace("Success")
	return activationModel, nil
}
