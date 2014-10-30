package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

type MachinesController struct {
	Controller
}

type PublicMachine struct {
	Id                     int
	Name                   string
	Description            string
	Price                  float32
	PriceUnit              string
	PriceCurrency          string
	Status                 string
	OccupiedByUserId       int
	OccupiedByUserFullName string
}

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
		var status string = "free"
		var occupiedByUserId int
		var occupiedByUsesFullName string
		if !machines[i].Available {
			status = "occupied"
			// Machine is not available, check if there is an activation with it
			activation, err := this.getActivation(machines[i].Id)
			if err != nil {
				// TODO: output unavail message
				// TODO: status = "unavailable"
				//occupiedBy = "Unavailable"
				status = "unavailable"
				//return nil, err
				occupiedByUserId = 0
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
			}
		}
		// Fill public machine struct for output
		machine := PublicMachine{
			Id:                     machines[i].Id,
			Name:                   machines[i].Name,
			Description:            machines[i].Description,
			Price:                  price,
			PriceUnit:              priceUnit,
			PriceCurrency:          "€", // TODO: add price currency table
			Status:                 status,
			OccupiedByUserId:       occupiedByUserId,
			OccupiedByUserFullName: occupiedByUsesFullName}
		// Append to array
		pubMachines = append(pubMachines, machine)
	}
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
