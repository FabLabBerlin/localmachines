// Fabsmith controllers package, handles all the API calls
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

// Field names for session variables
const SESSION_FIELD_NAME_USER_ID string = "user_id"
const SESSION_FIELD_NAME_USERNAME string = "username"

// Field names for request variables
const REQUEST_FIELD_NAME_USER_ID string = "user_id"
const REQUEST_FIELD_NAME_USERNAME string = "username"
const REQUEST_FIELD_NAME_PASSWORD string = "password"
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

// Checks if user is logged in before sending out any data, responds with
// "Not logged in" error if user not logged in
func (this *Controller) Prepare() {
	sessUser := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if sessUser == nil {
		this.CustomAbort(401, "Not logged in")
	}
}

// Return true if user is admin, if no args are passed, uses session user ID,
// if single user ID is passed, checks the passed one. Fails otherwise.
func (this *Controller) IsAdmin(userIds ...int64) bool {
	var userId int64
	var err error
	if len(userIds) == 0 {
		userId = this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	} else if len(userIds) == 1 {
		userId = userIds[0]
	} else {
		beego.Error("Expecting single or no value as input")
		return false
	}
	var roles *models.UserRoles
	roles, err = models.GetUserRoles(userId)
	if err != nil {
		return false
	}
	if roles.Admin {
		return true
	}
	return false
}

// Return true if user is staff, if no args are passed, uses session user ID,
// if single user ID is passed, checks the passed one. Fails otherwise.
func (this *Controller) IsStaff(userIds ...int64) bool {
	var userId int64
	var err error
	if len(userIds) == 0 {
		userId = this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	} else if len(userIds) == 1 {
		userId = userIds[0]
	} else {
		beego.Error("Expecting single or no value as input")
		return false
	}
	var roles *models.UserRoles
	roles, err = models.GetUserRoles(userId)
	if err != nil {
		return false
	}
	if roles.Staff {
		return true
	}
	return false
}
