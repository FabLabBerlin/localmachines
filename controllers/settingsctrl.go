package controllers

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/models/settings"
	"github.com/astaxie/beego"
)

type SettingsController struct {
	Controller
}

// @Title Get All
// @Description Get all settings
// @Success 200
// @Failure	500	Failed to get settings
// @Failure	401	Not authorized
// @router / [get]
func (this *SettingsController) GetAll() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	settings, err := settings.GetAllAt(locId)
	if err != nil {
		beego.Error("Failed to get all settings:", err)
		this.CustomAbort(500, "Failed to get all settings")
	}

	this.Data["json"] = settings.Data
	this.ServeJSON()
}

// @Title Post
// @Description post settings
// @Success 201 {object} models.User
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router / [post]
func (this *SettingsController) Post() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	var allSettings []*settings.Setting

	dec := json.NewDecoder(this.Ctx.Request.Body)
	defer this.Ctx.Request.Body.Close()
	if err := dec.Decode(&allSettings); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update Settings")
	}

	for _, setting := range allSettings {
		var err error
		var msg string

		setting.LocationId = locId

		if setting.Id == 0 {
			msg = "new setting"
			_, err = settings.Create(setting)
		} else {
			msg = "updating existing setting"
			err = setting.Update()
		}
		if err != nil {
			beego.Error(msg, err)
			this.CustomAbort(500, msg)
		}
	}

	this.ServeJSON()
}
