package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			"Get",
			`/:aid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			"Put",
			`/:rid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			"GetActive",
			`/active`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			"Start",
			`/start`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			"Close",
			`/:aid/close`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ActivationsController"],
		beego.ControllerComments{
			"PostFeedback",
			`/:aid/feedback`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:DebugController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:DebugController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"],
		beego.ControllerComments{
			"GetCustomers",
			`/customer`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"],
		beego.ControllerComments{
			"CreateCustomer",
			`/customer`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"],
		beego.ControllerComments{
			"UpdateCustomer",
			`/customer/:customerid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FastBillController"],
		beego.ControllerComments{
			"DeleteCustomer",
			`/customer/:customerid`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FeedbackController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:FeedbackController"],
		beego.ControllerComments{
			"PostFeedback",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:InvoicesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:InvoicesController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:InvoicesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:InvoicesController"],
		beego.ControllerComments{
			"Create",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:InvoicesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:InvoicesController"],
		beego.ControllerComments{
			"CreateDrafts",
			`/:iid/create_drafts`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:InvoicesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:InvoicesController"],
		beego.ControllerComments{
			"DownloadExcelExport",
			`/:id/download_excel`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MachineTypeController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MachineTypeController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			"Create",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			"Get",
			`/:mid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:MembershipsController"],
		beego.ControllerComments{
			"Update",
			`/:mid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			"Create",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			"Get",
			`/:id`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			"Put",
			`/:id`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ProductsController"],
		beego.ControllerComments{
			"ArchiveProduct",
			`/:productId/archive`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			"Create",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			"Get",
			`/:id`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			"Put",
			`/:id`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:PurchasesController"],
		beego.ControllerComments{
			"ArchivePurchase",
			`/:purchaseId/archive`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			"Get",
			`/:rid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			"Create",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			"Update",
			`/:rid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationRulesController"],
		beego.ControllerComments{
			"Delete",
			`/:rid`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			"Get",
			`/:rid([0-9]+)`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			"Create",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			"Put",
			`/:id`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:ReservationsController"],
		beego.ControllerComments{
			"ICalendar",
			`/icalendar`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:SettingsController"],
		beego.ControllerComments{
			"GetTermsUrl",
			`/terms_url`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:TutoringsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:TutoringsController"],
		beego.ControllerComments{
			"Start",
			`/:id/start`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:TutoringsController"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers:TutoringsController"],
		beego.ControllerComments{
			"Stop",
			`/:id/stop`,
			[]string{"post"},
			nil})

}
