// /api/locations
package locations

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/astaxie/beego"
)

type Controller struct {
	controllers.Controller
}

func (this *Controller) GetRouteLid() (lid int64, userRole user_roles.Role) {
	// Check if logged in
	suid, err := this.GetSessionUserId()
	if err != nil {
		beego.Info("Not logged in")
		return 0, user_roles.NOT_AFFILIATED
	}

	// Get requested location ID
	lid, err = this.GetInt64(":lid")
	if err != nil {
		beego.Error("Failed to get :lid", err)
		return 0, user_roles.NOT_AFFILIATED
	}

	uls, err := user_locations.GetAllForUser(suid)
	if err != nil {
		return 0, user_roles.NOT_AFFILIATED
	}

	for _, ul := range uls {
		if ul.LocationId == lid {
			return lid, ul.GetRole()
		}
	}

	return 0, user_roles.NOT_AFFILIATED
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
// @router /:lid [get]
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
