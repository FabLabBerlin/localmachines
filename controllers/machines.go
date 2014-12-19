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
	var response struct{ Machines []PublicMachine }
	reqUserId, err := this.GetInt(REQUEST_FIELD_NAME_USER_ID)
	var userId int
	if err != nil {
		beego.Info("No user ID set, attempt to get session user ID")
		userId, err = this.getSessionUserId()
		if err != nil {
			if beego.AppConfig.String("runmode") == "dev" {
				panic("Could not get session user ID")
			}
			this.serveErrorResponse("There was an error")
		}
	} else {
		if !this.isAdmin() && !this.isStaff() {
			this.serveErrorResponse("You don't have permissions to list other user's machines")
		} else {
			userId = int(reqUserId)
		}
	}
	machines, err := this.getUserMachines(userId)
	if err != nil {
		// TODO: serve empty array if no machines found, error on real error
		this.serveErrorResponse("No machines available")
	}
	response.Machines = machines
	this.Data["json"] = &response
	this.ServeJson()
}

func (this *MachinesController) getUserMachines(userId int) ([]PublicMachine, error) {
	beego.Trace("Attempt to get machines for user ID:", userId)
	machines := []models.Machine{}
	o := orm.NewOrm()
	num, err := o.Raw("SELECT * FROM machine INNER JOIN permission ON machine.id = permission.machine_id WHERE permission.user_id = ?",
		userId).QueryRows(&machines)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	beego.Trace("Got ", num, "machines")

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
				if activation.UserId == userId {
					status = MACHINE_STATUS_USED
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
