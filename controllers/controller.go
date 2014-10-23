// Fabsmith controllers package
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

// Root container for all fabsmith controllers.
// Contains common functions
type Controller struct {
	beego.Controller
}

// Checks if user is logged in before sending out any data
func (this *Controller) Prepare() {
	sessUser := this.GetSession("username")
	if sessUser == nil {
		response := struct {
			Status  string
			Message string
		}{"error", "Not logged in"}
		this.Data["json"] = &response
		this.ServeJson()
	}
}

// Returns user roles model for the currently logged in user
func (this *Controller) getSessionUserRoles() models.UserRoles {
	o := orm.NewOrm()
	sessUserId := this.GetSession("user_id").(int)
	rolesModel := models.UserRoles{UserId: sessUserId}
	beego.Info("Attempt to get user roles for user id", sessUserId)
	err := o.Read(&rolesModel)
	if err != nil {
		beego.Error(err)
	}
	return rolesModel
}

// Returns user model for the currently logged in user
func (this *Controller) getSessionUserData() *models.User {
	sessUserId := this.GetSession("user_id").(int)
	o := orm.NewOrm()
	userModel := new(models.User)
	beego.Info("Attempt to get user row for user id", sessUserId)
	err := o.Raw("SELECT * FROM user WHERE id = ?",
		sessUserId).QueryRow(&userModel)
	if err != nil {
		beego.Error(err)
	}
	return userModel
}
