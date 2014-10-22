package controllers

import (
	"github.com/astaxie/beego"
)

type UsersController struct {
	beego.Controller
}

func (this *UsersController) GetUsers() {
	type User struct {
		Id    int
		Name  string
		Email string
	}
	response := struct{ Users []User }{
		[]User{
			User{1, "John", "john@aber.du"},
			User{2, "Mollie", "mollie@alles.gut"}}}
	this.Data["json"] = &response
	this.ServeJson()
}
