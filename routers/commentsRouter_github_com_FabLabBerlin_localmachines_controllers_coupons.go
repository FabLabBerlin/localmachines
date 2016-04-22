package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/coupons:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/coupons:Controller"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/coupons:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/coupons:Controller"],
		beego.ControllerComments{
			"Generate",
			`/`,
			[]string{"post"},
			nil})

}
