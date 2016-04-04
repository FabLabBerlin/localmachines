package clients

import (
	"github.com/FabLabBerlin/localmachines/controllers"
)

type Admin struct {
	controllers.Controller
}

// @Title Get
// @Description Get Admin Panel
// @Success 200 Redirect
// @Failure	500	Internal Server Error
// @router / [get]
func (c *Admin) Get() {
	get(&c.Controller, "admin")
}
