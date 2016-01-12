// Fabsmith controllers package, handles all the API calls
package controllers

import (
	"github.com/astaxie/beego"
)

// Field names for session variables
const (
	SESSION_USER_ID         = "user_id"
	SESSION_USERNAME        = "username"
	SESSION_BROWSER         = "browser"
	SESSION_IP              = "ip"
	SESSION_ACCEPT_LANGUAGE = "accept_language"
	SESSION_ACCEPT_ENCODING = "accept_encoding"
)

// Field names for request variables
const (
	REQUEST_USER_ID       = "user_id"
	REQUEST_USERNAME      = "username"
	REQUEST_PASSWORD      = "password"
	REQUEST_MACHINE_ID    = "machine_id"
	REQUEST_ACTIVATION_ID = "activation_id"
)

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
