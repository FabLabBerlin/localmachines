package userctrls

import (
	"bytes"
	"github.com/FabLabBerlin/localmachines/lib/email"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"text/template"
)

// Inspired by Troy Hunt's write up about the topic:
// http://bit.ly/1CNpHDH
type ForgotPassword struct {
	Controller
}

const FORGOT_PASSWORD_TEMPLATE = `
Hello,


This is an automatically generated mail from the EASY LAB system.


{{.Message}}


If you have any questions, don't hesitate to drop us a mail:
support@fablab.berlin


Greetings, EASY LAB
`

const FORGOT_PASSWORD_STANDARD_TEXT = `
Someone has used the password reset function with your E-Mail
address, hopefully this was you. :)
`

const FORGOT_PASSWORD_WRONG_TEXT = `
Someone has used the password reset function with your E-Mail
address, although this E-Mail address is not registered in the
system.  So nothing happened basically. :)
`

var forgotPasswordTemplate *template.Template

func init() {
	forgotPasswordTemplate = template.Must(template.New("Forgot pw template").
		Parse(FORGOT_PASSWORD_TEMPLATE))
}

// @Title ForgotPassword
// @Description Send forgot password E-Mail
// @Param	email		body 	string	true		"Email"
// @Success 200
// @Failure	400	Not authorized
// @Failure	500	Internal Server Error
// @router /forgot_password [post]
func (c *ForgotPassword) ForgotPassword() {
	mail := email.New()
	addr := c.GetString("email")
	subject := "EASY LAB password recovery"

	key, err := users.AuthForgotPassword(addr)
	var message string
	if err == nil {
		message = FORGOT_PASSWORD_STANDARD_TEXT
		message += "\n\n"
		message += "Please follow this link to proceed:\n\n"
		message += "https://lab.fablab.berlin/machines/#/forgot_password/recover?key=" + key
		message += "\n"
	} else {
		beego.Error("forgot password:", err)
		// No feedback to user about this, to confuse people
		// that try to find accounts.
		message = FORGOT_PASSWORD_WRONG_TEXT
	}
	buf := bytes.NewBufferString("")
	beego.Info("email/to:", addr)
	beego.Info("subject:", subject)
	beego.Info("message:", message)
	err = forgotPasswordTemplate.Execute(buf, map[string]interface{}{
		"Message": message,
	})
	if err != nil {
		beego.Error("Error executing forgot password mail template:", err)
	}
	if err := mail.Send(addr, subject, message); err != nil {
		beego.Error("Error sending wrong forgot password mail:", err)
	}
	c.ServeJSON()
}

// @Title CheckPhone
// @Description Check phone number
// @Param	key			body 	string	true		"Reset Password Key"
// @Param	phone		body 	string	true		"Phone"
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /forgot_password/phone [post]
func (c *ForgotPassword) CheckPhone() {
	key := c.GetString("key")
	phone := c.GetString("phone")

	_, err := users.AuthCheckPhone(key, phone)
	if err == users.ErrAuthWrongKey ||
		err == users.ErrAuthOutdatedKey ||
		err == users.ErrAuthWrongPhone {
		c.CustomAbort(401, err.Error())
	} else if err != nil {
		beego.Error("auth check phone:", err)
		c.CustomAbort(500, "Internal Server Error")
	}
	c.ServeJSON()
}

// @Title ResetPassword
// @Description Reset password
// @Param	key			body 	string	true		"Reset Password Key"
// @Param	password	body 	string	true		"New password"
// @Param	phone		body 	string	true		"Phone"
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /forgot_password/phone [post]
func (c *ForgotPassword) ResetPassword() {
	key := c.GetString("key")
	phone := c.GetString("phone")
	pass := c.GetString("password")

	uid, err := users.AuthCheckPhone(key, phone)
	if err == users.ErrAuthWrongKey ||
		err == users.ErrAuthOutdatedKey ||
		err == users.ErrAuthWrongPhone {
		c.CustomAbort(401, err.Error())
	} else if err != nil {
		beego.Error("auth check phone:", err)
		c.CustomAbort(500, "Internal Server Error")
	}

	if err := users.AuthSetPassword(uid, pass); err != nil {
		beego.Error("set pw of ", uid, ":", err)
		c.CustomAbort(500, "Internal Server Error")
	}

	c.ServeJSON()
}
