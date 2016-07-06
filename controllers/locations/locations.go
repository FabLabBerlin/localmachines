// /api/locations
package locations

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
)

var runMode = beego.AppConfig.String("RunMode")

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
	var l *locations.Location
	for _, l = range ls {
		l.ClearPrivateData()
	}
	c.Data["json"] = ls
	c.ServeJSON()
}

// @Title Get
// @Description Get location
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:lid([0-9]+) [get]
func (c *Controller) Get() {
	locId, isLocAdmin := c.GetLocIdAdmin()

	l, err := locations.Get(locId)
	if err != nil {
		beego.Error("get location:", err)
		c.CustomAbort(500, "Internal Server Error")
	}
	if !isLocAdmin {
		l.ClearPrivateData()
	}
	c.Data["json"] = l
	c.ServeJSON()
}

// @Title MyIp
// @Description Returns client's IP address as string
// @Param	lid	path 	int	true	"Location"
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /my_ip [get]
func (c *Controller) MyIp() {
	beego.Info("Requesting MyIp: location=", c.GetString("location"))
	c.Ctx.WriteString(c.ClientIp())
	c.Finish()
}
