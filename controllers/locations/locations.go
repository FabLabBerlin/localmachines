// /api/locations
package locations

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
)

type Controller struct {
	controllers.Controller
}

// @Title Create
// @Description Create location
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [post]
func (this *Controller) Create() {
	l := locations.Location{
		Title:        this.GetString("title"),
		FirstName:    this.GetString("first_name"),
		LastName:     this.GetString("last_name"),
		Email:        this.GetString("email"),
		City:         this.GetString("city"),
		Organization: this.GetString("organization"),
		Phone:        this.GetString("phone"),
		Comments:     this.GetString("comments"),
	}

	if err := l.Save(); err != nil {
		beego.Error(err)
		this.CustomAbort(500, "Failed to save host")
	}

	this.Data = map[interface{}]interface{}{
		"Id": l.Id,
	}
	this.ServeJSON()
}

// @Title GetAll
// @Description Get all locations
// @Success 200 {object} models.locations.Location
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [get]
func (c *Controller) GetAll() {
	ls, err := locations.GetAll()
	if err != nil {
		c.CustomAbort(500, "Failed to get all locations")
	}
	c.Data["json"] = ls
	c.ServeJSON()
}
