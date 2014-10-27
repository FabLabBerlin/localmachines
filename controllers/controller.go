// Fabsmith controllers package, handles all the API calls
package controllers

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

// Field names for session variables
const SESSION_FIELD_NAME_USER_ID string = "user_id"

// Field names for request variables
const REQUEST_FIELD_NAME_USER_ID string = "user_id"
const REQUEST_FIELD_NAME_MACHINE_ID string = "machine_id"
const REQUEST_FIELD_NAME_ACTIVATION_ID string = "activation_id"

// Root container for all fabsmith controllers - contains common functions.
// It is used for almost every controller, except the login and logout
type Controller struct {
	beego.Controller
}

// Common status response message struct. Mostly used for
// {"status":"error", "message":"Error message"} JSON response
type ErrorResponse struct {
	Status  string
	Message string
}

// Creates new ErrorResponse instance with Status:"error" already set
func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{Status: "error"}
}

// Makes use of ErrorResponse struct, responds with simple error message JSON
func (this *Controller) serveErrorResponse(errorMessage string) {
	errorResponse := NewErrorResponse()
	errorResponse.Message = errorMessage
	this.Data["json"] = errorResponse
	this.ServeJson()
}

// Struct for simple "created" response: {"status":"created", "id":"created_id"}
type CreatedResponse struct {
	Status string
	Id     int
}

// Creates new CreatedResponse object with Status:"created" already set
func NewCreatedResponse() *CreatedResponse {
	return &CreatedResponse{Status: "created"}
}

// Serve created response with created element ID
func (this *Controller) serveCreatedResponse(createdId int) {
	createdResponse := NewCreatedResponse()
	createdResponse.Id = createdId
	this.Data["json"] = createdResponse
	this.ServeJson()
}

// Struct for simple "ok" JSON response {"status":"ok"}
type OkResponse struct {
	Status string
}

// Returns sible OkResponse instance with Status:"ok" set
func NewOkResponse() *OkResponse {
	return &OkResponse{Status: "ok"}
}

// Serves JSON OK message {"status":"ok"}
func (this *Controller) serveOkResponse() {
	okResponse := NewOkResponse()
	this.Data["json"] = okResponse
	this.ServeJson()
}

// Checks if user is logged in before sending out any data, responds with
// "Not logged in" error if user not logged in
func (this *Controller) Prepare() {
	sessUser := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if sessUser == nil {
		this.serveErrorResponse("Not logged in")
	}
}

// Return user ID
func (this *Controller) getSessionUserId() (int, error) {
	userId := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if userId == nil {
		beego.Critical("Could not get session user ID")
		return 0, errors.New("Session user ID not set")
	}
	return userId.(int), nil
}

// Get user data from database by user ID or if it is not given,
// by session user ID
func (this *Controller) getUser(userIds ...int) (interface{}, error) {
	var userId int
	var err error
	if len(userIds) == 0 {
		userId, err = this.getSessionUserId()
		if err != nil {
			return nil, err
		}
	} else if len(userIds) == 1 {
		userId = userIds[0]
	} else {
		beego.Critical("Wrong argument count")
		return nil, errors.New("Invalid argument count")
	}
	beego.Trace("Attempt to get user row for user ID", userId)
	userModel := models.User{Id: userId}
	o := orm.NewOrm()
	err = o.Read(&userModel)
	if err != nil {
		beego.Critical("Could not get user data from DB:", err)
		return nil, err
	}
	return userModel, nil
}

// Returns user roles model for the currently logged in user
func (this *Controller) getUserRoles(userIds ...int) (interface{}, error) {
	var userId int
	var err error
	if len(userIds) == 0 {
		userId, err = this.getSessionUserId()
		if err != nil {
			return nil, err
		}
	} else if len(userIds) == 1 {
		userId = userIds[0]
	} else {
		beego.Critical("Wrong argument count")
		return nil, err
	}
	beego.Trace("Attempt to get user roles for user id", userId)
	rolesModel := models.UserRoles{UserId: userId}
	o := orm.NewOrm()
	err = o.Read(&rolesModel)
	if err != nil {
		beego.Critical("Could not get roles model from DB:", err)
		return nil, err
	}
	return rolesModel, nil
}

// Return true if user is admin, if no args are passed, uses session user ID,
// if single user ID is passed, checks the passed one. Fails otherwise.
func (this *Controller) isAdmin(userIds ...int) bool {
	var userId int
	var err error
	if len(userIds) == 0 {
		userId, err = this.getSessionUserId()
		if err != nil {
			return false
		}
	} else if len(userIds) == 1 {
		userId = userIds[0]
	} else {
		beego.Critical("Expecting single or no value as input")
		return false
	}
	var rolesModelInterface interface{}
	rolesModelInterface, err = this.getUserRoles(userId)
	if err != nil {
		return false
	}
	rolesModel := rolesModelInterface.(models.UserRoles)
	if rolesModel.Admin {
		return true
	}
	return false
}

// Return true if user is staff, if no args are passed, uses session user ID,
// if single user ID is passed, checks the passed one. Fails otherwise.
func (this *Controller) isStaff(userIds ...int) bool {
	var userId int
	var err error
	if len(userIds) == 0 {
		userId, err = this.getSessionUserId()
		if err != nil {
			return false
		}
	} else if len(userIds) == 1 {
		userId = userIds[0]
	} else {
		beego.Critical("Expecting single or no value as input")
		return false
	}
	var rolesModelInterface interface{}
	rolesModelInterface, err = this.getUserRoles(userId)
	if err != nil {
		return false
	}
	rolesModel := rolesModelInterface.(models.UserRoles)
	if rolesModel.Staff {
		return true
	}
	return false
}
