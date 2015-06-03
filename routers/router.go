// @APIVersion 0.1.0
// @Title Fabsmith API
// @Description Makerspace machine management
// @Contact krisjanis.rijnieks@gmail.com
package routers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/controllers"
)

func init() {

	// Set main redirect in the MainController
	beego.Router("/", &controllers.MainController{})
	beego.Router("/crossdomain.xml", &controllers.CrossdomainController{})

	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/users",
			beego.NSInclude(
				&controllers.UsersController{},
			),
		),
		beego.NSNamespace("/memberships",
			beego.NSInclude(
				&controllers.MembershipsController{},
			),
		),
		beego.NSNamespace("/machines",
			beego.NSInclude(
				&controllers.MachinesController{},
			),
		),
		beego.NSNamespace("/activations",
			beego.NSInclude(
				&controllers.ActivationsController{},
			),
		),
		beego.NSNamespace("/hexabus",
			beego.NSInclude(
				&controllers.HexabusController{},
			),
		),
		beego.NSNamespace("/invoices",
			beego.NSInclude(
				&controllers.InvoicesController{},
			),
		),
		beego.NSNamespace("/urlswitch",
			beego.NSInclude(
				&controllers.UrlSwitchController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
