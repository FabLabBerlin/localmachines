package routers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/controllers"
)

func init() {
	// Route root request, serve Angular JS app
	beego.Router("/", &controllers.MainController{})

	// Route back office request
	beego.Router("/backoffice", &controllers.BackOfficeMainController{})

	// Route single method requests
	beego.Router("/api/login", &controllers.LoginController{}, "post:Login")
	beego.Router("/api/logout", &controllers.LogoutController{}, "*:Logout")
	beego.Router("/api/users", &controllers.UsersController{}, "get:GetUsers")
	beego.Router("/api/machines", &controllers.MachinesController{}, "get:GetMachines")

	// Route activation requests
	beego.Router("/api/activations", &controllers.ActivationsController{}, "get:GetActivations")
	beego.Router("/api/activations", &controllers.ActivationsController{}, "post:CreateActivation")
	beego.Router("/api/activations", &controllers.ActivationsController{}, "put:CloseActivation")
}
