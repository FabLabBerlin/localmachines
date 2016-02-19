package metrics

import (
	"github.com/FabLabBerlin/localmachines/models/metrics"
	"github.com/astaxie/beego"
	"strings"
)

// @Title Get Realtime
// @Description Get realtime metrics
// @Param	apikey		query	string	true	"API Key for Grafana consumption"
// @Param	location	query	string	true	"Location ID"
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /realtime [get]
func (c *Controller) GetRealtime() {
	locationId, _ := c.GetInt64("location")

	apiKey := beego.AppConfig.String("GrafanaApiKey")
	apiKey = strings.TrimSpace(apiKey)
	if len(apiKey) < 20 {
		beego.Error("GrafanaApiKey too short")
		c.CustomAbort(500, "Internal Server Error")
	}
	if c.GetString("apikey") != apiKey {
		c.CustomAbort(401, "Not authorized")
	}

	rt, err := metrics.NewRealtime(locationId)
	if err != nil {
		beego.Info(err.Error())
		c.CustomAbort(500, "Internal Server Error")
	}
	c.Data["json"] = rt
	c.ServeJSON()
}
