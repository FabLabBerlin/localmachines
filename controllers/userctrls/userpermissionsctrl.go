package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/astaxie/beego"
)

type UserPermissionsController struct {
	Controller
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
	this.ServeJSON()
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
	this.ServeJSON()
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
	this.ServeJSON()
}

// @Title Get User Machine Permissions
// @Description Get user machine permissions
// @Param	uid		path 	int	true	"User ID"
// @Success 200	models.Permission
// @Failure	403	Failed to update permissions
// @Failure	401	Not authorized
// @router /:uid/permissions [get]
func (this *UserPermissionsController) GetUserPermissions() {

	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(400, "Wrong uid in url or not authorized")
	}

	permissions, err := models.GetUserPermissions(uid)
	if err != nil {
		beego.Error("Failed to get user permissions")
		this.CustomAbort(403, "Failed to get permissions")
	}

	this.Data["json"] = permissions
	this.ServeJSON()
}

// @Title GetUserMachinePermissions
// @Description Get current saved machine permissions
// @Param	uid		path 	int	true		"User ID"
// @Success 200 {object} models.Machine
// @Failure	500	Internal Server Error
// @Failure	401	Unauthorized
// @router /:uid/machinepermissions [get]
func (this *UserPermissionsController) GetUserMachinePermissions() {

	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(400, "Wrong uid in url or not authorized")
	}

	// We need to get machine permissions first and then the machines
	machines := make([]*models.Machine, 0, 20)
	permissions, err := models.GetUserPermissions(uid)
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
	this.ServeJSON()
}
