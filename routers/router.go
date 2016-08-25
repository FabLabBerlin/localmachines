// @APIVersion 0.1.0
// @Title Fabsmith API
// @Description Makerspace machine management
// @Contact krisjanis.rijnieks@gmail.com
package routers

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/controllers/billing"
	"github.com/FabLabBerlin/localmachines/controllers/clients"
	"github.com/FabLabBerlin/localmachines/controllers/coupons"
	"github.com/FabLabBerlin/localmachines/controllers/custom_url"
	"github.com/FabLabBerlin/localmachines/controllers/locations"
	"github.com/FabLabBerlin/localmachines/controllers/machines"
	"github.com/FabLabBerlin/localmachines/controllers/metrics"
	"github.com/FabLabBerlin/localmachines/controllers/newsletters"
	"github.com/FabLabBerlin/localmachines/controllers/userctrls"
	locationModels "github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
)

// Init must be exportable for controller tests
func Init() {
	// Set main redirect in the MainController
	beego.Router("/", &controllers.MainController{})
	beego.Router("/crossdomain.xml", &controllers.CrossdomainController{})
	beego.Router("/apple-touch-icon.png", &controllers.AppleTouchIconController{})
	beego.Router("/favicon.png", &controllers.FaviconController{})
	beego.Router("/admin", &clients.Admin{})
	beego.Router("/machines", &clients.Machines{})
	beego.Router("/product", &clients.Machines{})
	beego.Router("/logout", &controllers.LogoutController{})
	// No need to create a router for favicon.ico and robots.txt as they are
	// handled by beego router.go

	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/activations",
			beego.NSInclude(
				&controllers.ActivationsController{},
			),
		),
		beego.NSNamespace("/billing",
			beego.NSInclude(
				&billing.Controller{},
			),
		),
		beego.NSNamespace("/coupons",
			beego.NSInclude(
				&coupons.Controller{},
			),
		),
		beego.NSNamespace("/debug",
			beego.NSInclude(
				&controllers.DebugController{},
			),
		),
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
		beego.NSNamespace("/user_locations",
			beego.NSInclude(
				&controllers.UserLocationsController{},
			),
		),
		beego.NSNamespace("/users",
			beego.NSInclude(
				&userctrls.ForgotPassword{},
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
		beego.NSNamespace("/newsletters",
			beego.NSInclude(
				&newsletters.Controller{},
			),
		),
		beego.NSNamespace("/products",
			beego.NSInclude(
				&controllers.ProductsController{},
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
	beego.Info("Router: loading locations...")
	locs, err := locationModels.GetAll()
	if err != nil {
		panic(err.Error())
	}
	for _, l := range locs {
		ns := beego.NewNamespace(l.Title,
			beego.NSInclude(
				&custom_url.Controller{},
			),
		)
		beego.AddNamespace(ns)
	}
	beego.Info("Router: custom urls added...")
}
