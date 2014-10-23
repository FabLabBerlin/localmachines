package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

type LoginController struct {
	beego.Controller
}

// Log in user, handle API /login request
func (this *LoginController) Login() {
	// Attempt to get stored session username
	sessUsername := this.GetSession("username")
	if sessUsername == nil {
		// If not set, user is not logged in
		if this.isUserValid() {
			// If user is valid, log in, save username in session
			reqUsername := this.GetString("username")
			this.SetSession("username", reqUsername)
			// Save the user ID in session as well
			userId := this.getUserId(reqUsername)
			this.SetSession("user_id", userId)
			// Output JSON
			response := struct{ Status string }{"ok"}
			this.Data["json"] = &response
			beego.Info("User", reqUsername, "successfully logged in")
		} else {
			// If user not valid, respond with error
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

func (this *LoginController) getUserId(username string) int {
	o := orm.NewOrm()
	userModel := new(models.User)
	beego.Info("Attempt to get user id for username ", username)
	err := o.Raw("SELECT id FROM user WHERE username = ?", username).QueryRow(userModel)
	if err != nil {
		beego.Error(err)
	}
	return userModel.Id
}

func (this *LoginController) getPassword(username string) string {
	o := orm.NewOrm()
	authModel := new(models.Auth)
	beego.Info("Attempt to get password from auth table for username", username)
	err := o.Raw("SELECT password FROM auth INNER JOIN user ON auth.user_id = user.id WHERE user.username = ?",
		username).QueryRow(&authModel)
	if err != nil {
		beego.Error(err)
	}
	return authModel.Password
}

func (this *LoginController) isUserValid() bool {
	// Get request variables
	username := this.GetString("username")
	password := this.GetString("password")
	// Get password from DB
	storedUserPassword := this.getPassword(username)
	// Check if passwords match
	if password == storedUserPassword {
		return true
	}
	return false
}
