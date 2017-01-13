package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "GetAllInvoices",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "GetInvoice",
			Router: `/invoices/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "GetMonth",
			Router: `/months/:year/:month`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "CreateDraft",
			Router: `/invoices/:id/draft`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "Cancel",
			Router: `/invoices/:id/cancel`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "Complete",
			Router: `/invoices/:id/complete`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "Send",
			Router: `/invoices/:id/send`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "SyncFastbillInvoices",
			Router: `/users/:uid/sync`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/monthly_earnings`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/monthly_earnings`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "CreateDrafts",
			Router: `/monthly_earnings/:iid/create_drafts`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			Method: "DownloadExcelExport",
			Router: `/monthly_earnings/:id/download_excel`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
