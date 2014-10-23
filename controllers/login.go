package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type LoginController struct {
	beego.Controller
}

// Log in user
func (this *LoginController) Login() {
	// Attempt to get stored session username
	sessUsername := this.GetSession("username")
	if sessUsername == nil {
		// If not set, user is not logged in
		if this.isUserValid() {
			// If user is valid, log in
			reqUsername := this.GetString("username")
			this.SetSession("username", reqUsername)
			response := struct{ Status string }{"ok"}
			this.Data["json"] = &response
			beego.Info("User", reqUsername, "successfully logged in")
		} else {
			// If not valid, respond with error
			response := struct {
				Status  string
				Message string
			}{"error",
				"Invalid username or password"}
			this.Data["json"] = &response
			beego.Info("Failed to authenticate user")
		}
	} else {
		response := struct{ Status string }{"logged"}
		this.Data["json"] = &response
	}
	// Serve JSON response
	this.ServeJson()
}

func (this *LoginController) isUserValid() bool {
	// Get request variables
	username := this.GetString("username")
	password := this.GetString("password")
	// Get password from DB
	userModel := new(models.User)
	modelPassword := userModel.GetPasswordForUsername(username)
	beego.Info(modelPassword)
	// Check if passwords match
	if username == username &&
		password == modelPassword {
		return true
	}
	return false
}
