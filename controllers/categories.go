package controllers

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models/categories"
	"github.com/astaxie/beego"
)

type CategoriesController struct {
	Controller
}

// @Title GetAll
// @Description Get all categories
// @Success 200 {object} models.locations.Location
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [get]
func (cc *CategoriesController) GetAll() {
	locId, authorized := cc.GetLocIdMember()
	if !authorized {
		cc.Fail(401, "Not authorized")
	}

	cs, err := categories.GetAll(locId)
	if err != nil {
		cc.Fail(500, "Failed to get all locations")
	}
	cc.Data["json"] = cs
	cc.ServeJSON()
}

// @Title Create
// @Description Create machine tpye
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [post]
func (cc *CategoriesController) Create() {
	locId, authorized := cc.GetLocIdAdmin()
	if !authorized {
		cc.Fail(401, "Not authorized")
	}

	c := categories.Category{
		LocationId: locId,
		Name:       cc.GetString("name"),
	}

	if err := c.Create(); err != nil {
		beego.Error(err)
		cc.Fail(500, "Failed to save")
	}

	cc.ServeJSON()
}

// @Title Put
// @Description Update machine tpye with id
// @Success	200	ok
// @Failure	400	Variable message
// @Failure	401	Unauthorized
// @Failure	403	Variable message
// @router /:id [put]
func (cc *CategoriesController) Put() {
	locId, authorized := cc.GetLocIdAdmin()
	if !authorized {
		cc.Fail(401, "Not authorized")
	}

	dec := json.NewDecoder(cc.Ctx.Request.Body)
	c := categories.Category{}
	if err := dec.Decode(&c); err == nil {
		beego.Info("c: ", c)
	} else {
		beego.Error("Failed to decode json", err)
		cc.Fail(400, "Failed to decode json")
	}

	if !cc.IsAdminAt(locId) || c.LocationId != locId {
		beego.Error("Not authorized to update cc machine category")
		cc.Fail(403)
	}

	if err := c.Update(); err != nil {
		cc.Fail("Failed to update:", err)
	}

	cc.Data["json"] = "ok"
	cc.ServeJSON()
}

// @Title Archive Machine Type
// @Description Archive a machine category
// @Param	id	query	int	true	"Machine category ID"
// @Success 200 string
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:id/archive [put]
func (cc *CategoriesController) Archive() {
	locId, authorized := cc.GetLocIdAdmin()
	if !authorized {
		cc.Fail(401, "Not authorized")
	}

	id, err := cc.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get :id variable")
		cc.Fail("400")
	}

	c, err := categories.Get(id)
	if err != nil {
		beego.Error("Failed to get category")
		cc.Fail("500")
	}

	if !cc.IsAdminAt(locId) || c.LocationId != locId {
		beego.Error("Unauthorized attempt to archive machine category")
		cc.Fail("401")
	}

	err = c.Archive()
	if err != nil {
		cc.Fail(500, "Failed to archive")
	}

	cc.ServeJSON()
}

// @Title Unarchive Machine category
// @Description Unarchive a machine category
// @Param	id	query	int	true	"Machine category ID"
// @Success 200 string
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:id/unarchive [put]
func (cc *CategoriesController) Unarchive() {
	locId, authorized := cc.GetLocIdAdmin()
	if !authorized {
		cc.Fail(401, "Not authorized")
	}

	id, err := cc.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get :id variable")
		cc.Fail("400")
	}

	c, err := categories.Get(id)
	if err != nil {
		beego.Error("Failed to get category")
		cc.Fail("500")
	}

	if !cc.IsAdminAt(locId) || c.LocationId != locId {
		beego.Error("Unauthorized attempt to unarchive category")
		cc.Fail("401")
	}

	err = c.Unarchive()
	if err != nil {
		cc.Fail(500, "Failed to unarchive")
	}

	cc.ServeJSON()
}
