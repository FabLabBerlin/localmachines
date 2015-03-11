package controllers

import (
	"github.com/astaxie/beego"
)

type MachinesMainController struct {
	beego.Controller
}

// Set runmode specific template variables and serve machines template
func (this *MachinesMainController) Get() {
	if beego.AppConfig.String("runmode") == "dev" {
		this.Data["AppTitle"] = "Fabsmith Machine Interface Dev"
		this.TplNames = "dev/machines/index.html"
	} else {
		this.Data["AppTitle"] = "Fabsmith Machine Interface"
		this.TplNames = "prod/machines/index.html"
	}

	this.Data["AppDescription"] = "Authenticate and activate your machines."
	this.Render()
}
