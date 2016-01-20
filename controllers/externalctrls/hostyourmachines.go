package externalctrls

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models/external"
	"github.com/astaxie/beego"
)

type HostYourMachinesController struct {
	controllers.Controller
}

func (this *HostYourMachinesController) Prepare() {
	beego.Info("Skipping global login check")
}

// @Title Host Your Machines
// @Description Stores information provided by a potential host into the database
// @Param	first_name		query 	string	true		"First name"
// @Param	last_name		query 	string	true		"Last name"
// @Param	email			query 	string	true		"Email"
// @Param	location		query 	string	true		"Location"
// @Param	organization	query 	string	false		"Organization"
// @Param	phone			query 	string	false		"Phone"
// @Param	comments		query 	string	false		"Comments"
// @Success 200
// @Failure 500 Internal server error
// @router /hostyourmachines [post]
func (this *HostYourMachinesController) HostYourMachines() {
	host := external.Host{}
	host.FirstName = this.GetString("first_name")
	host.LastName = this.GetString("last_name")
	host.Email = this.GetString("email")
	host.Location = this.GetString("location")
	host.Organization = this.GetString("organization")
	host.Phone = this.GetString("phone")
	host.Comments = this.GetString("comments")

	id, err := host.Add()
	if err != nil {
		beego.Error(err)
		this.CustomAbort(500, "Failed to save host")
	}
	beego.Trace("Created host with ID", id)

	this.ServeJson()
}
