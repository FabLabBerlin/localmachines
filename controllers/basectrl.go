// controllers package handles all API calls (/api)
package controllers

import (
	"errors"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/memcache"
	"github.com/boj/redistore"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	store       *redistore.RediStore
)

func init() {
	runmode := beego.AppConfig.String("runmode")
	if runmode == "dev" || runmode == "prod" {
		var err error
		secret := []byte(beego.AppConfig.String("sessionhashkey"))
		store, err = redistore.NewRediStoreWithPool(redis.GetPool(), secret)
		if err != nil {
			panic(err.Error())
		}
		store.SetMaxAge(86400)
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

func (this *Controller) DestroySession() {
	if runmodeTest {
		this.Controller.DestroySession()
	} else {
		session, err := store.Get(this.Ctx.Request, SESSION_NAME)
		if err != nil {
			beego.Error("GetSession:", err)
		}
		delete(session.Values, SESSION_USER_ID)
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

func (this *Controller) GetSessionUserId() (int64, error) {
	tmp := this.GetSession(SESSION_USER_ID)
	isWs := this.Ctx.Input.IsWebsocket()
	if sid, ok := tmp.(int64); ok {
		ip := this.GetSession(SESSION_IP)
		if ip != this.ClientIp() {
			beego.Error("GetSessionUserId: wrong IP")
			return 0, errors.New("user not correctly logged in")
		}
		browser := this.GetSession(SESSION_BROWSER)
		if browser != this.Ctx.Input.UserAgent() {
			beego.Error("GetSessionUserId: wrong browser, ip=", ip)
			return 0, errors.New("user not correctly logged in")
		}
		accLang := this.GetSession(SESSION_ACCEPT_LANGUAGE)
		if accLang != this.Ctx.Input.Header("Accept-Language") && !isWs {
			beego.Error("GetSessionUserId: wrong Accept-Language, ip=", ip)
			beego.Info("expected", accLang)
			beego.Info("but got", this.Ctx.Input.Header("Accept-Language"))
			return 0, errors.New("user not correctly logged in")
		}
		return sid, nil
	} else {
		return 0, errors.New("User not logged in")
	}
}

func (this *Controller) GetSessionLocationId() (locId int64, ok bool) {
	locId, ok = this.GetSession(SESSION_LOCATION_ID).(int64)
	return
}

func (this *Controller) SetSessionLocationId(locId int64) {
	c := &http.Cookie{
		Name:    "location",
		Value:   strconv.FormatInt(locId, 10),
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, 1),
		Secure:  beego.AppConfig.String("runmode") == "prod",
	}
	http.SetCookie(this.Ctx.ResponseWriter, c)
}

func (this *Controller) SetLogged(username string, userId int64, locationId int64) {
	this.SetSession(SESSION_USERNAME, username)
	this.SetSession(SESSION_USER_ID, userId)
	this.SetSession(SESSION_LOCATION_ID, locationId)
	this.SetSession(SESSION_BROWSER, this.Ctx.Input.UserAgent())
	this.SetSession(SESSION_IP, this.ClientIp())
	//this.SetSession(SESSION_ACCEPT_ENCODING, this.Ctx.Input.Header("Accept-Encoding"))
	this.SetSession(SESSION_ACCEPT_LANGUAGE, this.Ctx.Input.Header("Accept-Language"))
	this.SetSessionLocationId(locationId)
}

func (this *Controller) IsLogged() bool {
	_, err := this.GetSessionUserId()
	return err == nil
}

// Return true if user is super admin, if no args are passed, uses session user ID,
// if single user ID is passed, checks the passed one. Fails otherwise.
func (this *Controller) IsSuperAdmin() bool {
	beego.Info("IsSuperAdmin()")
	role := this.globalUserRole()
	beego.Info("glboal user role:", role)
	return role == user_roles.SUPER_ADMIN
}

// Return true if user is admin at that location, if only the location id is
// passed, uses session user ID, if single user ID is passed, checks the passed
// one. Fails otherwise.
func (this *Controller) IsAdminAt(locationId int64) bool {
	role := this.localUserRole(locationId)
	return role == user_roles.ADMIN ||
		role == user_roles.SUPER_ADMIN ||
		this.IsSuperAdmin()
}

func (this *Controller) IsStaffAt(locationId int64) bool {
	role := this.localUserRole(locationId)
	return role == user_roles.STAFF ||
		role == user_roles.ADMIN ||
		role == user_roles.SUPER_ADMIN ||
		this.IsSuperAdmin()
}

func (this *Controller) IsApiAt(locationId int64) bool {
	role := this.localUserRole(locationId)
	return role == user_roles.API ||
		role == user_roles.ADMIN ||
		role == user_roles.STAFF ||
		role == user_roles.SUPER_ADMIN ||
		this.IsSuperAdmin()
}

// Return true if user is member at that location, if only the location id is
// passed, uses session user ID, if single user ID is passed, checks the passed
// one. Fails otherwise.
func (this *Controller) IsMemberAt(locationId int64) bool {
	role := this.localUserRole(locationId)
	return role == user_roles.MEMBER ||
		role == user_roles.ADMIN ||
		role == user_roles.STAFF ||
		role == user_roles.SUPER_ADMIN ||
		this.IsSuperAdmin()
}

func (this *Controller) globalUserRole() user_roles.Role {
	userId, ok := this.getUserId()
	if !ok {
		beego.Info("globalUserRole: couldn't get user id")
		return user_roles.NOT_AFFILIATED
	}
	user, err := users.GetUser(userId)
	if err != nil {
		beego.Info("globalUserRole: couldn't get user:", err)
		return user_roles.NOT_AFFILIATED
	}
	return user.GetRole()
}

func (this *Controller) localUserRole(locationId int64) user_roles.Role {
	userId, ok := this.getUserId()
	if !ok {
		return user_roles.NOT_AFFILIATED
	}
	uls, err := user_locations.GetAllForUser(userId)
	if err != nil {
		return user_roles.NOT_AFFILIATED
	}
	for _, ul := range uls {
		if ul.LocationId == locationId && ul.UserId == userId {
			return ul.GetRole()
		}
	}
	return user_roles.NOT_AFFILIATED
}

func (this *Controller) getUserId() (userId int64, ok bool) {
	userId, err := this.GetSessionUserId()
	if err != nil {
		return 0, false
	}
	return userId, true
}

// GetLocIdAdmin gets the location id, if passed as URL parameter, otherwise
// it will be 0.  0 being synonym for all locations.  Also it returns whether
// the user is allowed to perform admin tasks at that location.
func (this *Controller) GetLocIdAdmin() (locId int64, authorized bool) {
	locId, err := this.GetInt64("location")
	if err == nil {
		if !this.IsAdminAt(locId) {
			return
		}
	} else {
		locId = 0
		if !this.IsSuperAdmin() {
			return
		}
	}
	return locId, true
}

func (this *Controller) GetLocIdStaff() (locId int64, authorized bool) {
	locId, err := this.GetInt64("location")
	if err == nil {
		if !this.IsStaffAt(locId) {
			return
		}
	} else {
		locId = 0
		if !this.IsSuperAdmin() {
			return
		}
	}
	return locId, true
}

func (this *Controller) GetLocIdApi() (locId int64, authorized bool) {
	locId, err := this.GetInt64("location")
	if err == nil {
		if !this.IsApiAt(locId) {
			return
		}
	} else {
		locId = 0
		if !this.IsSuperAdmin() {
			return
		}
	}
	return locId, true
}

// GetLocIdAdmin gets the location id, if passed as URL parameter, otherwise
// it will be 0.  0 being synonym for all locations.  Also it returns whether
// the user is allowed to perform member tasks at that location.
func (this *Controller) GetLocIdMember() (locId int64, authorized bool) {
	locId, err := this.GetInt64("location")
	if err == nil {
		if !this.IsMemberAt(locId) {
			return
		}
	} else {
		locId = 0
		if !this.IsSuperAdmin() {
			return
		}
	}
	return locId, true
}

func (this *Controller) ClientIp() (ip string) {
	ip = this.Ctx.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = this.Ctx.Input.IP()
	}
	return
}
