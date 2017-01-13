package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/newsletters:Controller"] = append(beego.GlobalControllerRouter["github.com/FabLabBerlin/localmachines/controllers/newsletters:Controller"],
		beego.ControllerComments{
			Method: "EasylabDev",
			Router: `/easylab_dev`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}
