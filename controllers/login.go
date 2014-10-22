package controllers

import (
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Login() {
	response := struct{ SessionToken string }{"63b88a8f355f5d83844ac327d28cc25a"}
	this.Data["json"] = &response
	this.ServeJson()
}
