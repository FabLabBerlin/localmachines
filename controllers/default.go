package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	if beego.AppConfig.String("runmode") == "dev" {
		this.Data["AppTitle"] = "Fabsmith Dev"
	} else {
		this.Data["AppTitle"] = "Fabsmith"
	}
	this.Data["AppDescription"] = "Fabsmith - the fab lab locksmith"
	this.TplNames = "index.html"
	this.Render()
}
