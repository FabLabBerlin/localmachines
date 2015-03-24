package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

// Redirect to default html interface
func (this *MainController) Get() {

	// Check for app config value
	redirectUrl := beego.AppConfig.String("redirecturl")
	if redirectUrl != "" {
		beego.Trace("Redirect URL:", redirectUrl)
	} else {
		redirectUrl = "machines/"
	}

	beego.Info("Redirecting...")
	this.Redirect(redirectUrl, 302)
}
