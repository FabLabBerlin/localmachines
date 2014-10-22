package controllers

import (
	"github.com/astaxie/beego"
)

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Logout() {
	response := struct{ Status string }{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}
