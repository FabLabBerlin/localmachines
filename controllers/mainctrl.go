package controllers

import (
	"errors"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
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

func (this *Controller) GetSessionUserId() (int64, error) {
	tmp := this.GetSession(SESSION_USER_ID)
	if sid, ok := tmp.(int64); ok {
		browser := this.GetSession(SESSION_BROWSER)
		if browser != this.Ctx.Input.UserAgent() {
			beego.Error("GetSessionUserId: wrong browser")
			return 0, errors.New("user not correctly logged in")
		}
		ip := this.GetSession(SESSION_IP)
		if ip != this.Ctx.Input.IP() {
			beego.Error("GetSessionUserId: wrong IP")
			return 0, errors.New("user not correctly logged in")
		}
		/*accEnc := this.GetSession(SESSION_ACCEPT_ENCODING)
		if h := this.Ctx.Input.Header("Accept-Encoding"); accEnc != h {
			beego.Error("GetSessionUserId: wrong Accept-Encoding:", accEnc, "vs", h)
			return 0, errors.New("user not correctly logged in")
		}*/
		accLang := this.GetSession(SESSION_ACCEPT_LANGUAGE)
		if accLang != this.Ctx.Input.Header("Accept-Language") {
			beego.Error("GetSessionUserId: wrong Accept-Language")
			return 0, errors.New("user not correctly logged in")
		}
		return sid, nil
	} else {
		return 0, errors.New("User not logged in")
	}
}

func (this *Controller) GetSessionLocationId() int64 {
	return this.GetSession(SESSION_LOCATION_ID).(int64)
}

func (this *Controller) SetLogged(username string, userId int64, locationId int64) {
	this.SetSession(SESSION_USERNAME, username)
	this.SetSession(SESSION_USER_ID, userId)
	this.SetSession(SESSION_LOCATION_ID, locationId)
	this.SetSession(SESSION_BROWSER, this.Ctx.Input.UserAgent())
	this.SetSession(SESSION_IP, this.Ctx.Input.IP())
	//this.SetSession(SESSION_ACCEPT_ENCODING, this.Ctx.Input.Header("Accept-Encoding"))
	this.SetSession(SESSION_ACCEPT_LANGUAGE, this.Ctx.Input.Header("Accept-Language"))
}

func (this *Controller) IsLogged() bool {
	_, err := this.GetSessionUserId()
	return err == nil
}

// Return true if user is admin, if no args are passed, uses session user ID,
// if single user ID is passed, checks the passed one. Fails otherwise.
func (this *Controller) IsAdmin(userIds ...int64) bool {
	var userId int64
	var err error
	if len(userIds) == 0 {
		userId, err = this.GetSessionUserId()
		if err != nil {
			return false
		}
	} else if len(userIds) == 1 {
		userId = userIds[0]
	} else {
		beego.Error("Expecting single or no value as input")
		return false
	}
	user, err := models.GetUser(userId)
	if err != nil {
		return false
	}
	return user.GetRole() == user_roles.ADMIN || user.GetRole() == user_roles.SUPER_ADMIN
}

// Return true if user is staff, if no args are passed, uses session user ID,
// if single user ID is passed, checks the passed one. Fails otherwise.
func (this *Controller) IsStaff(userIds ...int64) bool {
	var userId int64
	var err error
	if len(userIds) == 0 {
		userId, err = this.GetSessionUserId()
		if err != nil {
			return false
		}
	} else if len(userIds) == 1 {
		userId = userIds[0]
	} else {
		beego.Error("Expecting single or no value as input")
		return false
	}
	var user *models.User
	user, err = models.GetUser(userId)
	if err != nil {
		return false
	}
	return user.GetRole() == user_roles.STAFF
}
