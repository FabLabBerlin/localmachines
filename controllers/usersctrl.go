package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
	"strings"
	"time"
)

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
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {object} models.LoginResponse
// @Failure 401 Failed to authenticate
// @router /login [post]
func (this *UsersController) Login() {
	var userId int64
	var err error
	sessUsername := this.GetSession(SESSION_FIELD_NAME_USERNAME)
	if sessUsername == nil {
		username := this.GetString("username")
		password := this.GetString("password")
		userId, err = models.AuthenticateUser(username, password)
		if err != nil {
			this.CustomAbort(401, "Failed to authenticate")
		} else {
			this.SetSession(SESSION_FIELD_NAME_USERNAME, username)
			this.SetSession(SESSION_FIELD_NAME_USER_ID, userId)
			this.Data["json"] = models.LoginResponse{"ok", userId}
		}
	} else {
		userId = this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
		this.Data["json"] = models.LoginResponse{"logged", userId}
	}
	this.ServeJson()
}

// @Title LoginUid
// @Description Logs user into the system by using NFC UID
// @Param	uid		query 	string	true		"The NFC UID"
// @Success 200 {object} models.LoginResponse
// @Failure 401 Failed to authenticate
// @router /loginuid [post]
func (this *UsersController) LoginUid() {
	var username string
	var userId int64
	var err error

	sessUsername := this.GetSession(SESSION_FIELD_NAME_USERNAME)

	if sessUsername == nil {
		uid := this.GetString("uid")
		beego.Trace("uid", uid)
		username, userId, err = models.AuthenticateUserUid(uid)
		if err != nil {
			beego.Error(err)
			this.CustomAbort(401, "Failed to authenticate")
		} else {
			this.SetSession(SESSION_FIELD_NAME_USERNAME, username)
			this.SetSession(SESSION_FIELD_NAME_USER_ID, userId)
			this.Data["json"] = models.LoginResponse{"ok", userId}
		}
	} else {
		var ok bool
		userId, ok = this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
		if !ok {
			beego.Error("Could not get session user ID")
			this.CustomAbort(401, "Failed to authenticate")
		}
		this.Data["json"] = models.LoginResponse{"logged", userId}
	}

	this.ServeJson()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {object} models.StatusResponse
// @router /logout [get]
func (this *UsersController) Logout() {
	sessUsername := this.GetSession(SESSION_FIELD_NAME_USERNAME)
	beego.Info("Logging out")
	this.DestroySession()
	if sessUsername == nil {
		beego.Info("User was not logged in")
	} else {
		beego.Info("Logged out user", sessUsername)
	}
	this.Data["json"] = models.StatusResponse{"ok"}
	this.ServeJson()
}

// @Title GetAll
// @Description Get all users
// @Success 200 {object} models.User
// @Failure	403	Failed to get all users
// @router / [get]
func (this *UsersController) GetAll() {

	// Check if logged in
	uid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if uid == nil {
		beego.Info("Attempt to get all users while not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	if !this.IsAdmin(uid.(int64)) && !this.IsStaff(uid.(int64)) {
		beego.Error("Not authorized to get all users")
		this.CustomAbort(401, "Not authorized")
	}

	users, err := models.GetAllUsers()
	if err != nil {
		this.CustomAbort(403, "Failed to get all users")
	}
	this.Data["json"] = users
	this.ServeJson()
}

type UserSignupRequest struct {
	User     models.User
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

	// Get body as array of models.User
	// Attempt to decode passed json
	jsonDecoder := json.NewDecoder(this.Ctx.Request.Body)
	data := UserSignupRequest{}
	if err = jsonDecoder.Decode(&data); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Attempt to create the user
	if userId, err = models.CreateUser(&data.User); err != nil {
		beego.Error("Failed to create user:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Set the password
	if err = models.AuthSetPassword(userId, data.Password); err != nil {
		beego.Error("Failed to set password for user ID", userId)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = userId
	this.ServeJson()
}

// @Title Post
// @Description create user and associated tables
// @Param	email		query 	string	true		"The new user's E-Mail"
// @Success 201 {object} models.User
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router / [post]
func (this *UsersController) Post() {
	email := this.GetString("email")

	if !this.IsAdmin() {
		beego.Error("Unauthorized attempt to delete user")
		this.CustomAbort(401, "Unauthorized")
	}

	user := models.User{Email: email}

	// Attempt to create the user.
	// The CreateUser function takes (or should take)
	// care of validating the email.
	if userId, err := models.CreateUser(&user); err != nil {
		beego.Error("Failed to create user:", err)
		this.CustomAbort(500, "Internal Server Error")
	} else {
		user.Id = userId
		this.Data["json"] = user
		this.ServeJson()
	}
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.User
// @Failure	403	Variable message
// @Failure	401	Unauthorized
// @router /:uid [get]
func (this *UsersController) Get() {
	var err error
	var user *models.User
	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get :uid")
	}

	suid, ok := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	if !ok {
		beego.Error("Can't get data if not logged in")
		this.CustomAbort(401, "Unauthorized")
	} else if uid == suid {

		// Request user ID and session user ID match.
		// The user is logged in and deserves to get her data.
		user, err = models.GetUser(uid)
		if err != nil {
			beego.Error("Failed to get user data")
			this.CustomAbort(403, "Failed to get user data")
		} else {
			this.Data["json"] = user
		}
	} else {
		// Requested user ID and stored session ID does not match,
		// meaning that the logged user is trying to access other user data.
		// Don't allow to get data of another user unless session user is admin or staff.
		if !this.IsAdmin(suid) && !this.IsStaff(suid) {
			beego.Error("Unauthorized attempt to get other user data")
			this.CustomAbort(401, "Unauthorized")
		} else {
			user, err = models.GetUser(uid)
			if err != nil {
				beego.Error("Failed to get other user data")
				this.CustomAbort(403, "Failed to get other user data")
			} else {
				this.Data["json"] = user
			}
		}
	}
	this.ServeJson()
}

type UserPutRequest struct {
	User models.User
}

// @Title Put
// @Description Update user with uid
// @Param	uid		path 	int	true		"User ID"
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
	sessUserId, ok := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	if !ok {
		beego.Error("Failed to get session user ID")
		this.CustomAbort(500, "Internal Server Error")
	}

	if req.User.Id != sessUserId {
		if !this.IsAdmin() {
			beego.Error("Unauthorized attempt update user")
			this.CustomAbort(401, "Not authorized")
		}
	}

	// Do not allow change user role if not admin and self
	existingUser, err := models.GetUser(req.User.Id)
	if err != nil {
		beego.Error("User does not exist, user ID:", req.User.Id)
		this.CustomAbort(500, "Internal Server Error")
	}

	if existingUser.UserRole != req.User.UserRole {
		if sessUserId == req.User.Id {
			beego.Error("User can't change his own user role")
			this.CustomAbort(500, "Internal Server Error")
		} else if !this.IsAdmin() {
			beego.Error("User is not authorized to change UserRole")
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	err = models.UpdateUser(&req.User)
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
			this.CustomAbort(403, "lastAdmin")
		}
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title GetUserMachinePermissions
// @Description Get current saved machine permissions
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.Machine
// @Failure	500	Internal Server Error
// @Failure	401	Unauthorized
// @router /:uid/machinepermissions [get]
func (this *UsersController) GetUserMachinePermissions() {

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Unauthorized")
	}

	// Get requested user ID
	var err error
	var ruid int64
	ruid, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(500, "Internal Server Error")
	}

	// Attempt to get user sessio id as int64
	suidInt64, ok := suid.(int64)
	if !ok {
		beego.Error("Failed to get session user ID")
		this.CustomAbort(500, "Internal Server Error")
	}

	// If the session user ID and the request user ID does not match
	// check whether the session user is an admin, return if not
	if suidInt64 != ruid {
		if !this.IsAdmin() {
			beego.Error("Not authorized to access other user machine permissions")
			this.CustomAbort(401, "Unauthorized")
		}
	}

	// We need to get machine permissions first and then the machines
	var permissions *[]models.Permission
	var machines []*models.Machine
	permissions, err = models.GetUserPermissions(ruid)
	if err != nil {
		beego.Error("Failed to get user machine permissions: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	for _, permission := range *permissions {
		machine, err := models.GetMachine(permission.MachineId)
		if err != nil {
			beego.Warning("Failed to get machine ID", permission.MachineId)
			// Just don't add the machine permission if not exists in db
			//this.CustomAbort(403, "Failed to get user machines")
		} else {
			machines = append(machines, machine)
		}
	}

	// Serve machines
	this.Data["json"] = machines
	this.ServeJson()
}

// @Title GetUserMachines
// @Description Get user machines, all machines for admin user
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.Machine
// @Failure	500	Internal Server Error
// @Failure	401	Unauthorized
// @router /:uid/machines [get]
func (this *UsersController) GetUserMachines() {

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Unauthorized")
	}

	// Get requested user ID
	var err error
	var ruid int64
	ruid, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	suidInt64, ok := suid.(int64)
	if !ok {
		beego.Error("Could not get session user ID as int64")
		this.CustomAbort(500, "Internal Server Error")
	}

	if suidInt64 != ruid {
		if !this.IsAdmin() {
			beego.Error("Not authorized")
			this.CustomAbort(401, "Unauthorized")
		}
	}

	// Get the machines!
	var machines []*models.Machine
	if !this.IsAdmin(ruid) {

		// If the requested user roles is not admin
		// we need to get machine permissions first and then the machines
		var permissions *[]models.Permission
		permissions, err = models.GetUserPermissions(ruid)
		if err != nil {
			beego.Error("Failed to get user machine permissions: ", err)
			this.CustomAbort(500, "Internal Server Error")
		}
		for _, permission := range *permissions {
			machine, err := models.GetMachine(permission.MachineId)
			if err != nil {
				beego.Warning("Failed to get machine ID", permission.MachineId)
			} else {
				machines = append(machines, machine)
			}
		}
	} else {

		// List all machines if the requested user is admin
		machines, err = models.GetAllMachines()
		if err != nil {
			beego.Error("Failed to get all machines: ", err)
			this.CustomAbort(500, "Internal Server Error")
		}
	}

	// Serve machines
	this.Data["json"] = machines
	this.ServeJson()
}

// @Title GetUserBill
// @Description Get a user PayAsYouGo data (Machines, usage and price per machine and total price)
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.Machine
// @Failure	401	Unauthorized
// @Failure	500	Internal Server Error
// @router /:uid/bill [get]
func (this *UsersController) GetUserBill() {
	// Get requested user ID
	var err error
	var ruid int64
	ruid, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Unauthorized")
	}

	suidInt64, ok := suid.(int64)
	if !ok {
		beego.Error("Could not get session user ID as int64")
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() && suidInt64 != ruid {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Unauthorized")
	}

	startTime, err := models.GetUserActivationsStartTime(suidInt64)
	if err != nil {
		beego.Error("GetUserActivationsStartTime:", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	startTime = startTime.Add(-86400 * time.Second)

	endTime := time.Now().Add(86400 * 30 * time.Second)
	invoice, err := models.CalculateInvoiceSummary(startTime, endTime)
	if err != nil {
		beego.Error("CalculateInvoiceSummary:", err)
	}

	var userSummary *models.UserSummary

	for _, us := range invoice.UserSummaries {
		if us.User.Id == suid {
			userSummary = us
		}
	}

	this.Data["json"] = userSummary
	this.ServeJson()
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

// @Title PostUserMemberships
// @Description Post user membership
// @Param	uid							path 		int			true		"User ID"
// @Param	membershipId		query 	int			true		"Membership ID"
// @Param	startDate				query 	string	true		"Membership ID"
// @Success 200 {object} models.UserMembership
// @Failure	403	Failed to get user memberships
// @Failure	401	Not authorized
// @router /:uid/memberships [post]
func (this *UsersController) PostUserMemberships() {

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	// Get requested user ID
	ruid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get :uid")
	}
	if ruid <= 0 {
		beego.Error("Bad :uid")
		this.CustomAbort(500, "Bad uid")
	}

	// Get requested user membership ID
	membershipId, err := this.GetInt64("membershipId")
	if err != nil {
		beego.Error("Failed to get membership ID")
		this.CustomAbort(403, "Failed to get membership ID")
	}

	// Get requested start date
	startDate, err := time.ParseInLocation("2006-01-02",
		this.GetString("startDate"),
		time.UTC)
	if err != nil {
		beego.Error("Failed to parse startDate=%v", startDate)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Create user membership by using the model function
	var userMembershipId int64
	userMembershipId, err = models.CreateUserMembership(ruid, membershipId, startDate)
	if err != nil {
		beego.Error("Error creating user membership:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Read the user membership back
	var userMembership *models.UserMembership
	userMembership, err = models.GetUserMembership(userMembershipId)
	if err != nil {
		beego.Error("Failed to get user membership:", err)
		this.CustomAbort(500, "Failed to get user membership")
	}

	this.Data["json"] = userMembership
	this.ServeJson()
}

// @Title GetUserMemberships
// @Description Get user memberships
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.UserMembershipList
// @Failure	403	Failed to get user memberships
// @Failure	401	Not authorized
// @router /:uid/memberships [get]
func (this *UsersController) GetUserMemberships() {

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Get requested user ID
	var err error
	var ruid int64
	ruid, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get :uid")
	}

	// We need the user roles in order to understand
	// whether we are allowed to access other user machines

	suidInt64, ok := suid.(int64)
	if !ok {
		beego.Error("Failed to get int64 value out of session ID")
		this.CustomAbort(500, "Internal Server Error")
	}

	if suidInt64 != ruid {
		if !this.IsAdmin() {

			// The currently logged in user is not allowed to access
			// other user machines
			beego.Error("Not authorized")
			this.CustomAbort(401, "Not authorized")
		}
	}

	// If the requested user roles is not admin and staff
	// we need to get machine permissions first and then the machines
	var ums *models.UserMembershipList
	ums, err = models.GetUserMemberships(ruid)
	if err != nil {
		beego.Error("Failed to get user machine permissions")
		this.CustomAbort(403, "Failed to get user machines")
	}

	// Serve machines
	this.Data["json"] = ums
	this.ServeJson()
}

// @Title DeleteUserMembership
// @Description Delete user membership
// @Param	uid		path 	int	true		"User ID"
// @Param	umid	path	int	true		"User Membership ID"
// @Success 200
// @Failure	403	Failed to get user memberships
// @Failure	401	Not authorized
// @router /:uid/memberships/:umid [delete]
func (this *UsersController) DeleteUserMembership() {
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	umid, err := this.GetInt64(":umid")
	if err != nil {
		beego.Error("Failed to get :umid")
		this.CustomAbort(403, "Failed to get :umid")
	}
	beego.Trace("User membership ID:", umid)

	err = models.DeleteUserMembership(umid)
	if err != nil {
		beego.Error("Failed to delete user membership")
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title Put
// @Description Update UserMembership
// @Param	uid		path 	int	true						"User Membership Id"
// @Param	body	body	models.UserMembership	true	"User Membership model"
// @Success	200	ok
// @Failure	400	Variable message
// @Failure	401	Unauthorized
// @Failure	403	Variable message
// @router /:uid/memberships/:umid [put]
func (this *UsersController) PutUserMembership() {
	dec := json.NewDecoder(this.Ctx.Request.Body)
	var userMembership models.UserMembership
	if err := dec.Decode(&userMembership); err == nil {
		beego.Info("userMembership: ", userMembership)
	} else {
		beego.Error("Failed to decode json", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	if err := models.UpdateUserMembership(&userMembership); err != nil {
		beego.Error("UpdateMembership: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title GetUserName
// @Description Get user name data only
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.UserNameResponse
// @Failure	403	Failed to get user name
// @Failure	401	Not loggen
// @router /:uid/name [get]
func (this *UsersController) GetUserName() {

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Get the user name data
	var err error
	var uid int64
	uid, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get user name")
	}
	var user *models.User
	user, err = models.GetUser(uid)
	if err != nil {
		beego.Error("Failed not get user name")
		this.CustomAbort(403, "Failed not get user name")
	}
	response := models.UserNameResponse{}
	response.UserId = user.Id
	response.FirstName = user.FirstName
	response.LastName = user.LastName
	this.Data["json"] = response
	this.ServeJson()
}

// @Title PostUserPassword
// @Description Post user password
// @Param	uid		path 	int	true		"User ID"
// @Success 200
// @Failure	403	Failed to get user
// @Failure	401	Not authorized
// @router /:uid/password [post]
func (this *UsersController) PostUserPassword() {
	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Get requested user ID
	ruid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get :uid")
	}

	if !this.IsAdmin() && suid != ruid {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	err = models.AuthSetPassword(ruid, this.GetString("password"))
	if err != nil {
		beego.Error("Unable to update password: ", err)
		this.CustomAbort(403, "Unable to update password")
	}

	this.Data["json"] = models.StatusResponse{"Password changed successfully!"}
	this.ServeJson()
}

// @Title UpdateNfcUid
// @Description Update user NFC UID
// @Param	uid		path 	int	true		"User ID"
// @Param	nfcuid		query 	int	true		"NFC UID"
// @Success 200	string	ok
// @Failure	403	Failed to update NFC UID
// @Failure	401	Not authorized
// @router /:uid/nfcuid [put]
func (this *UsersController) UpdateNfcUid() {
	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not authorized")
	}

	// Get requested user ID
	ruid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :ruid")
		this.CustomAbort(403, "Failed to update NFC UID")
	}

	// User can't change NFC UID of others if she is not an admin
	if !this.IsAdmin() && suid != ruid {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	// Get the NFC UID
	nfcuid := this.GetString("nfcuid")
	if nfcuid == "" {
		beego.Error("Empty nfcuid")
		this.CustomAbort(403, "Failed to update NFC UID")
	}

	err = models.AuthUpdateNfcUid(ruid, nfcuid)
	if err != nil {
		beego.Error("Unable to update NFC UID: ", err)
		this.CustomAbort(403, "Failed to update NFC UID")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title GetUserRoles
// @Description Get user roles
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.UserRoles
// @Failure	403	Failed to get user roles
// @Failure	401	Not authorized
// @router /:uid/roles [get]
/*
func (this *UsersController) GetUserRoles() {

	var sessionRoles *models.UserRoles
	var userRoles *models.UserRoles
	var err error

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Don't give the roles to someone not admin
	sessionRoles, err = models.GetUserRoles(suid.(int64))
	if err != nil {
		beego.Error("Failed to get session user roles")
		this.CustomAbort(403, "Failed tp get user roles")
	}

	var uid int64
	uid, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get user roles")
	}

	if !sessionRoles.Admin && !sessionRoles.Staff {
		if uid != suid.(int64) {
			beego.Error("Unauthorized attempt to get user roles")
			this.CustomAbort(401, "Not authorized")
		}
	}

	if userRoles, err = models.GetUserRoles(uid); err == nil {
		this.Data["json"] = userRoles
		this.ServeJson()
	} else {
		beego.Error("Unable to retrieve user roles")
		this.CustomAbort(500, "Internal Server Error")
	}
}
*/

// @Title CreateUserPermission
// @Description Create a permission for a user to allow him/her to use a machine
// @Param	uid		path 	int	true		"User ID"
// @Param	mid		query 	int	true		"Machine ID"
// @Success 200 string ok
// @Failure	403	Failed to create permission
// @Failure	401	Not authorized
// @router /:uid/permissions [post]
func (this *UsersController) CreateUserPermission() {

	// TODO: think about bulk permission creation or another way
	// 		 that does not use a separate table maybe.

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Only admin
	if !this.IsAdmin() {
		this.CustomAbort(401, "Not authorized")
	}

	// Get user ID for the machine permission to be made
	var err error
	var userId int64
	userId, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID")
		this.CustomAbort(403, "Failed to create permission")
	}

	// Get machine ID for the permission being made
	var machineId int64
	machineId, err = this.GetInt64("mid")
	if err != nil {
		beego.Error("Failed to get queried machine ID")
		this.CustomAbort(403, "Failed to create permission")
	}

	// Create a new user permission record in the database
	err = models.CreateUserPermission(userId, machineId)
	if err != nil {
		beego.Error("Failed to create machine user permission", err)
		this.CustomAbort(403, "Failed to create user machine permission")
	}

	// We are done!
	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title DeleteUserPermission
// @Description Delete user machine permission
// @Param	uid		path 	int	true		"User ID"
// @Param	mid		query 	int	true		"Machine ID"
// @Success 200 string ok
// @Failure	403	Failed to delete permission
// @Failure	401	Not authorized
// @router /:uid/permissions [delete]
func (this *UsersController) DeleteUserPermission() {

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Only admin
	if !this.IsAdmin() {
		this.CustomAbort(401, "Not authorized")
	}

	// Get user ID for the machine permission to be made
	var err error
	var userId int64
	userId, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID:", err)
		this.CustomAbort(403, "Failed to create permission")
	}

	// Get machine ID for the permission being made
	var machineId int64
	machineId, err = this.GetInt64("mid")
	if err != nil {
		beego.Error("Failed to get queried machine ID")
		this.CustomAbort(403, "Failed to create permission")
	}

	err = models.DeleteUserPermission(userId, machineId)
	if err != nil {
		beego.Error("Failed to delete permission:", err)
		this.CustomAbort(403, "Failed to delete permission")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title Update User Machine Permissions
// @Description Update user machine permissions
// @Param	uid		path 	int	true	"User ID"
// @Param	model	body	models.Permission	true	"Permissions Array"
// @Success 200	ok
// @Failure	403	Failed to update permissions
// @Failure	401	Not authorized
// @router /:uid/permissions [put]
func (this *UsersController) UpdateUserPermissions() {

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Only admin can do this
	if !this.IsAdmin() {
		this.CustomAbort(401, "Not authorized")
	}

	// Get request user ID
	var err error
	var userId int64
	userId, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID:", err)
		this.CustomAbort(403, "Failed to update permissions")
	}

	// Get body as array of models.Permission
	// Attempt to decode passed json
	jsonDecoder := json.NewDecoder(this.Ctx.Request.Body)
	permissions := []models.Permission{}
	if err = jsonDecoder.Decode(&permissions); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update permissions")
	}

	// Make sure that the user IDs of all the permissions are the same
	// as the request user ID
	for i := 0; i < len(permissions); i++ {
		permissions[i].UserId = userId
	}

	// Update permissions
	err = models.UpdateUserPermissions(userId, &permissions)
	if err != nil {
		beego.Error("Failed to update permissions:", err)
		this.CustomAbort(403, "Failed to update permissions")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @Title Get User Machine Permissions
// @Description Get user machine permissions
// @Param	uid		path 	int	true	"User ID"
// @Success 200	models.Permission
// @Failure	403	Failed to update permissions
// @Failure	401	Not authorized
// @router /:uid/permissions [get]
func (this *UsersController) GetUserPermissions() {

	// Check if logged in
	suid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Get session user ID as int64
	var suidInt64 int64
	var found bool
	suidInt64, found = suid.(int64)
	if !found {
		beego.Error("Could not cast session ID to int64")
		this.CustomAbort(403, "Failed to get permissions")
	}

	// Get request user ID
	var err error
	var userId int64
	userId, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID:", err)
		this.CustomAbort(403, "Failed to get permissions")
	}

	// Allow to get other user permissions only if admin
	if userId != suidInt64 {
		if !this.IsAdmin() {
			beego.Error("Not authorized to get other user permissions")
			this.CustomAbort(401, "Not authorized")
		}
	}

	var permissions *[]models.Permission
	permissions, err = models.GetUserPermissions(userId)
	if err != nil {
		beego.Error("Failed to get user permissions")
		this.CustomAbort(403, "Failed to get permissions")
	}

	this.Data["json"] = permissions
	this.ServeJson()
}
