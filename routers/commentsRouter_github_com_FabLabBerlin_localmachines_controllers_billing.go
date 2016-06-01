package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"GetMonth",
			`/months/:year/:month`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"GetUser",
			`/months/:year/:month/users/:uid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"GetStatuses",
			`/months/:year/:month/users/:uid/statuses`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"CreateDraft",
			`/months/:year/:month/users/:uid/draft`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"Send",
			`/months/:year/:month/users/:uid/invoices/:id/send`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"Update",
			`/months/:year/:month/users/:uid/update`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"SyncFastbillInvoices",
			`/months/:year/:month/sync`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"GetAll",
			`/monthly_earnings`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"Create",
			`/monthly_earnings`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"CreateDrafts",
			`/monthly_earnings/:iid/create_drafts`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/billing:Controller"],
		beego.ControllerComments{
			"DownloadExcelExport",
			`/monthly_earnings/:id/download_excel`,
			[]string{"get"},
			nil})

}
