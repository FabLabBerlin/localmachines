package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
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

// @Title Post
// @Description create user and associated tables
// @Param	email		query 	string	true		"The new user's E-Mail"
// @Success 201 {object} models.User
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router / [post]
func (this *UsersController) Post() {
	email := this.GetString("email")

	// TODO: validate email

	sid := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	if !this.IsAdmin(sid) && !this.IsStaff(sid) {
		beego.Error("Unauthorized attempt to delete user")
		this.CustomAbort(401, "Unauthorized")
	}

	user := models.User{Email: email}
	o := orm.NewOrm()
	if err := o.Begin(); err == nil {
		id, err := o.Insert(&user)
		if err != nil {
			beego.Error("Cannot create user: ", err)
			o.Rollback()
			this.CustomAbort(500, "Internal Server Error")
		}
		user.Id = id
		user.Email = email

		/*userRoles := models.UserRoles{
			UserId: user.Id,
			Admin:  false,
			Staff:  false,
			Member: false,
		}
		if _, err := o.Insert(&userRoles); err != nil {
			beego.Error("Cannot create user roles: ", err)
			o.Rollback()
			this.CustomAbort(500, "Internal Server Error")
		}*/
	}
	if err := o.Commit(); err == nil {
		this.Data["json"] = user
		this.ServeJson()
	} else {
		beego.Error("Error committing new user")
		this.CustomAbort(500, "Internal Server Error")
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

	sid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	if uid == sid {

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
		if !this.IsAdmin(sid.(int64)) && !this.IsStaff(sid.(int64)) {
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

// @Title Delete
// @Description delete user with uid
// @Param	uid		path 	int	true		"User ID"
// @Success 200
// @Failure	403	Variable message
// @Failure	401	Unauthorized
// @router /:uid [delete]
func (this *UsersController) Delete() {
	sid := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	if !this.IsAdmin(sid) && !this.IsStaff(sid) {
		beego.Error("Unauthorized attempt to delete user")
		this.CustomAbort(401, "Unauthorized")
	}

	uid, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get :uid")
	}

	if err := models.DeleteUserAuth(uid); err != nil {
		beego.Error("Failed to delete user auth")
		this.CustomAbort(403, "Failed to delete :uid completely, please retry")
	}
	if err := models.DeleteUser(uid); err != nil {
		beego.Error("Failed to delete user")
		this.CustomAbort(403, "Failed to delete :uid")
	}
}

type UserPutRequest struct {
	User models.User
}

// @Title Put
// @Description Update user with uid
// @Param	uid		path 	int	true		"User ID"
// @Success	200
// @Failure	400	Variable message
// @Failure	401	Unauthorized
// @Failure	403	Variable message
// @router /:uid [put]
func (this *UsersController) Put() {

	if !this.IsAdmin() {
		beego.Error("Unauthorized attempt update user")
		this.CustomAbort(401, "Unauthorized")
	}

	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := UserPutRequest{}
	if err := dec.Decode(&req); err == nil {
		beego.Info("req: ", req)
	} else {
		beego.Error("Failed to decode json")
		this.CustomAbort(400, "Failed to decode json")
	}

	// If the user is trying to disable his admin status
	// do not allow to do that
	sessUserId, ok := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	if !ok {
		beego.Error("Failed to get session user ID")
		this.CustomAbort(403, "Failed to update user")
	}

	if sessUserId == req.User.Id && req.User.UserRole != "admin" {
		beego.Error("User can't unadmin itself")
		this.CustomAbort(403, "selfAdmin")
	}

	err := models.UpdateUser(&req.User)
	if err != nil {
		beego.Error("Failed to update user:", err)
		this.CustomAbort(403, "lastAdmin")
	}
}

// @Title GetUserMachines
// @Description Get user machines
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.Machine
// @Failure	403	Failed to get user machines
// @Failure	401	Not authorized
// @router /:uid/machines [get]
func (this *UsersController) GetUserMachines() {

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

	if suid.(int64) != ruid {
		if !this.IsAdmin(suid.(int64)) && !this.IsStaff(suid.(int64)) {

			// The currently logged in user is not allowed to access
			// other user machines
			beego.Error("Not authorized")
			this.CustomAbort(401, "Not authorized")
		}
	}

	// Get the machines!
	var machines []*models.Machine
	if !this.IsAdmin(ruid) && !this.IsStaff(ruid) {

		// If the requested user roles is not admin and staff
		// we need to get machine permissions first and then the machines
		var permissions *[]models.Permission
		permissions, err = models.GetUserPermissions(ruid)
		if err != nil {
			beego.Error("Failed to get user machine permissions: ", err)
			this.CustomAbort(403, "Failed to get user machines")
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
	} else {

		// The requested user is also an admin, list all machines
		machines, err = models.GetAllMachines()
		if err != nil {
			beego.Error("Failed to get all machines")
			this.CustomAbort(403, "Failed to get all machines")
		}
	}

	// Serve machines
	this.Data["json"] = machines
	this.ServeJson()
}

// @Title PostUserMemberships
// @Description Post user membership
// @Param	uid		path 	int	true		"User ID"
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

	// Get requested user membership Id
	userMembershipId, err := this.GetInt64("UserMembershipId")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get :uid")
	}

	// Get requested start date
	startDate, err := time.Parse("2006-01-02", this.GetString("StartDate"))
	if err != nil {
		beego.Error("Failed to parse startDate")
		this.CustomAbort(400, "Failed to obtain start date")
	}

	o := orm.NewOrm()
	um := models.UserMembership{
		UserId:       ruid,
		MembershipId: userMembershipId,
		StartDate:    startDate,
	}

	if _, err := o.Insert(&um); err != nil {
		beego.Error("Error creating new user membership: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}
}

// @Title GetUserMemberships
// @Description Get user memberships
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.UserMembership
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

	if suid.(int64) != ruid {
		if !this.IsAdmin() && !this.IsStaff() {

			// The currently logged in user is not allowed to access
			// other user machines
			beego.Error("Not authorized")
			this.CustomAbort(401, "Not authorized")
		}
	}

	// If the requested user roles is not admin and staff
	// we need to get machine permissions first and then the machines
	ums, err := models.GetUserMemberships(ruid)
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

	if err := models.DeleteUserMembership(umid); err != nil {
		beego.Error("Failed to delete user membership")
		this.CustomAbort(403, "Failed to :umid")
	}
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
