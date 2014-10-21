package controllers

import (
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

type LoginResponse struct {
	SessionToken string
}

func (this *LoginController) Post() {
	response := &LoginResponse{"63b88a8f355f5d83844ac327d28cc25a"}
	this.Data["json"] = &response
	this.ServeJson()
}
