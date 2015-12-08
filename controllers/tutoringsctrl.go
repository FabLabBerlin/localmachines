package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models/purchases"
)

type TutoringsController struct {
	Controller
}

// @Title Start
// @Description Start tutoring
// @Param	id		path 	int	true		"Tutoring ID"
// @Success 200
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:id/start [post]
func (this *TutoringsController) Start() {
	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get :id variable")
		this.CustomAbort(400, "Incorrect id")
	}

	err = purchases.StartTutoringPurchase(id)
	if err != nil {
		beego.Error("Failed to start tutoring purchase:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.ServeJson()
}

// @Title Stop
// @Description Stop tutoring
// @Param	id		path 	int	true		"Tutoring ID"
// @Success 200
// @Failure	400	Bad Request
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /:id/stop [post]
func (this *TutoringsController) Stop() {
	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get :id variable")
		this.CustomAbort(400, "Incorrect id")
	}

	err = purchases.StopTutoringPurchase(id)
	if err != nil {
		beego.Error("Failed to stop tutoring purchase:", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	this.ServeJson()
}
