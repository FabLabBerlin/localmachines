// Fabsmith controllers package
package controllers

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

// Root container for all fabsmith controllers.
// Contains common functions
type Controller struct {
	beego.Controller
}

type StatusResponse struct {
	Status  string
	Message string
}

type CreatedResponse struct {
	Status string
	Id     int
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

// Makes use of SimpleResponse struct, responds with simple message JSON
func (this *Controller) serveStatusResponse(status string, message string) {
	this.Data["json"] = &StatusResponse{status, message}
	this.ServeJson() // Exit!
}

// Serve created response with created element ID
func (this *Controller) serveCreatedResponse(id int) {
	this.Data["json"] = &CreatedResponse{"created", id}
	this.ServeJson()
}

func (this *Controller) serveSuccessResponse() {
	this.Data["json"] = &struct{ Status string }{"ok"}
	this.ServeJson()
}

// Checks if request has user_id variable set and returns it
func (this *Controller) requestHasUserId() (bool, int) {
	beego.Info("Checking for user_id request variable")
	var isUserIdSet bool = false
	userId, err := this.GetInt("user_id")
	if err == nil {
		beego.Info("Found", userId)
		isUserIdSet = true
	} else {
		beego.Error(err)
		beego.Info("Not found")
	}
	return isUserIdSet, int(userId)
}

// Checks if request variable machine_id is set and returns it
func (this *Controller) requestHasMachineId() (bool, int) {
	beego.Info("Checking for machine_id request variable")
	var isMachineIdSet bool = false
	machineId, err := this.GetInt("machine_id")
	if err == nil {
		beego.Info("Machine ID", machineId)
		isMachineIdSet = true
	} else {
		beego.Error(err)
		beego.Info("Machine ID not found")
	}
	return isMachineIdSet, int(machineId)
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

// Get user data by id
func (this *Controller) getUser(userId int) (*models.User, error) {
	o := orm.NewOrm()
	userModel := models.User{Id: userId}
	err := o.Read(&userModel)
	if err != nil {
		beego.Error(err)
	}
	return &userModel, err
}

// Return true if user is admin
func (this *Controller) isAdmin() bool {
	rolesModel := this.getSessionUserRoles()
	if rolesModel.Admin {
		return true
	}
	return false
}

// Return true if user is staff
func (this *Controller) isStaff() bool {
	rolesModel := this.getSessionUserRoles()
	if rolesModel.Staff {
		return true
	}
	return false
}

// Return session user ID and perform relevant checks
func (this *Controller) getSessionUserId() (int, error) {
	userId := this.GetSession("user_id")
	if userId == nil {
		beego.Error("Could not get session user ID")
		return int(0), errors.New("Could not get session user ID")
	}
	return userId.(int), nil
}
