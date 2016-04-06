package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/locations:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/locations:Controller"],
		beego.ControllerComments{
			"Create",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/locations:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/locations:Controller"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/locations:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/locations:Controller"],
		beego.ControllerComments{
			"Get",
			`/:lid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/locations:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/locations:Controller"],
		beego.ControllerComments{
			"PostLocalIp",
			`/:lid/local_ip`,
			[]string{"post"},
			nil})

}
