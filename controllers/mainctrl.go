package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

// Redirect to default html interface
func (this *MainController) Get() {
	this.Redirect("machines/", 302)
}
