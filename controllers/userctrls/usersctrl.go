// /api/users
package userctrls

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/lib/mailchimp"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type Controller struct {
	controllers.Controller
	i int
}

func (this *Controller) GetRouteUid() (uid int64, authorized bool) {
	// Check if logged in
	suid, err := this.GetSessionUserId()
	if err != nil {
		beego.Info("GetRouteUid: Not logged in")
		return 0, false
	}

	// Get requested user ID
	ruid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid", err)
		return 0, false
	}

	if suid != ruid {
		if !this.IsAdminForUser(ruid) {
			beego.Error("Not authorized")
			return 0, false
		}
	}

	if ruid <= 0 {
		beego.Error("uid < 0")
		return 0, false
	}
	return ruid, true
}

func (this *Controller) IsAdminForUser(uid int64) (authorized bool) {
	uls, err := user_locations.GetAllForUser(uid)
	if err != nil {
		beego.Error("IsAdminForUser: user_locations.GetAllForUser:", err)
		this.Abort("500")
	}
	for _, ul := range uls {
		if this.IsAdminAt(ul.LocationId) {
			return true
		}
	}
	return false
}

type UsersController struct {
	Controller
}

// Override our custom root controller's Prepare method as it is checking
// if we are logged in and we don't want that here at this point
func (this *UsersController) Prepare() {
	beego.Info("Skipping global login check")
}

// @Title login
// @Description Logs user into the system
// @Param	username		body 	string	true		"The username for login"
// @Param	password		body 	string	true		"The password for login"
// @Param	location		body	int		true		"Location ID"
// @Param	admin			body	bool	false		"Define whether login for Admin Panel"
// @Success 200 {object} models.LoginResponse
// @Failure 400 Bad Request
// @Failure 401 Failed to authenticate
// @router /login [post]
func (this *UsersController) Login() {
	locId, err := this.GetInt64("location")
	if err != nil {
		beego.Error("get location:", err)
		this.CustomAbort(400, "Bad Request")
	}

	username := this.GetString("username")
	pw := this.GetString("password")

	resp := models.LoginResponse{
		LocationId: locId,
	}

	if resp.UserId, err = this.GetSessionUserId(); err == nil {
		beego.Info("already logged")
		sessionLocationId, _ := this.GetSessionLocationId()
		if sessionLocationId != locId {
			beego.Error("sessionLocationId=", sessionLocationId, "locId=", locId)
			this.CustomAbort(500, "session location id not matching with requested location id")
		}
		resp.Status = "logged"
	} else {
		reqAdmin := this.GetString("admin") != ""
		userId, unregisteredAtLocation, err := login(locId, reqAdmin, username, pw)
		if err != nil {
			beego.Error("login:", err, ", location=", locId)
			this.CustomAbort(401, "Not authorized")
		}
		this.SetLogged(username, userId, locId)
		resp.UserId = userId
		if unregisteredAtLocation {
			resp.Status = "unregistered"
		} else {
			resp.Status = "ok"
		}
	}

	this.Data["json"] = resp
	this.ServeJSON()
}

func login(locId int64, reqAdmin bool, username, password string) (
	userId int64,
	unregisteredAtLocation bool,
	err error,
) {
	userId, err = users.AuthenticateUser(username, password)
	if err != nil {
		return 0, false, fmt.Errorf("Failed to authenticate: %v", err)
	} else {
		userLocations, err := user_locations.GetAllForUser(userId)
		if err != nil {
			return 0, false, fmt.Errorf("Failed to get user locations: %v", err)
		}
		var userLocation *user_locations.UserLocation
		for _, ul := range userLocations {
			if ul.LocationId == locId {
				userLocation = ul
				break
			}
		}
		if reqAdmin {
			if userLocation == nil || userLocation.GetRole() != user_roles.ADMIN {
				u, err := users.GetUser(userId)
				if err != nil {
					return 0, false, fmt.Errorf("Failed to get user: %v", err)
				}
				if !u.SuperAdmin {
					return 0, false, fmt.Errorf("User is not admin at this location")
				}
			}
		}

		return userId, userLocation == nil, nil
	}
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {object} models.StatusResponse
// @router /logout [get]
func (this *UsersController) Logout() {
	sessUsername := this.GetSession(controllers.SESSION_USERNAME)
	beego.Info("Logging out")
	this.DestroySession()
	if sessUsername == nil {
		beego.Info("User was not logged in")
	} else {
		beego.Info("Logged out user", sessUsername)
	}
	this.Data["json"] = models.StatusResponse{
		Status: "ok",
	}
	this.ServeJSON()
}

// @Title GetCurrentUser
// @Description Get current user
// @Success 200 users.User
// @Failure 400 Failed to authenticate
// @Failure 500 Internal Server Error
// @router /current [get]
func (this *UsersController) GetCurrentUser() {
	uid, err := this.GetSessionUserId()
	if err != nil {
		beego.Error("GetUser Not logged in")
	} else {
		u, err := users.GetUser(uid)
		if err != nil {
			beego.Error("models.GetUser:", err)
			this.CustomAbort(500, "Internal Server Error")
		}
		this.Data["json"] = u
	}
	this.ServeJSON()
}

// @Title GetAll
// @Description Get all users
// @Success 200 users.User
// @Failure	500	Failed to get all users
// @router / [get]
func (this *UsersController) GetAll() {

	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	users, err := users.GetAllUsersAt(locId)
	if err != nil {
		this.CustomAbort(500, "Failed to get all users")
	}
	this.Data["json"] = users
	this.ServeJSON()
}

type UserSignupRequest struct {
	User     users.User
	Password string
}

// @Title Signup
// @Description Accept user signup, create a zombie user with no privileges for later access
// @Param	model	body	UserSignupRequest	true	"User model and password"
// @Success 200 string ok
// @Failure 500 Internal Server Error
// @router /signup [post]
func (this *UsersController) Signup() {
	var err error
	var userId int64

	locId, _ := this.GetInt64("location")

	newsletter := this.GetString("newsletter") == "true"

	// Get body as array of models.User
	// Attempt to decode passed json
	jsonDecoder := json.NewDecoder(this.Ctx.Request.Body)
	data := UserSignupRequest{}
	if err = jsonDecoder.Decode(&data); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Attempt to create the user
	if userId, err = users.CreateUser(&data.User); err == users.ErrUsernameExists {
		this.CustomAbort(400, err.Error())
	} else if err == users.ErrEmailExists {
		this.CustomAbort(400, err.Error())
	} else if err != nil {
		beego.Error("Failed to create user:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	if newsletter {
		beego.Info("signup up for lab newsletter")
		go func() {
			if err := mailchimp.LocationSubscribe(locId, data.User.Email); err != nil {
				beego.Error("signup for location", locId, "newsletter:", err)
			}
		}()
	} else {
		beego.Info("no signup up for lab newsletter")
	}

	// Set the password
	if err = users.AuthSetPassword(userId, data.Password); err != nil {
		beego.Error("Failed to set password for user ID", userId)
		this.CustomAbort(500, "Internal Server Error")
	}

	if locId > 0 {
		ul := &user_locations.UserLocation{
			UserId:     userId,
			LocationId: locId,
			UserRole:   user_roles.MEMBER.String(),
		}
		if _, err := user_locations.Create(ul); err != nil {
			beego.Error("Failed to create user location for new user:", err)
		}
	}

	this.Data["json"] = userId
	this.ServeJSON()
}

// @Title Post
// @Description create user and associated tables
// @Param	email		query 	string	true	"The new user's E-Mail"
// @Param	location	query 	int64	false	"Make user member of location id"
// @Success 201 users.User
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router / [post]
func (this *UsersController) Post() {
	email := this.GetString("email")

	locId, err := this.GetInt64("location")
	if err != nil {
		beego.Error("No location specified")
		this.Abort("400")
	}

	if !this.IsAdminAt(locId) {
		beego.Error("Unauthorized attempt to delete user")
		this.Abort("401")
	}

	user := users.User{
		Email: email,
	}

	// Attempt to create the user.
	// The CreateUser function takes (or should take)
	// care of validating the email.
	if userId, err := users.CreateUser(&user); err != nil {
		beego.Error("Failed to create user:", err)
		this.CustomAbort(500, "Internal Server Error")
	} else {
		user.Id = userId
		ul := &user_locations.UserLocation{
			UserId:     userId,
			LocationId: locId,
			UserRole:   user_roles.MEMBER.String(),
		}
		if _, err := user_locations.Create(ul); err != nil {
			beego.Error("Failed to create user location for new user:", err)
		}
		this.Data["json"] = user
		this.ServeJSON()
	}
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	int	true		"User ID"
// @Success 200 users.User
// @Failure	403	Variable message
// @Failure	401	Unauthorized
// @router /:uid [get]
func (this *UsersController) Get() {
	if this.GetString(":uid") == "current" {
		this.GetCurrentUser()
		return
	}

	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	user, err := users.GetUser(uid)
	if err != nil {
		beego.Error("Failed to get user data")
		this.CustomAbort(403, "Failed to get user data")
	}

	this.Data["json"] = user
	this.ServeJSON()
}

type UserPutRequest struct {
	User users.User
}

// @Title Put
// @Description Update user with uid
// @Param	uid		path 	int			true	"User ID"
// @Param	body	body	models.User	true	"User model"
// @Success	200	ok
// @Failure	400	Variable message
// @Failure	401	Unauthorized
// @Failure	403	Variable message
// @router /:uid [put]
func (this *UsersController) Put() {

	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := UserPutRequest{}
	if err := dec.Decode(&req); err == nil {
		beego.Info("req: ", req)
	} else {
		beego.Error("Failed to decode json", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// If the user is trying update his own information
	// let him do so. Check that by comparing session user ID
	// with the one passed as :uid
	sessUserId, err := this.GetSessionUserId()
	if err != nil {
		beego.Error("Failed to get session user ID")
		this.CustomAbort(500, "Internal Server Error")
	}

	if req.User.Id != sessUserId {
		if !this.IsAdminForUser(req.User.Id) {
			beego.Error("Unauthorized attempt update user")
			this.CustomAbort(401, "Not authorized")
		}
	}

	// Do not allow change user role if not admin and self
	existingUser, err := users.GetUser(req.User.Id)
	if err != nil {
		beego.Error("User does not exist, user ID:", req.User.Id)
		this.CustomAbort(500, "Internal Server Error")
	}

	if existingUser.SuperAdmin != req.User.SuperAdmin {
		beego.Error("User is not authorized to change UserRole")
		this.CustomAbort(403, "Internal Server Error")
	}

	err = req.User.Update()
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			beego.Error("Failed to update user due to duplicate entry:", err)
			this.CustomAbort(400, "duplicateEntry")
		} else if strings.Contains(err.Error(), "same username") {
			beego.Error("Failed to update username:", err)
			this.CustomAbort(500, "Internal Server Error")
		} else if strings.Contains(err.Error(), "same email") {
			beego.Error("Failed to update user email:", err)
			this.CustomAbort(500, "Internal Server Error")
		} else {
			beego.Error("Failed to update user:", err)
			this.Abort("500")
		}
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title GetUserMachines
// @Description Get user machines, all machines for admin user
// @Param	uid		path 	int	true		"User ID"
// @Param	location	query	int	false		"Location ID"
// @Success 200 machine.Machine
// @Failure	500	Internal Server Error
// @Failure	401	Unauthorized
// @router /:uid/machines [get]
func (this *UsersController) GetUserMachines() {
	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(401, "Wrong uid in url or not authorized")
	}
	locationId, err := this.GetInt64("location")
	if err != nil {
		beego.Error("location id not specified")
		this.Abort("400")
	}

	// List all machines if the requested user is admin
	allMachines, err := machine.GetAll()
	if err != nil {
		beego.Error("Failed to get all machines: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Get the machines!
	machines := make([]*machine.Machine, 0, len(allMachines))
	if !this.IsStaffAt(locationId) {
		permissions, err := user_permissions.Get(uid)
		if err != nil {
			beego.Error("Failed to get user machine permissions: ", err)
			this.CustomAbort(500, "Internal Server Error")
		}
		for _, machine := range allMachines {
			for _, permission := range *permissions {
				machine.Locked = true
				if machine.TypeId == permission.CategoryId {
					machine.Locked = false
					break
				}
			}
			machines = append(machines, machine)
		}
	} else {
		machines = allMachines
	}

	filteredByLocation := make([]*machine.Machine, 0, len(machines))
	for _, m := range machines {
		if locationId <= 0 || locationId == m.LocationId {
			m.HideSensitiveData()
			filteredByLocation = append(filteredByLocation, m)
		}
	}

	// Serve machines
	this.Data["json"] = filteredByLocation
	this.ServeJSON()
}

// @Title GetUserBill
// @Description Get a user PayAsYouGo data (Machines, usage and price per machine and total price)
// @Param	uid		path 	int	true		"User ID"
// @Success 200 machine.Machine
// @Failure	401	Unauthorized
// @Failure	500	Internal Server Error
// @router /:uid/bill [get]
func (this *UsersController) GetUserBill() {
	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(401, "Wrong uid in url or not authorized")
	}
	locId, authorized := this.GetLocIdMember()
	if !authorized {
		this.CustomAbort(401, "Not authorized for this location")
	}
	invs, err := invutil.GetAllOfUserAt(locId, uid)
	if err != nil {
		beego.Error("invutil.GetAllOfUserAt:", err)
		this.Abort("500")
	}

	for _, inv := range invs {
		for _, p := range inv.Purchases {
			if m := p.Machine; m != nil {
				m.HideSensitiveData()
			}
		}
	}

	this.Data["json"] = invs
	this.ServeJSON()
}

type TotalUsage struct {
	TotalTime  int
	TotalPrice float64
	Details    []MachineUsage
}

type MachineUsage struct {
	MachineId   int64
	MachineName string
	Price       float64
	Time        int
}

type MachinesAffectedArray struct {
	MachinesIds []int
}

// @Title GetUserNames
// @Description Get user name data only
// @Param	uids		query 	string	true		"User IDs"
// @Success 200 {object} models.UserNamesResponse
// @Failure	400	Bad request
// @Failure	401	Not logged in
// @Failure	500	Internal Server Error
// @router /names [get]
func (this *UsersController) GetUserNames() {

	// Check if logged in
	if !this.IsLogged() {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Get the user names data
	var uids []int64

	if uidsString := this.GetString("uids"); len(uidsString) > 0 {
		uidsList := strings.Split(uidsString, ",")
		uids = make([]int64, 0, len(uidsList))
		for _, uid := range uidsList {
			id, err := strconv.ParseInt(uid, 10, 64)
			if err != nil {
				beego.Error("Failed to parse uids:", err)
				this.CustomAbort(400, "Failed to get user name")
			}
			uids = append(uids, id)
		}
	}

	response := models.UserNamesResponse{
		Users: make([]models.UserNameResponse, 0, len(uids)),
	}

	for _, uid := range uids {
		user, err := users.GetUser(uid)
		if err != nil {
			beego.Error("Failed not get user name:", err)
			this.CustomAbort(500, "Failed not get user name")
		}
		data := models.UserNameResponse{
			UserId:    user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
		response.Users = append(response.Users, data)
	}

	this.Data["json"] = response
	this.ServeJSON()
}

// @Title PostUserPassword
// @Description Post user password
// @Param	uid		path 	int	true		"User ID"
// @Success 200
// @Failure	403	Failed to get user
// @Failure	401	Not authorized
// @router /:uid/password [post]
func (this *UsersController) PostUserPassword() {
	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(400, "Wrong uid in url or not authorized")
	}

	err := users.AuthSetPassword(uid, this.GetString("password"))
	if err != nil {
		beego.Error("Unable to update password: ", err)
		this.CustomAbort(403, "Unable to update password")
	}

	this.Data["json"] = models.StatusResponse{
		Status: "Password changed successfully!",
	}
	this.ServeJSON()
}
