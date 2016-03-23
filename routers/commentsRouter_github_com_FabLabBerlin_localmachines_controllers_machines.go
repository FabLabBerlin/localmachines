package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"Get",
			`/:mid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"Create",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"Update",
			`/:mid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"PostImage",
			`/:mid/image`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"ReportBroken",
			`/:mid/report_broken`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"UnderMaintenanceOn",
			`/:mid/under_maintenance/on`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"UnderMaintenanceOff",
			`/:mid/under_maintenance/off`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"Search",
			`/search`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/machines:Controller"],
		beego.ControllerComments{
			"ApplyConfig",
			`/:mid/apply_config`,
			[]string{"post"},
			nil})

}
