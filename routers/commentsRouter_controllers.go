package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:aid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			Method: "Close",
			Router: `/:aid/close`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:rid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			Method: "GetLast",
			Router: `/:uid/last`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			Method: "GetActive",
			Router: `/active`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			Method: "Start",
			Router: `/start`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"],
		beego.ControllerComments{
			Method: "Archive",
			Router: `/:id/archive`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:CategoriesController"],
		beego.ControllerComments{
			Method: "Unarchive",
			Router: `/:id/unarchive`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:DebugController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:DebugController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"],
		beego.ControllerComments{
			Method: "GetCustomers",
			Router: `/customer`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"],
		beego.ControllerComments{
			Method: "CreateCustomer",
			Router: `/customer`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"],
		beego.ControllerComments{
			Method: "UpdateCustomer",
			Router: `/customer/:customerid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"],
		beego.ControllerComments{
			Method: "DeleteCustomer",
			Router: `/customer/:customerid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FeedbackController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FeedbackController"],
		beego.ControllerComments{
			Method: "PostFeedback",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:mid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/:mid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			Method: "SetArchived",
			Router: `/:mid/set_archived`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			Method: "GetAllRunning",
			Router: `/all_running`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			Method: "ArchiveProduct",
			Router: `/:productId/archive`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			Method: "ArchivePurchase",
			Router: `/:purchaseId/archive`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:rid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/:rid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:rid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			Method: "Cancel",
			Router: `/:id/cancel`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:rid([0-9]+)`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			Method: "ICalendar",
			Router: `/icalendar`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"],
		beego.ControllerComments{
			Method: "GetFastbillTemplates",
			Router: `/fastbill_templates`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"],
		beego.ControllerComments{
			Method: "GetTermsUrl",
			Router: `/terms_url`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"],
		beego.ControllerComments{
			Method: "GetVatPercent",
			Router: `/vat_percent`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:TutoringsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:TutoringsController"],
		beego.ControllerComments{
			Method: "Start",
			Router: `/:id/start`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:TutoringsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:TutoringsController"],
		beego.ControllerComments{
			Method: "Stop",
			Router: `/:id/stop`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:UserLocationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:UserLocationsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
