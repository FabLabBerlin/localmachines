package controllers

import (
	"github.com/astaxie/beego"
)

type UsersController struct {
	beego.Controller
}

type UsersResponseUser struct {
	Id    int
	Name  string
	Email string
}

type UsersResponse struct {
	Users []UsersResponseUser
}

func (this *UsersController) Get() {
	response := &UsersResponse{
		[]UsersResponseUser{
			UsersResponseUser{1, "John", "john@aber.du"},
			UsersResponseUser{2, "Mollie", "mollie@alles.gut"}}}
	this.Data["json"] = &response
	this.ServeJson()
}
