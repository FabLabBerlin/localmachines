package controllers

import (
	"github.com/FabLabBerlin/localmachines/models/machine"
)

type MachineTypeController struct {
	Controller
}

// @Title GetAll
// @Description Get all machine_types
// @Success 200 {object} locations.Location
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [get]
func (c *MachineTypeController) GetAll() {
	ls, err := machine.GetAllTypes()
	if err != nil {
		c.CustomAbort(500, "Failed to get all locations")
	}
	c.Data["json"] = ls
	c.ServeJSON()
}
