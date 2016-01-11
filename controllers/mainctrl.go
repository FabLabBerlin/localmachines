package controllers

import (
	"errors"
	"github.com/FabLabBerlin/localmachines/models"
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
		return sid, nil
	} else {
		return 0, errors.New("User not logged in")
	}
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
	var user *models.User
	user, err = models.GetUser(userId)
	if err != nil {
		return false
	}
	return user.UserRole == models.ADMIN
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
	return user.UserRole == models.STAFF
}
