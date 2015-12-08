package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
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
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	settings, err := models.GetAllSettings()
	if err != nil {
		beego.Error("Failed to get all settings:", err)
		this.CustomAbort(500, "Failed to get all settings")
	}

	this.Data["json"] = settings.Data
	this.ServeJson()
}

// @Title Post
// @Description post settings
// @Success 201 {object} models.User
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router / [post]
func (this *SettingsController) Post() {
	if !this.IsAdmin() {
		beego.Error("Unauthorized attempt to change setting")
		this.CustomAbort(401, "Unauthorized")
	}

	var settings []*models.Setting

	dec := json.NewDecoder(this.Ctx.Request.Body)
	defer this.Ctx.Request.Body.Close()
	if err := dec.Decode(&settings); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update Settings")
	}

	for _, setting := range settings {
		var err error
		var msg string
		if setting.Id == 0 {
			msg = "new setting"
			_, err = models.CreateSetting(setting)
		} else {
			msg = "updating existing setting"
			err = setting.Update()
		}
		if err != nil {
			beego.Error(msg, err)
			this.CustomAbort(500, msg)
		}
	}

	this.ServeJson()
}
