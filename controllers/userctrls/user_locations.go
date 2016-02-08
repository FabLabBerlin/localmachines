package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/astaxie/beego"
)

type UserLocationsController struct {
	Controller
}

// @Title GetUserLocations
// @Description Get user locations
// @Param	uid		path 	int	true		"User ID"
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:uid/locations [get]
func (c *UserLocationsController) GetUserLocations() {
	uid, authorized := c.GetRouteUid()
	if !authorized {
		c.CustomAbort(400, "Wrong uid in url or not authorized")
	}
	var err error
	ls, err := locations.GetAll()
	if err != nil {
		beego.Error("locations:", err)
		c.CustomAbort(500, "Cannot get user locations")
	}
	uls, err := user_locations.GetAllForUser(uid)
	if err != nil {
		beego.Error("get user locations:", err)
		c.CustomAbort(500, "Cannot get user locations")
	}
	ulsById := make(map[int64]*user_locations.UserLocation)
	for _, ul := range uls {
		ulsById[ul.LocationId] = ul
	}
	for _, l := range ls {
		if _, ok := ulsById[l.Id]; !ok {
			emptyUl := &user_locations.UserLocation{
				LocationId: l.Id,
				Location:   l,
				UserId:     uid,
				UserRole:   models.NOT_AFFILIATED,
			}
			uls = append(uls, emptyUl)
		}
	}
	c.Data["json"] = uls
	c.ServeJSON()
}

// @Title PutUserLocation
// @Description Put user location
// @Param	uid		path 	int	true		"User ID"
// @Param	lid		path 	int	true		"Location ID"
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:uid/locations/:lid [put]
func (c *UserLocationsController) PutUserLocation() {
	if !c.IsAdmin() {
		c.CustomAbort(401, "Not authorized")
	}
	dec := json.NewDecoder(c.Ctx.Request.Body)
	defer c.Ctx.Request.Body.Close()
	var ul user_locations.UserLocation
	if err := dec.Decode(&ul); err != nil {
		beego.Error("json decode:", err)
		c.CustomAbort(400, "Wrong data")
	}
	if ul.Id == 0 {
		if _, err := user_locations.Create(&ul); err != nil {
			beego.Error("create:", err)
			c.CustomAbort(500, "Internal Server Error")
		}
	} else if err := ul.Update(); err != nil {
		beego.Error("update:", err)
		c.CustomAbort(500, "Internal Server Error")
	}
	c.ServeJSON()
}
