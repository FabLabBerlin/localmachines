package locations

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
)

type Controller struct {
	controllers.Controller
}

// @Title GetAll
// @Description Get all locations
// @Success 200 {object} locations.Location
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [get]
func (c *Controller) GetAll() {
	if !c.IsAdmin() && !c.IsStaff() {
		beego.Error("Not authorized to get all locations")
		c.CustomAbort(401, "Not authorized")
	}

	ls, err := locations.GetAll()
	if err != nil {
		c.CustomAbort(500, "Failed to get all locations")
	}
	c.Data["json"] = ls
	c.ServeJson()
}
