package controllers

import (
	"github.com/astaxie/beego"
)

type ActivationsController struct {
	beego.Controller
}

type ActivationsResponsePost struct {
	Status string
}

type ActivationsResponseGetItem struct {
	Id        int
	MachineId int
	UserId    int
}

type ActivationsResponseGet struct {
	Activations []ActivationsResponseGetItem
}

type ActivationsResponsePut struct {
	Status string
}

func (this *ActivationsController) Post() {
	response := &ActivationsResponsePost{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}

func (this *ActivationsController) Get() {
	response := &ActivationsResponseGet{
		[]ActivationsResponseGetItem{
			ActivationsResponseGetItem{1, 2, 3},
			ActivationsResponseGetItem{3, 2, 1}}}
	this.Data["json"] = &response
	this.ServeJson()
}

func (this *ActivationsController) Put() {
	response := &ActivationsResponsePut{"ok"}
	this.Data["json"] = &response
	this.ServeJson()
}
