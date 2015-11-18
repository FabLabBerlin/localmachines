package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type TutoringController struct {
	Controller
}

// @Title Get Tutor List
// @Description Get all settings
// @Success 200	{object}	models.TutorList
// @Failure 500	Failed to get tutors
// @Failure 401	Not authorized
// @router /tutors [get]
func (this *TutoringController) GetAllTutors() {
	var msg string

	if !this.IsAdmin() {
		msg = "Not authorized"
		beego.Error(msg)
		this.CustomAbort(401, msg)
	}

	tutors, err := models.GetTutorList()
	if err != nil {
		msg = "Failed to get tutors"
		beego.Error(msg)
		this.CustomAbort(500, msg)
	}

	this.Data["json"] = tutors
	this.ServeJson()
}

// @Title Get All Tutoring Purchases
// @Description Get all tutoring purchases
// @Success 200	{object}	models.TutoringPurchaseList
// @Failure 500	Failed to get tutoring purchases
// @Failure 401 Not authorized
// @router /purchases [get]
func (this *TutoringController) GetAllTutoringPurchases() {
	var msg string

	if !this.IsAdmin() {
		msg = "Not authorized"
		beego.Error(msg)
		this.CustomAbort(401, msg)
	}

	purchases, err := models.GetTutoringPurchaseList()
	if err != nil {
		msg = "Failed to get tutoring purchases"
		beego.Error(msg)
		this.CustomAbort(500, msg)
	}

	this.Data["json"] = purchases
	this.ServeJson()
}
