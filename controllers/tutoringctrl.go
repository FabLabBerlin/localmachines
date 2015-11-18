package controllers

import (
	"encoding/json"
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

// @Title Put Tutor
// @Description Add or Update tutor
// @Success 200
// @Failure 500	Failed to update tutor
// @Failure 401 Not authorized
// @router /tutor [put]
func (this *TutoringController) PutTutor() {
	var msg string

	if !this.IsAdmin() {
		msg = "Not authorized"
		beego.Error(msg)
		this.CustomAbort(401, msg)
	}

	// Get body as array of models.User
	// Attempt to decode passed json
	jsonDecoder := json.NewDecoder(this.Ctx.Request.Body)
	tutor := models.Tutor{}

	err := jsonDecoder.Decode(&tutor)
	if err != nil {
		msg = "Failed to decode json"
		beego.Error(msg)
		this.CustomAbort(500, msg)
	}

	isNewTutor := false
	if tutor.Id == 0 {
		isNewTutor = true
	}

	if isNewTutor {
		err = models.CreateNewTutor(&tutor)
		if err != nil {
			msg = "Failed to create tutor"
			beego.Error(msg)
			this.CustomAbort(500, msg)
		}
	} else {
		/*
			err = models.UpdateTutor(&tutor)
			if err != nil {
				msg = "Failed to update tutor"
				beego.Error(msg)
				this.CustomAbort(500, msg)
			}
		*/
	}

	this.ServeJson()
}
