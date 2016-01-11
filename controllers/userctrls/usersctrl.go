package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
)

type UsersController struct {
	controllers.Controller
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
	if sessUserId, err := this.GetSessionUserId(); err != nil {
		username := this.GetString("username")
		password := this.GetString("password")
		userId, err := models.AuthenticateUser(username, password)
		if err != nil {
			this.CustomAbort(401, "Failed to authenticate")
		} else {
			this.SetLogged(username, userId)
			this.Data["json"] = models.LoginResponse{"ok", userId}
		}
	} else {
		this.Data["json"] = models.LoginResponse{"logged", sessUserId}
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
	if sessUserId, err := this.GetSessionUserId(); err == nil {
		uid := this.GetString("uid")
		username, userId, err := models.AuthenticateUserUid(uid)
		if err != nil {
			beego.Error(err)
			this.CustomAbort(401, "Failed to authenticate")
		} else {
			this.SetLogged(username, userId)
			this.Data["json"] = models.LoginResponse{"ok", userId}
		}
	} else {
		this.Data["json"] = models.LoginResponse{"logged", sessUserId}
	}

	this.ServeJson()
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
	this.Data["json"] = models.StatusResponse{"ok"}
	this.ServeJson()
}

// @Title GetAll
// @Description Get all users
// @Success 200 {object} models.User
// @Failure	403	Failed to get all users
// @router / [get]
func (this *UsersController) GetAll() {

	if !this.IsAdmin() && !this.IsStaff() {
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

	suid, err := this.GetSessionUserId()
	if err != nil {
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
	sessUserId, err := this.GetSessionUserId()
	if err != nil {
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
			this.CustomAbort(403, "lastAdmin")
		}
	}

	this.Data["json"] = "ok"
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
	suid, err := this.GetSessionUserId()
	if err != nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Unauthorized")
	}

	// Get requested user ID
	var ruid int64
	ruid, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	if suid != ruid {
		if !this.IsAdmin() {
			beego.Error("Not authorized")
			this.CustomAbort(401, "Unauthorized")
		}
	}

	// List all machines if the requested user is admin
	allMachines, err := models.GetAllMachines(false)
	if err != nil {
		beego.Error("Failed to get all machines: ", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	// Get the machines!
	machines := make([]*models.Machine, 0, len(allMachines))
	if !this.IsAdmin(ruid) {
		permissions, err := models.GetUserPermissions(ruid)
		if err != nil {
			beego.Error("Failed to get user machine permissions: ", err)
			this.CustomAbort(500, "Internal Server Error")
		}
		for _, permission := range *permissions {
			for _, machine := range allMachines {
				if machine.Id == permission.MachineId {
					machines = append(machines, machine)
					break
				}
			}
		}
	} else {
		machines = allMachines
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
	suid, err := this.GetSessionUserId()
	if err != nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Unauthorized")
	}

	if !this.IsAdmin() && suid != ruid {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Unauthorized")
	}

	startTime, err := purchases.GetUserStartTime(suid)
	if err != nil {
		beego.Error("GetUserStartTime:", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	startTime = startTime.Add(-86400 * time.Second)

	endTime := time.Now().Add(86400 * 30 * time.Second)
	invoice, err := invoices.CalculateSummary(startTime, endTime)
	if err != nil {
		beego.Error("Calculate invoice summary:", err)
	}

	var userSummary *invoices.UserSummary

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
	beego.Info("uids:", this.GetString("uids"))
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
		user, err := models.GetUser(uid)
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
	suid, err := this.GetSessionUserId()
	if err != nil {
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
	suid, err := this.GetSessionUserId()
	if err != nil {
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
