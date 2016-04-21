package controllers

import (
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/astaxie/beego"
)

type UserLocationsController struct {
	Controller
}

// @Title GetAll
// @Description Get all user locations
// @Param	location	query 	int	true		"User ID"
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [get]
func (c *UserLocationsController) GetAll() {
	locId, authorized := c.GetLocIdAdmin()
	if !authorized {
		c.CustomAbort(401, "Not authorized")
	}
	uls, err := user_locations.GetAllForLocation(locId)
	if err != nil {
		beego.Error("get user locations:", err)
		c.CustomAbort(500, "Cannot get user locations")
	}
	c.Data["json"] = uls
	c.ServeJSON()
}