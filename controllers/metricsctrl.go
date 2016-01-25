package controllers

import (
	"github.com/FabLabBerlin/localmachines/models/metrics"
	"github.com/astaxie/beego"
)

type MetricsController struct {
	Controller
}

// @Title Get All
// @Description Get all metrics
// @Success 200
// @Failure	500	Failed to get metrics
// @Failure	401	Not authorized
// @router / [get]
func (this *MetricsController) GetAll() {
	// Only admin can use this API call
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	data, err := metrics.FetchData()
	if err != nil {
		beego.Error("Failed to get metrics data:", err)
		this.CustomAbort(500, "Failed to get metrics data")
	}

	resp, err := metrics.NewResponse(data)
	if err != nil {
		beego.Error("Failed to get metrics response:", err)
		this.CustomAbort(500, "Failed to get metrics response")
	}

	this.Data["json"] = resp
	this.ServeJSON()
}
