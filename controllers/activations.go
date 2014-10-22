package controllers

import (
	"github.com/astaxie/beego"
)

type ActivationsController struct {
	beego.Controller
}

func (this *ActivationsController) CreateActivation() {
	response := struct{ Status string }{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}

func (this *ActivationsController) GetActivations() {
	type Item struct {
		Id        int
		MachineId int
		UserId    int
	}
	response := struct{ Activations []Item }{
		[]Item{
			Item{1, 2, 3},
			Item{3, 2, 1}}}
	this.Data["json"] = &response
	this.ServeJson()
}

func (this *ActivationsController) UpdateActivation() {
	response := struct{ Status string }{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}
