package routers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/controllers"
)

func init() {
	
	// Set main redirect in the MainController
	beego.Router("/", &controllers.MainController{})

	// Route API requests
	beego.Router("/api/login", &controllers.LoginController{}, "post:Login")
	beego.Router("/api/logout", &controllers.LogoutController{}, "*:Logout")
	beego.Router("/api/users", &controllers.UsersController{}, "get:GetUsers")
	beego.Router("/api/machines", &controllers.MachinesController{}, "get:GetMachines")

	beego.Router("/api/activations", &controllers.ActivationsController{}, "get:GetActivations")
	beego.Router("/api/activations", &controllers.ActivationsController{}, "post:CreateActivation")
	beego.Router("/api/activations", &controllers.ActivationsController{}, "put:CloseActivation")
}
