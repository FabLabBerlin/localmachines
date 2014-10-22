package controllers

import (
	"github.com/astaxie/beego"
)

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Logout() {
	this.DestroySession()
	response := struct{ Status string }{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}
