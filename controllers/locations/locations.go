// /api/locations
package locations

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
	"strings"
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

	if err := l.Create(); err != nil {
		beego.Error(err)
		this.CustomAbort(500, "Failed to save host")
	}

	this.Data = map[interface{}]interface{}{
		"Id": l.Id,
	}
	this.ServeJSON()
}

// @Title Put
// @Description Update location with id
// @Param	id		path 	int					true	"Location ID"
// @Param	body	body	locations.Location	true	"Location model"
// @Success	200	ok
// @Failure	400	Variable message
// @Failure	401	Unauthorized
// @Failure	403	Variable message
// @router /:id [put]
func (this *Controller) Put() {

	dec := json.NewDecoder(this.Ctx.Request.Body)
	l := locations.Location{}
	if err := dec.Decode(&l); err == nil {
		beego.Info("l: ", l)
	} else {
		beego.Error("Failed to decode json", err)
		this.Fail(400, "Failed to decode json")
	}

	if !this.IsSuperAdmin() {
		beego.Error("Not authorized to update a location")
		this.Fail(403)
	}

	if err := l.Update(); err != nil {
		this.Fail("Failed to update location:", err)
	}

	this.Data["json"] = "ok"
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
	if !c.IsSuperAdmin() {
		for _, l = range ls {
			l.ClearPrivateData()
		}
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

// @Title PostImage
// @Description Post machine image
// @Param	mid	path	int	true	"Machine ID"
// @Success 200 string ok
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500 Internal Server Error
// @router /:lid([0-9]+)/image [post]
func (c *Controller) PostImage() {
	locId, isLocAdmin := c.GetLocIdAdmin()
	if !isLocAdmin {
		c.CustomAbort(401, "Unauthorized")
	}

	if fmt.Sprintf("%v", locId) != c.GetString(":lid") {
		c.CustomAbort(400, "Client error")
	}

	dataUri := c.GetString("Image")

	i := strings.LastIndex(c.GetString("Filename"), ".")
	var fileExt string
	if i >= 0 {
		fileExt = c.GetString("Filename")[i:]
	} else {
		c.CustomAbort(500, "File name has no proper extension")
	}

	fn := imageFilename(c.GetString(":lid"), fileExt)
	if err := models.UploadImage(fn, dataUri); err != nil {
		beego.Error("upload image:", err)
		c.CustomAbort(500, "Internal Server Error")
	}

	l, err := locations.Get(locId)
	if err != nil {
		beego.Error("get location:", err)
		c.CustomAbort(500, "Internal Server Error")
	}

	l.Logo = fn
	if err = l.Update(); err != nil {
		beego.Error("Failed updating location:", err)
		c.CustomAbort(500, "Failed to update location")
	}

	c.ServeJSON()
}

func imageFilename(locationId string, fileExt string) string {
	return "location-logo-" + locationId + fileExt
}

// @Title Debug
// @Description Post debug infos, everybody can do it
// @Param	location 	query 	int		true	"Location"
// @Param	message		body 	string	true	"Message"
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /debug [post]
func (c *Controller) Debug() {
	locId := c.GetString("location")
	ip := c.ClientIp()
	msg := c.GetString("message")
	beego.Info("[locId="+locId, "ip=", ip, "]:", msg)
	c.Ctx.WriteString("")
	c.Finish()
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

// @Title Jabber Connect
// @Description Return server's Jabber ID
// @Param	lid	path 	int	true	"Location"
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /jabber_connect [get]
func (c *Controller) JabberConnect() {
	beego.Info("Requesting Server Jabber ID: location=", c.GetString("location"))
	c.Ctx.WriteString(beego.AppConfig.String("XmppUser"))
	c.Finish()
}
