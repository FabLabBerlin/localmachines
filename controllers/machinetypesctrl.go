package controllers

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/astaxie/beego"
)

type MachineTypeController struct {
	Controller
}

// @Title GetAll
// @Description Get all machine_types
// @Success 200 {object} models.locations.Location
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

// @Title Create
// @Description Create machine tpye
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [post]
func (c *MachineTypeController) Create() {
	t := machine.Type{
		Name: c.GetString("name"),
	}

	if err := t.Create(); err != nil {
		beego.Error(err)
		c.CustomAbort(500, "Failed to save")
	}

	c.ServeJSON()
}

// @Title Put
// @Description Update machine tpye with id
// @Success	200	ok
// @Failure	400	Variable message
// @Failure	401	Unauthorized
// @Failure	403	Variable message
// @router /:id [put]
func (this *MachineTypeController) Put() {

	dec := json.NewDecoder(this.Ctx.Request.Body)
	t := machine.Type{}
	if err := dec.Decode(&t); err == nil {
		beego.Info("t: ", t)
	} else {
		beego.Error("Failed to decode json", err)
		this.Fail(400, "Failed to decode json")
	}

	if !this.IsSuperAdmin() {
		beego.Error("Not authorized to update a machine type")
		this.Fail(403)
	}

	if err := t.Update(); err != nil {
		this.Fail("Failed to update:", err)
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}
