package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/astaxie/beego"
)

type UserPermissionsController struct {
	controllers.Controller
}

// @Title CreateUserPermission
// @Description Create a permission for a user to allow him/her to use a machine
// @Param	uid		path 	int	true		"User ID"
// @Param	mid		query 	int	true		"Machine ID"
// @Success 200 string ok
// @Failure	403	Failed to create permission
// @Failure	401	Not authorized
// @router /:uid/permissions [post]
func (this *UserPermissionsController) CreateUserPermission() {

	// TODO: think about bulk permission creation or another way
	// 		 that does not use a separate table maybe.

	// Only admin
	if !this.IsAdmin() {
		this.CustomAbort(401, "Not authorized")
	}

	// Get user ID for the machine permission to be made
	userId, err := this.GetInt64(":uid")
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
func (this *UserPermissionsController) DeleteUserPermission() {

	// Only admin
	if !this.IsAdmin() {
		this.CustomAbort(401, "Not authorized")
	}

	// Get user ID for the machine permission to be made
	userId, err := this.GetInt64(":uid")
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
func (this *UserPermissionsController) UpdateUserPermissions() {

	// Only admin can do this
	if !this.IsAdmin() {
		this.CustomAbort(401, "Not authorized")
	}

	// Get request user ID
	userId, err := this.GetInt64(":uid")
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
func (this *UserPermissionsController) GetUserPermissions() {

	// Check if logged in
	suid, err := this.GetSessionUserId()
	if err != nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Not logged in")
	}

	// Get session user ID as int64
	var found bool
	if !found {
		beego.Error("Could not cast session ID to int64")
		this.CustomAbort(403, "Failed to get permissions")
	}

	// Get request user ID
	var userId int64
	userId, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID:", err)
		this.CustomAbort(403, "Failed to get permissions")
	}

	// Allow to get other user permissions only if admin
	if userId != suid {
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

// @Title GetUserMachinePermissions
// @Description Get current saved machine permissions
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.Machine
// @Failure	500	Internal Server Error
// @Failure	401	Unauthorized
// @router /:uid/machinepermissions [get]
func (this *UserPermissionsController) GetUserMachinePermissions() {

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
		beego.Error("Failed to get :uid")
		this.CustomAbort(500, "Internal Server Error")
	}

	// If the session user ID and the request user ID does not match
	// check whether the session user is an admin, return if not
	if suid != ruid {
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
