package controllers

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
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
// @router / [get]
func (this *SettingsController) GetAll() {
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
// @Success 201
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

// @Title GetTermsUrl
// @Description Get URL of Terms and Conditions
// @Success 200
// @Failure	400	Bad request
// @Failure	500	Internal Server Error
// @router /terms_url [get]
func (this *SettingsController) GetTermsUrl() {
	locId, err := this.GetInt64("location")
	if err != nil {
		beego.Error("parse int location:", err)
		this.CustomAbort(400, "Bad request")
	}

	s, err := settings.GetAllAt(locId)
	if err != nil {
		beego.Error("Failed to get all settings:", err)
		this.CustomAbort(500, "Failed to get all settings")
	}

	this.Data["json"] = s.GetString(locId, settings.TERMS_URL)
	this.ServeJSON()
}

// @Title GetVatPercent
// @Description Get VAT Percent
// @Success 200
// @Failure	400	Bad request
// @Failure	500	Internal Server Error
// @router /vat_percent [get]
func (this *SettingsController) GetVatPercent() {
	locId, err := this.GetInt64("location")
	if err != nil {
		beego.Error("parse int location:", err)
		this.CustomAbort(400, "Bad request")
	}

	s, err := settings.GetAllAt(locId)
	if err != nil {
		beego.Error("Failed to get all settings:", err)
		this.CustomAbort(500, "Failed to get all settings")
	}

	this.Data["json"] = s.GetFloat(locId, settings.VAT)
	this.ServeJSON()
}

// @Title GetFastbillTemplates
// @Description Get Fastbill Templates
// @Success 200
// @Failure	400	Bad request
// @Failure	500	Internal Server Error
// @router /fastbill_templates [get]
func (this *SettingsController) GetFastbillTemplates() {
	locId, authorized := this.GetLocIdAdmin()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	ts := []fastbill.Template{}
	if locId == 1 {
		var err error

		ts, err = fastbill.ListTemplates()
		if err != nil {
			beego.Error("list templates:", err)
			this.Abort("500")
		}
	}

	this.Data["json"] = ts
	this.ServeJSON()
}
