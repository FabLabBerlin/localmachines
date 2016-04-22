// /api/coupons
package coupons

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models/coupons"
	"github.com/astaxie/beego"
)

type Controller struct {
	controllers.Controller
}

// @Title GetAll
// @Description Get all coupons
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [get]
func (this *Controller) GetAll() {
	locId, isLocAdmin := this.GetLocIdAdmin()
	if !isLocAdmin {
		this.Abort("401")
	}
	if locId <= 0 {
		this.Abort("400")
	}
	cs, err := coupons.GetAllCouponsAt(locId)
	if err != nil {
		beego.Error("get all:", err)
		this.Abort("500")
	}
	this.Data["json"] = cs
	this.ServeJSON()
}

// @Title Generate
// @Description Generate coupons
// @Param	location	query	int64	true	"Location ID"
// @Param	static_code	body	string	false	"Static code (optional)"
// @Param	n			body	int64	true	"Number of coupons"
// @Param	value		body	float64	true	"Value of coupon in ¤"
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [post]
func (this *Controller) Generate() {
	locId, isLocAdmin := this.GetLocIdAdmin()
	if !isLocAdmin {
		this.Abort("401")
	}
	if locId <= 0 {
		this.Abort("400")
	}
	n, err := this.GetInt64("n")
	if err != nil {
		beego.Error("n:", err)
		this.Abort("400")
	}
	value, err := this.GetFloat("value")
	if err != nil {
		beego.Error("value:", err)
		this.Abort("400")
	}
	staticCode := this.GetString("static_code")
	cs, err := coupons.Generate(locId, staticCode, int(n), value)
	if err != nil {
		beego.Error("generate:", err)
		this.Abort("500")
	}
	this.Data["json"] = cs
	this.ServeJSON()
}
