package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

type MachinesController struct {
	Controller
}

type PublicMachine struct {
	Id               int
	Name             string
	Description      string
	Price            float32
	PriceUnit        string
	PriceCurrency    string
	Status           string
	OccupiedByUserId int
}

// Output JSON list with available machines for a user
func (this *MachinesController) GetMachines() {
	var response struct{ Machines []PublicMachine }
	// Cehck if there is user_id in the request variables
	userId, isSet := this.requestHasUserId()
	if isSet {
		// Check if admin
		roles := this.getSessionUserRoles()
		// If user role is admin, staff or request user_id matches session user_id
		if roles.Admin || roles.Staff ||
			userId == this.GetSession("user_id").(int) {
			// Use request user_id
			machines, err := this.getUserMachines(userId)
			if err != nil {
				errResponse := ErrorResponse{"error", "No machines available"}
				this.Data["json"] = &errResponse
				this.ServeJson()
			}
			response.Machines = machines
		} else {
			// User has no permission to get machine info for another user
			errResponse := ErrorResponse{"error", "Not authorized"}
			this.Data["json"] = &errResponse
			this.ServeJson()
		}
	} else {
		// Use current session user_id
		machines, err := this.getUserMachines(this.GetSession("user_id").(int))
		if err != nil {
			errResponse := ErrorResponse{"error", "No machines available"}
			this.Data["json"] = &errResponse
			this.ServeJson()
		}
		response.Machines = machines
	}
	// Respond
	this.Data["json"] = &response
	this.ServeJson()
}

func (this *MachinesController) requestHasUserId() (int, bool) {
	beego.Info("Checking for user_id request variable")
	var isUserIdSet bool = false
	userId, err := this.GetInt("user_id")
	if err == nil {
		beego.Info("Found", userId)
		isUserIdSet = true
	} else {
		beego.Error(err)
		beego.Info("Not found")
	}
	return int(userId), isUserIdSet
}

func (this *MachinesController) getUserMachines(userId int) ([]PublicMachine, error) {
	o := orm.NewOrm()
	var machines []models.Machine
	num, err := o.Raw("SELECT * FROM machine INNER JOIN permission ON machine.id = permission.machine_id WHERE permission.user_id = ?",
		userId).QueryRows(&machines)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	beego.Info("Got ", num, "machines")

	// Interpret machines as PublicMachine
	var pubMachines []PublicMachine
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
			}
		}
		// Fill public machine struct for output
		machine := PublicMachine{
			Id:               machines[i].Id,
			Name:             machines[i].Name,
			Description:      machines[i].Description,
			Price:            price,
			PriceUnit:        priceUnit,
			PriceCurrency:    "€", // TODO: add price currency table
			Status:           status,
			OccupiedByUserId: occupiedByUserId}
		// Append to array
		pubMachines = append(pubMachines, machine)
	}
	return pubMachines, nil
}

func (this *MachinesController) getActivation(machineId int) (*models.Activation, error) {
	o := orm.NewOrm()
	activationModel := new(models.Activation)
	beego.Info("Attempt to get activation for machine ID", machineId)
	err := o.Raw("SELECT * FROM activation WHERE machine_id = ? AND active = 1",
		machineId).QueryRow(activationModel)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return activationModel, nil
}
