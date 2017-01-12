package userctrls

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
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

	// Get user ID for the machine permission to be made
	userId, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID")
		this.Abort("400")
	}

	// Get machine ID for the permission being made
	machineId, err := this.GetInt64("mid")
	if err != nil {
		beego.Error("Failed to get queried machine ID")
		this.Abort("400")
	}

	machine, err := machine.Get(machineId)
	if err != nil {
		beego.Error("get machine:", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(machine.LocationId) {
		this.Abort("401")
	}

	// Create a new user permission record in the database
	err = user_permissions.Create(userId, machineId)
	if err != nil {
		beego.Error("Failed to create machine user permission", err)
		this.Abort("500")
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
	// Get user ID for the machine permission to be made
	userId, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID:", err)
		this.Abort("400")
	}

	// Get machine ID for the permission being made
	machineId, err := this.GetInt64("mid")
	if err != nil {
		beego.Error("Failed to get queried machine ID")
		this.Abort("400")
	}

	machine, err := machine.Get(machineId)
	if err != nil {
		beego.Error("get machine:", err)
		this.Abort("500")
	}

	if !this.IsAdminAt(machine.LocationId) {
		this.Abort("401")
	}

	err = user_permissions.Delete(userId, machineId)
	if err != nil {
		beego.Error("Failed to delete permission:", err)
		this.Abort("500")
	}

	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title Update User Machine Permissions
// @Description Update user machine permissions
// @Param	uid			path 	int					true	"User ID"
// @Param	model		body	models.Permission	true	"Permissions Array"
// @Param	location	query	int					true	"Location ID"
// @Success 200	ok
// @Failure	401	Not authorized
// @Failure	403	Forbidden
// @Failure	500 Internal Server Error
// @router /:uid/permissions [put]
func (this *UserPermissionsController) UpdateUserPermissions() {

	// Get request user ID
	userId, err := this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get requested user ID:", err)
		this.Abort("400")
	}

	locationId, err := this.GetInt64("location")
	if err != nil {
		beego.Error("No location specified")
		this.Abort("400")
	}

	if !this.IsAdminAt(locationId) {
		this.Abort("401")
	}

	// Get body as array of models.Permission
	// Attempt to decode passed json
	jsonDecoder := json.NewDecoder(this.Ctx.Request.Body)
	permissions := []user_permissions.Permission{}
	if err = jsonDecoder.Decode(&permissions); err != nil {
		beego.Error("Failed to decode json:", err)
		this.Abort("400")
	}

	// Make sure that the user IDs of all the permissions are the same
	// as the request user ID and location ID
	for i := 0; i < len(permissions); i++ {
		permissions[i].LocationId = locationId
		permissions[i].UserId = userId
	}

	// Update permissions
	err = user_permissions.Update(userId, locationId, &permissions)
	if err != nil {
		beego.Error("Failed to update permissions:", err)
		this.Abort("500")
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

	permissions, err := user_permissions.Get(uid)
	if err != nil {
		beego.Error("Failed to get user permissions")
		this.CustomAbort(403, "Failed to get permissions")
	}

	this.Data["json"] = permissions
	this.ServeJSON()
}
