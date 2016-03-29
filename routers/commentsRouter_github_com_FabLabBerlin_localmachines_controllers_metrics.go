package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			"GetActivations",
			`/activations`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/metrics:Controller"],
		beego.ControllerComments{
			"GetRealtime",
			`/realtime`,
			[]string{"get"},
			nil})

}