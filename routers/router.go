package routers

import (
	"github.com/kr15h/fabsmith/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
