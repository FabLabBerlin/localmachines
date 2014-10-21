package controllers

import (
	"github.com/astaxie/beego"
)

type MachinesController struct {
	beego.Controller
}

type MachinesResponseMachine struct {
	Id               int
	Name             string
	Description      string
	Price            float64
	PriceUnit        string
	PriceCurrency    string
	Status           string
	OccupiedByUserId int
}

type MachinesResponse struct {
	Machines []MachinesResponseMachine
}

func (this *MachinesController) Get() {
	response := &MachinesResponse{
		[]MachinesResponseMachine{
			MachinesResponseMachine{1, "Laser Cutter",
				"Cuts all sorts of things",
				1.0, "minute", "€", "free", 0},
			MachinesResponseMachine{2, "3D Printer",
				"Prints 3D objects with molten plastic",
				15.0, "hour", "€", "occupied", 2}}}
	this.Data["json"] = &response
	this.ServeJson()
}
