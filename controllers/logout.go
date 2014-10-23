package controllers

import (
	"github.com/astaxie/beego"
)

type LogoutController struct {
	beego.Controller
}

// Logout user, handle API /logout request
func (this *LogoutController) Logout() {
	sessUsername := this.GetSession("username")
	this.DestroySession()
	if sessUsername == nil {
		beego.Info("User not logged in")
	} else {
		beego.Info("Logged out user", sessUsername)
	}
	// Respond
	response := struct{ Status string }{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}
