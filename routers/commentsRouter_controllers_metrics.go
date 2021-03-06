package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			Method: "GetActivations",
			Router: `/activations`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			Method: "GetHeatmap",
			Router: `/heatmap`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			Method: "GetMachineCapacities",
			Router: `/machine_capacities`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			Method: "GetMachineEarnings",
			Router: `/machine_earnings`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			Method: "GetMemberships",
			Router: `/memberships`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			Method: "GetRetention",
			Router: `/retention`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			Method: "GetRealtime",
			Router: `/realtime`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
