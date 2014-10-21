package controllers

import (
	"github.com/astaxie/beego"
)

type LogoutController struct {
	beego.Controller
}

type LogoutResponse struct {
	Status string
}

func (this *LogoutController) Any() {
	response := &LogoutResponse{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}
