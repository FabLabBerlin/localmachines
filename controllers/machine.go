package controllers

/*
import (
	"github.com/kr15h/fabsmith/models"
	"github.com/astaxie/beego"
)
*/

type MachineController struct {
	Controller
}

// @Title Get
// @Description Get machine by machine ID
// @Param	mid		path 	int	true		"Machine ID"
// @Success 200 {object} models.Machine
// @Failure	403	Failed to get machine
// @Failure	401	Not authorized
// @router /:mid [get]
func (this *MachineController) Get() {

} 

// @Title GetAll
// @Description Get all machines
// @Success 200 {object} models.Machine
// @Failure	403	Failed to get all machines
// @Failure	401 Not authorized
// @router / [get]
func (this *MachineController) GetAll() {

}