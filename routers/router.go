package routers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/login", &controllers.LoginController{})
	beego.Router("/api/logout", &controllers.LogoutController{}, "*:Any")
	beego.Router("/api/users", &controllers.UsersController{})
	beego.Router("/api/machines", &controllers.MachinesController{})
	beego.Router("/api/activations", &controllers.ActivationsController{})
}
