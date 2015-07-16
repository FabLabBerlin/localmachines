package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type PaymillController struct {
	Controller
}

// Override our custom root controller's Prepare method as it is checking
// if we are logged in and we don't want that here at this point
func (this *PaymillController) Prepare() {
	beego.Info("Skipping global login check")
}

// @Title Get All
// @Description Get all activations
// @Success 200 {object} models.GetPublicKeyResponse
// @router /public [get]
func (this *PaymillController) GetAll() {
	resp := models.GetPublicKeyResponse{
		PublicKey: beego.AppConfig.String("paymillpublic"),
	}

	this.Data["json"] = resp
	this.ServeJson()
}
