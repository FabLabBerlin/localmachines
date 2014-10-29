package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

type LoginController struct {
	Controller
}

// Override our custom root controller's Prepare method as it is checking
// if we are logged in and we don't want that here at this point
func (this *LoginController) Prepare() {
	beego.Info("Skipping login check")
}

// Log in user, handle API /login request
func (this *LoginController) Login() {
	// Attempt to get stored session username
	sessUsername := this.GetSession(SESSION_FIELD_NAME_USERNAME)
	if sessUsername == nil {
		// If not set, user is not logged in
		if this.isUserValid() {
			// If user is valid, log in, save username in session
			reqUsername := this.GetString(REQUEST_FIELD_NAME_USERNAME)
			this.SetSession(SESSION_FIELD_NAME_USERNAME, reqUsername)
			// Save the user ID in session as well
			userId, err := this.getUserId(reqUsername)
			if err != nil {
				if beego.AppConfig.String("runmode") == "dev" {
					panic("User valid, but could not get user ID")
				}
				// This is really strange - respond with error in case of bad spirits...
				this.serveErrorResponse("Invalid username or password")
			}
			this.SetSession(SESSION_FIELD_NAME_USER_ID, userId)
			beego.Info("User", reqUsername, "successfully logged in")
			this.serveOkResponse()
		} else {
			// If user not valid, respond with error
			this.serveErrorResponse("Invalid username or password")
			beego.Info("Failed to authenticate user")
		}
	} else {
		this.Data["json"] = struct{ Status string }{"logged"}
		this.ServeJson()
		//this.serveErrorResponse("Already logged")
	}
}

func (this *LoginController) getUserId(username string) (int, error) {
	userModel := models.User{}
	userModel.Username = username
	beego.Trace("Attempt to get user id for username ", username)
	o := orm.NewOrm()
	err := o.Read(&userModel, "Username")
	if err != nil {
		beego.Error("Could not get user ID with username", username, ":", err)
		return int(0), err
	}
	return userModel.Id, nil
}

func (this *LoginController) getPassword(username string) (string, error) {
	beego.Info("Attempt to get password from auth table for username", username)
	authModel := models.Auth{}
	o := orm.NewOrm()
	err := o.Raw("SELECT password FROM auth INNER JOIN user ON auth.user_id = user.id WHERE user.username = ?",
		username).QueryRow(&authModel)
	if err != nil {
		beego.Error("Could not read into AuthModel:", err)
		return "", err
	}
	return authModel.Password, nil
}

func (this *LoginController) isUserValid() bool {
	// Get request variables
	username := this.GetString("username")
	password := this.GetString("password")

	beego.Trace("POST username:", username)
	beego.Trace("POST password:", password)

	// Get password from DB
	storedUserPassword, err := this.getPassword(username)
	if err != nil {
		beego.Error("Could not get password for user", username)
		return false
	}
	// Check if passwords match
	if password == storedUserPassword {
		return true
	}
	return false
}
