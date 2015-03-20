package controllers

import (
	"github.com/kr15h/fabsmith/models"
	"github.com/astaxie/beego"
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
	var userId int
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
		userId = this.GetSession(SESSION_FIELD_NAME_USER_ID).(int)
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
	userRoles, err := models.GetUserRoles(uid.(int))
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
	uid, err := this.GetInt(":uid")
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
		userRoles, err := models.GetUserRoles(uid)
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
	var ruid int
	ruid, err = this.GetInt(":uid")
	if err != nil {
		beego.Error("Failed to get :uid")
		this.CustomAbort(403, "Failed to get :uid")
	}

	// We need the user roles in order to understand
	// whether we are allowed to access other user machines
	var sessionUserRoles *models.UserRoles
	sessionUserRoles, err = models.GetUserRoles(suid.(int))
	if err != nil {
		beego.Error("Failed to get session user roles")
		this.CustomAbort(403, "Failed to get session user roles")
	}
	if (suid.(int) != ruid) {
		if !sessionUserRoles.Admin && !sessionUserRoles.Staff {
			
			// The currently logged in user is not allowed to access
			// other user machines
			beego.Error("Not authorized")
			this.CustomAbort(401, "Not authorized")
		}
	} 
	
	// Get requested user roles
	var requestedUserRoles *models.UserRoles
	if suid.(int) == ruid {
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
