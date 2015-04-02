package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
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
	sid := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)
	if !this.IsAdmin(sid) && !this.IsStaff(sid) {
		beego.Error("Unauthorized attempt to delete user")
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
	o := orm.NewOrm()
	if _, err := o.Update(&req.User); err != nil {
		beego.Error("Failed to update user")
		this.CustomAbort(400, "Failed to update user")
	}
	/*if _, err := o.Update(&req.UserRoles); err != nil {
		beego.Error("Failed to update user roles", err)
		this.CustomAbort(400, "Failed to update user roles")
	}*/
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
		var permissions []*models.Permission
		permissions, err = models.GetUserPermissions(ruid)
		if err != nil {
			beego.Error("Failed to get user machine permissions")
			this.CustomAbort(403, "Failed to get user machines")
		}
		for _, permission := range permissions {
			machine, err := models.GetMachine(permission.MachineId)
			if err != nil {
				beego.Error("Failed to get machine ID", permission.MachineId)
				this.CustomAbort(403, "Failed to get user machines")
			}
			machines = append(machines, machine)
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

// @Title GetUserMemberships
// @Description Get user memberships
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} []models.UserMembership
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
