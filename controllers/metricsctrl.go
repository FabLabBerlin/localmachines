package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
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

	data, err := models.FetchMetricsData()
	if err != nil {
		beego.Error("Failed to get metrics data:", err)
		this.CustomAbort(500, "Failed to get metrics data")
	}

	resp, err := models.NewMetricsResponse(data)
	if err != nil {
		beego.Error("Failed to get metrics response:", err)
		this.CustomAbort(500, "Failed to get metrics response")
	}

	this.Data["json"] = resp
	this.ServeJson()
}
