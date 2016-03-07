package userctrls

import (
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
)

type ForgotPassword struct {
	Controller
}

// @Title ForgotPassword
// @Description Create new user location
// @Param	email		body 	string	true		"Email"
// @Success 200
// @Failure	400	Not authorized
// @Failure	500	Internal Server Error
// @router /forgot_password [post]
func (c *ForgotPassword) ForgotPassword() {
	err := users.AuthForgotPassword(c.GetString("email"))
	if err != nil {
		beego.Error("forgot password:", err)
		// No feedback to user about this!
	}
	c.ServeJSON()
}
