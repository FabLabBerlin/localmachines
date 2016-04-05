package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
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
	uls, err := user_locations.GetAllForUser(uid)
	if err != nil {
		beego.Error("get user locations:", err)
		c.CustomAbort(500, "Cannot get user locations")
	}
	for _, ul := range uls {
		ul.Location.ClearPrivateData()
	}
	c.Data["json"] = uls
	c.ServeJSON()
}

// @Title PostUserLocation
// @Description Create new user location
// @Param	uid		path 	int	true		"User ID"
// @Param	lid		path 	int	true		"Location ID"
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:uid/locations/:lid [post]
func (c *UserLocationsController) PostUserLocation() {
	routeUid, _ := c.GetRouteUid()
	lid, err := c.GetInt64(":lid")
	if err != nil {
		beego.Error("get int:", err)
		c.CustomAbort(400, "Client Error")
	}
	sessionUid, err := c.GetSessionUserId()
	if err != nil {
		beego.Error("get session user id:", err)
		c.CustomAbort(500, "Internal Server Error")
	}
	if c.IsAdminAt(lid) || (c.IsLogged() && routeUid == sessionUid) {
		ul := user_locations.UserLocation{
			UserId:     routeUid,
			LocationId: lid,
			UserRole:   user_roles.MEMBER.String(),
		}
		if _, err := user_locations.Create(&ul); err != nil {
			beego.Error("create:", err)
			c.CustomAbort(500, "Internal Server Error")
		}
		c.ServeJSON()
	} else {
		c.CustomAbort(401, "Not authorized")
	}
}

// @Title PutUserLocation
// @Description Update existing user location
// @Param	uid		path 	int	true		"User ID"
// @Param	lid		path 	int	true		"Location ID"
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:uid/locations/:lid [put]
func (c *UserLocationsController) PutUserLocation() {
	lid, err := c.GetInt64(":lid")
	if err != nil {
		beego.Error("get int:", err)
		c.CustomAbort(400, "Client Error")
	}
	if lid <= 0 {
		beego.Error("lid:", err)
		c.CustomAbort(400, "Bad request")
	}
	dec := json.NewDecoder(c.Ctx.Request.Body)
	defer c.Ctx.Request.Body.Close()
	var ul user_locations.UserLocation
	if err := dec.Decode(&ul); err != nil {
		beego.Error("json decode:", err)
		c.CustomAbort(400, "Wrong data")
	}
	if ul.LocationId != lid {
		beego.Error("Mustn't change the user location")
		c.CustomAbort(401, "Not authorized")
	}
	if !c.IsAdminAt(ul.LocationId) {
		c.CustomAbort(401, "Not authorized")
	}
	if err := ul.Update(); err != nil {
		beego.Error("update:", err)
		c.CustomAbort(500, "Internal Server Error")
	}
	c.ServeJSON()
}

// @Title DeleteUserLocation
// @Description Delete user location
// @Param	uid		path 	int	true		"User ID"
// @Param	lid	path	int	true			"Location ID"
// @Success 200
// @Failure	400 Client Error
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:uid/locations/:lid [delete]
func (c *UserLocationsController) DeleteUserLocation() {
	uid, _ := c.GetRouteUid()
	lid, err := c.GetInt64(":lid")
	if err != nil {
		beego.Error("get int:", err)
		c.Abort("400")
	}

	if !c.IsAdminAt(lid) {
		c.Abort("401")
	}

	err = user_locations.Delete(uid, lid)
	if err != nil {
		beego.Error("Failed to delete user location")
		c.Abort("500")
	}

	c.Data["json"] = "ok"
	c.ServeJSON()
}
