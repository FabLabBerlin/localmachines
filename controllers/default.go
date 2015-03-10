package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	if beego.AppConfig.String("runmode") == "dev" {
		this.Data["AppTitle"] = "Fabsmith Machine Interface Dev"
	} else {
		this.Data["AppTitle"] = "Fabsmith Machine Interface"
	}

	this.Data["AppDescription"] = "Authenticate and activate your machines."
	this.TplNames = "machineinterface.html"
	this.Render()
}
