package controllers

import (
	"github.com/astaxie/beego"
)

type BackOfficeMainController struct {
	beego.Controller
}

func (this *BackOfficeMainController) Get() {
	if beego.AppConfig.String("runmode") == "dev" {
		this.Data["AppTitle"] = "Fabsmith Back Office Dev"
	} else {
		this.Data["AppTitle"] = "Fabsmith Back Office"
	}
	this.Data["AppDescription"] = "The Fabsmith Back Office let's you do things..."
	this.TplNames = "backoffice.html"
	this.Render()
}
