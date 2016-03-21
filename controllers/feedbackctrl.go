package controllers

import (
	"github.com/FabLabBerlin/localmachines/lib/feedback"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
)

type FeedbackController struct {
	Controller
}

// @Title PostFeedback
// @Description Post feedback, either as ZenDesk ticket or as E-Mail
// @Param location query   int     true     "Location id"
// @Param name     query   string  true     "Full name of the customer"
// @Param email    query   string  true     "Email of the customer"
// @Param subject  query   string  true     "Subject of the Ticket"
// @Param message  query   string  true     "Message of the Ticket"
// @Success 200 {object}
// @Failure 500 Internal Server Error
// @Failure 401 Not authorized
// @router  / [post]
func (this *FeedbackController) PostFeedback() {
	locationId, err := this.GetInt64("location")
	if err != nil {
		this.CustomAbort(400, "Bad request")
	}
	location, err := locations.Get(locationId)
	if err != nil {
		beego.Error("locations get: %v", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	fb := feedback.Feedback{
		Location:       location,
		RequesterName:  this.GetString("name"),
		RequesterEmail: this.GetString("email"),
		Subject:        this.GetString("subject"),
		Message:        this.GetString("message"),
	}
	if err := fb.Send(); err != nil {
		beego.Error("Error sending feedback:", err)
		this.CustomAbort(500, "Internal Server Error")
	}
	this.ServeJSON()
}
