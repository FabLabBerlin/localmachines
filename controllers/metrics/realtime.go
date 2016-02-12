package metrics

import (
	"github.com/FabLabBerlin/localmachines/models/metrics"
	"github.com/astaxie/beego"
)

// @Title Get Realtime
// @Description Get realtime metrics
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /realtime [get]
func (c *Controller) GetRealtime() {
	if !c.IsAdmin() {
		c.CustomAbort(401, "Not authorized")
	}
	rt, err := metrics.NewRealtime()
	if err != nil {
		beego.Info(err.Error())
		c.CustomAbort(500, "Internal Server Error")
	}
	c.Data["json"] = rt
	c.ServeJSON()
}
