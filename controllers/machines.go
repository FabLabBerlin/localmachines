package controllers

import (
	"github.com/astaxie/beego"
)

type MachinesController struct {
	beego.Controller
}

func (this *MachinesController) GetMachines() {
	type Machine struct {
		Id               int
		Name             string
		Description      string
		Price            float64
		PriceUnit        string
		PriceCurrency    string
		Status           string
		OccupiedByUserId int
	}
	response := struct{ Machines []Machine }{
		[]Machine{
			Machine{1, "Laser Cutter",
				"Cuts all sorts of things",
				1.0, "minute", "€", "free", 0},
			Machine{2, "3D Printer",
				"Prints 3D objects with molten plastic",
				15.0, "hour", "€", "occupied", 2}}}
	this.Data["json"] = &response
	this.ServeJson()
}
