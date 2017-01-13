package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:mid`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/:mid`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "SetArchived",
			Router: `/:mid/set_archived`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "PostImage",
			Router: `/:mid/image`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "ReportBroken",
			Router: `/:mid/report_broken`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "UnderMaintenanceOn",
			Router: `/:mid/under_maintenance/on`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "UnderMaintenanceOff",
			Router: `/:mid/under_maintenance/off`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "Search",
			Router: `/search`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			Method: "ApplyConfig",
			Router: `/:mid/apply_config`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}
