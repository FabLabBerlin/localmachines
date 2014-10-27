package controllers

import (
	"github.com/astaxie/beego"
)

type LogoutController struct {
	Controller
}

// Override this, because we want to be able logout if logged in
func (this *LogoutController) Prepare() {
	beego.Info("Skipping login check")
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
	this.serveOkResponse()
}
