package controllers

import (
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

	// Check if user is admin or staff
	userRoles, err := models.GetUserRoles(uid.(int64))
	if err != nil {
		beego.Error("Failed to get user roles")
		this.CustomAbort(403, "Failed to get user roles")
	}

	if !userRoles.Admin && !userRoles.Staff {
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
// @Description create user
// @Param	email		query 	string	true		"The new user's E-Mail"
// @Success 201 {object} models.User
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router / [post]
func (this *UsersController) Post() {
	email := this.GetString("email")
	var err error

	beego.Info("email:", email)

	uid := this.GetSession(SESSION_FIELD_NAME_USER_ID)
	beego.Info("uid:", uid)
	// Check if user is admin or staff
	userRoles, err := models.GetUserRoles(uid.(int64))
	if !userRoles.Admin && !userRoles.Staff {
		beego.Error("Not authorized to create user")
		this.CustomAbort(401, "Not authorized")
	}

	o := orm.NewOrm()
	user := models.User{Email: email}
	id, err := o.Insert(&user)
	if err == nil {
		user.Id = id
		user.Email = email
		this.Data["json"] = user
		this.ServeJson()
	} else {
		beego.Error("Cannot create user: ", err)
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
		userRoles, err := models.GetUserRoles(sid.(int64))
		if err != nil {
			beego.Error("Failed to get user roles")
			this.CustomAbort(403, "Failed to get user roles")
		} else {
			if !userRoles.Admin && !userRoles.Staff {
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
	}
	this.ServeJson()
}

// @Title Get
// @Description delete user with uid
// @Param	uid		path 	int	true		"User ID"
// @Success 200
// @Failure	403	Variable message
// @Failure	401	Unauthorized
// @router /:uid [delete]
func (this *UsersController) Delete() {
	if !this.IsAdmin() && !this.IsStaff() {
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
	var sessionUserRoles *models.UserRoles
	sessionUserRoles, err = models.GetUserRoles(suid.(int64))
	if err != nil {
		beego.Error("Failed to get session user roles")
		this.CustomAbort(403, "Failed to get session user roles")
	}
	if suid.(int64) != ruid {
		if !sessionUserRoles.Admin && !sessionUserRoles.Staff {

			// The currently logged in user is not allowed to access
			// other user machines
			beego.Error("Not authorized")
			this.CustomAbort(401, "Not authorized")
		}
	}

	// Get requested user roles
	var requestedUserRoles *models.UserRoles
	if suid.(int64) == ruid {
		requestedUserRoles = sessionUserRoles
	} else {
		requestedUserRoles, err = models.GetUserRoles(ruid)
		if err != nil {
			beego.Error("Failed to get requested user roles")
			this.CustomAbort(403, "Failed to get user machines")
		}
	}

	// Get the machines!
	var machines []*models.Machine
	if !requestedUserRoles.Admin && !requestedUserRoles.Staff {

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

// @Title GetUserRoles
// @Description Get user roles
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.UserRoles
// @Failure	403	Failed to get user roles
// @Failure	401	Not authorized
// @router /:uid/roles [get]
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

	userRoles, err = models.GetUserRoles(uid)
	this.Data["json"] = userRoles
	this.ServeJson()
}
