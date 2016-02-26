package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

type DebugController struct {
	Controller
}

// @Title Get
// @Description Get debug route
// @Param	aid		path 	int	true		"Activation ID"
// @Success 200 {object}
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [get]
func (this *DebugController) Get() {
	if !this.IsSuperAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}
	beego.Error("cookies:", this.Ctx.Request.Cookies())
	json, err := json.Marshal(this.Ctx.Request.Header)
	if err != nil {
		beego.Error("marshal json:", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	beego.Error("headers:", string(json))
	beego.Error("host:", this.Ctx.Request.Host)
	beego.Error("remote addr:", this.Ctx.Request.RemoteAddr)
	beego.Error("request uri:", this.Ctx.Request.RequestURI)
	this.ServeJSON()
}
