package controllers

import (
	"github.com/astaxie/beego"
)

type MachinesMainController struct {
	beego.Controller
}

func (this *MachinesMainController) Get() {
	if beego.AppConfig.String("runmode") == "dev" {
		this.Data["AppTitle"] = "Fabsmith Machine Interface Dev"
	} else {
		this.Data["AppTitle"] = "Fabsmith Machine Interface"
	}

	this.Data["AppDescription"] = "Authenticate and activate your machines."
	this.TplNames = "machines.html"
	this.Render()
}
