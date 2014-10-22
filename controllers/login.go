package controllers

import (
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Login() {
	username := this.GetSession("username")

	if username == nil {
		// TODO: check username and password
		if this.isUserValid() {
			this.SetSession("username", "kris")
			response := struct{ Status string }{"ok"}
			this.Data["json"] = &response
		} else {
			response := struct {
				Status  string
				Message string
			}{"error",
				"Invalid username or password"}
			this.Data["json"] = &response
		}
	} else {
		response := struct{ Status string }{"logged"}
		this.Data["json"] = &response
	}

	this.ServeJson()
}

func (this *LoginController) isUserValid() bool {
	username := this.GetString("username")
	password := this.GetString("password")

	if username == "kris" &&
		password == "bfd59291e825b5f2bbf1eb76569f8fe7" {
		return true
	}

	return false
}
