// controllers package handles all API calls (/api)
package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/memcache"
	"github.com/bradfitz/gomemcache/memcache"
	gsm "github.com/bradleypeabody/gorilla-sessions-memcache"
	"strings"
)

// Field names for session variables
const (
	SESSION_USER_ID         = "user_id"
	SESSION_USERNAME        = "username"
	SESSION_LOCATION_ID     = "location_id"
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

const SESSION_NAME = "easylab"

var (
	runmodeTest bool
	store       *gsm.MemcacheStore
)

func init() {
	runmode := beego.AppConfig.String("runmode")
	if runmode == "dev" || runmode == "prod" {
		dsn := beego.AppConfig.String("SessionProviderConfig")
		memcacheClient := memcache.New(dsn)
		secret := []byte(beego.AppConfig.String("sessionhashkey"))
		store = gsm.NewMemcacheStore(memcacheClient, "fabsmith_", secret)
	} else {
		runmodeTest = true
	}
}

// Creates new ErrorResponse instance with Status:"error" already set
func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{Status: "error"}
}

func (this *Controller) GetSession(key string) interface{} {
	if runmodeTest {
		return this.Controller.GetSession(key)
	} else {
		session, err := store.Get(this.Ctx.Request, SESSION_NAME)
		if err != nil {
			beego.Error("GetSession:", err)
		}
		return session.Values[key]
	}
}

func (this *Controller) SetSession(key string, value interface{}) {
	if runmodeTest {
		this.Controller.SetSession(key, value)
	} else {
		session, err := store.Get(this.Ctx.Request, SESSION_NAME)
		if err != nil {
			beego.Error("GetSession:", err)
		}
		session.Values[key] = value
		err = session.Save(this.Ctx.Request, this.Ctx.ResponseWriter)
		if err != nil {
			beego.Error("Error saving session:", err)
		}
	}
}

// Checks if user is logged in before sending out any data, responds with
// "Not logged in" error if user not logged in
func (this *Controller) Prepare() {
	path := this.Ctx.Request.URL.Path
	if !strings.HasPrefix(path, "/api") {
		return
	}
	switch path {
	case "/api/users/current", "/api/users/forgot_password", "/api/users/forgot_password/phone", "/api/users/forgot_password/reset", "/api/machine_types", "/api/machines/search", "/api/locations", "/api/locations/my_ip", "/api/metrics/realtime", "/api/reservations/icalendar", "/api/settings/terms_url":
	default:
		sessUser := this.GetSession(SESSION_USER_ID)
		if sessUser == nil {
			this.CustomAbort(401, "Not logged in")
		}
	}
}
