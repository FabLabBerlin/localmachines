// @APIVersion 0.1.0
// @Title Fabsmith API
// @Description Makerspace machine management
// @Contact krisjanis.rijnieks@gmail.com
package routers

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/controllers/locations"
	"github.com/FabLabBerlin/localmachines/controllers/machines"
	"github.com/FabLabBerlin/localmachines/controllers/metrics"
	"github.com/FabLabBerlin/localmachines/controllers/userctrls"
	"github.com/astaxie/beego"
)

func init() {
	Init()
}

// Init must be exportable for controller tests
func Init() {
	// Set main redirect in the MainController
	beego.Router("/", &controllers.MainController{})
	beego.Router("/crossdomain.xml", &controllers.CrossdomainController{})
	beego.Router("/apple-touch-icon.png", &controllers.AppleTouchIconController{})
	beego.Router("/favicon.png", &controllers.FaviconController{})
	// No need to create a router for favicon.ico and robots.txt as they are
	// handled by beego router.go

	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/locations",
			beego.NSInclude(
				&locations.Controller{},
			),
		),
		beego.NSNamespace("/machines",
			beego.NSInclude(
				&machines.Controller{},
			),
		),
		beego.NSNamespace("/users",
			beego.NSInclude(
				&userctrls.UsersController{},
				&userctrls.UserDashboardController{},
				&userctrls.UserLocationsController{},
				&userctrls.UserMembershipsController{},
				&userctrls.UserPermissionsController{},
			),
		),
		beego.NSNamespace("/memberships",
			beego.NSInclude(
				&controllers.MembershipsController{},
			),
		),
		beego.NSNamespace("/machine_types",
			beego.NSInclude(
				&controllers.MachineTypeController{},
			),
		),
		beego.NSNamespace("/products",
			beego.NSInclude(
				&controllers.ProductsController{},
			),
		),
		beego.NSNamespace("/activations",
			beego.NSInclude(
				&controllers.ActivationsController{},
			),
		),
		beego.NSNamespace("/invoices",
			beego.NSInclude(
				&controllers.InvoicesController{},
			),
		),
		beego.NSNamespace("/fastbill",
			beego.NSInclude(
				&controllers.FastBillController{},
			),
		),
		beego.NSNamespace("/feedback",
			beego.NSInclude(
				&controllers.FeedbackController{},
			),
		),
		beego.NSNamespace("/metrics",
			beego.NSInclude(
				&metrics.Controller{},
			),
		),
		beego.NSNamespace("/reservations",
			beego.NSInclude(
				&controllers.ReservationsController{},
			),
		),
		beego.NSNamespace("/reservation_rules",
			beego.NSInclude(
				&controllers.ReservationRulesController{},
			),
		),
		beego.NSNamespace("/settings",
			beego.NSInclude(
				&controllers.SettingsController{},
			),
		),
		beego.NSNamespace("/purchases",
			beego.NSInclude(
				&controllers.PurchasesController{},
			),
		),
		beego.NSNamespace("/tutorings",
			beego.NSInclude(
				&controllers.TutoringsController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
