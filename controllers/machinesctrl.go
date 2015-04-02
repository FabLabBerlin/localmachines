package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

type MachinesController struct {
	Controller
}

// @Title GetAll
// @Description Get all machines
// @Success 200 {object} models.Machine
// @Failure	403	Failed to get all machines
// @Failure	401 Not authorized
// @router / [get]
func (this *MachinesController) GetAll() {

}

// @Title Get
// @Description Get machine by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 {object} models.Machine
// @Failure	403	Failed to get machine
// @Failure	401	Not authorized
// @router /:mid [get]
func (this *MachinesController) Get() {

}

// @Title Create
// @Description Create machine
// @Param	mname	query	string	true	"Machine Name"
// @Success 200 {object} models.MachineCreatedResponse
// @Failure	403	Failed to create machine
// @Failure	401	Not authorized
// @router / [post]
func (this *MachinesController) Create() {
	machineName := this.GetString("mname")
	beego.Trace(machineName)

	var err error

	uid := this.GetSession(SESSION_FIELD_NAME_USER_ID).(int64)

	// Check if user is admin or staff
	userRoles, err := models.GetUserRoles(uid)
	if !userRoles.Admin && !userRoles.Staff {
		beego.Error("Not authorized to create machine")
		this.CustomAbort(401, "Not authorized")
	}

	// All clear - create machine in the database
	o := orm.NewOrm()
	machine := models.Machine{Name: machineName, Available: true}
	id, err := o.Insert(&machine)
	if err == nil {
		response := models.MachineCreatedResponse{MachineId: id}
		this.Data["json"] = response
		this.ServeJson()
	} else {
		beego.Error("Cannot create user: ", err)
		this.CustomAbort(403, "Failed to create machine")
	}
}
