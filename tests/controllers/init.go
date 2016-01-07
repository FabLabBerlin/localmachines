package controllerTest

import (
	"github.com/FabLabBerlin/localmachines/routers"
	"github.com/FabLabBerlin/localmachines/tests/setup"
)

func init() {
	setup.ConfigDir("/..")
	setup.ConfigDB()
	routers.Init()
}
