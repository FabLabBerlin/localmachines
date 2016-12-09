package userctrls

import (
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
)

// @Title ByNfcId
// @Description Post nfc id
// @Body	nfcId		body 	string	""		""
// @Success 200
// @Failure	403	Failed to get user
// @Failure	401	Not authorized
// @router /by_nfc_id [post]
func (this *UsersController) GetByNfcId() {
	locId, authorized := this.GetLocIdApi()
	if !authorized {
		beego.Error("not having api user role")
		this.Fail(400)
	}

	if locId != 15 {
		beego.Error("only allowed for location 15 to test")
		this.Fail(400)
	}

	nfcId := this.GetString("nfcId")

	uid, err := users.AuthGetByNfcId(nfcId)
	if err != nil {
		beego.Error("Unable to get nfcId: ", err)
		this.Fail(400)
	}

	uls, err := user_locations.GetAllForUser(uid)
	if err != nil {
		beego.Error("Unable to get user locations: ", err)
		this.Fail(500)
	}

	found := false
	for _, ul := range uls {
		if ul.LocationId == locId {
			found = true
			break
		}
	}

	if !found {
		this.Fail(403)
	}

	u, err := users.GetUser(uid)
	if err != nil {
		beego.Error("Unable to get user: ", err)
		this.CustomAbort(500, "Unable to get user")
	}

	this.Data["json"] = u
	this.ServeJSON()
}

// @Title PostNfcId
// @Description Post nfc id
// @Param	uid		path 	int	true		"User ID"
// @Success 200
// @Failure	403	Failed to get user
// @Failure	401	Not authorized
// @router /:uid/nfc_id [post]
func (this *UsersController) PostNfcId() {
	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(400, "Wrong uid in url or not authorized")
	}

	err := users.AuthSetNfcId(uid, this.GetString("nfcId"))
	if err != nil {
		beego.Error("Unable to set nfcId: ", err)
		this.CustomAbort(403, "Unable to update nfcId")
	}

	this.Data["json"] = models.StatusResponse{
		Status: "Password changed successfully!",
	}
	this.ServeJSON()
}
