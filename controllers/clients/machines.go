package clients

import (
	"github.com/FabLabBerlin/localmachines/controllers"
)

type Machines struct {
	controllers.Controller
}

// @Title Get
// @Description Get Machines Page
// @Success 200 Redirect
// @Failure	500	Internal Server Error
// @router / [get]
func (c *Machines) Get() {
	get(&c.Controller, "machines")
}
