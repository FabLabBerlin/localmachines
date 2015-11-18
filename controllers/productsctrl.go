package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type ProductsController struct {
	Controller
}

// @Title Get All
// @Description Get all products
// @Success 200
// @Failure	500	Failed to get products
// @Failure	401	Not authorized
// @router / [get]
func (this *ProductsController) GetAll() {
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	products, err := models.GetAllProducts()
	if err != nil {
		beego.Error("Failed to get all products:", err)
		this.CustomAbort(500, "Failed to get all products")
	}

	this.Data["json"] = products
	this.ServeJson()
}
