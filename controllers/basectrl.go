// Fabsmith controllers package, handles all the API calls
package controllers

import (
	"github.com/astaxie/beego"
)

// Field names for session variables
const SESSION_USER_ID string = "user_id"
const SESSION_USERNAME string = "username"

// Field names for request variables
const REQUEST_USER_ID string = "user_id"
const REQUEST_USERNAME string = "username"
const REQUEST_PASSWORD string = "password"
const REQUEST_MACHINE_ID string = "machine_id"
const REQUEST_ACTIVATION_ID string = "activation_id"

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
	sessUser := this.GetSession(SESSION_USER_ID)
	if sessUser == nil {
		this.CustomAbort(401, "Not logged in")
	}
}
