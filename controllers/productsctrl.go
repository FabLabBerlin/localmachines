package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models/products"
)

type ProductsController struct {
	Controller
}

// @Title Create
// @Description Create product with name of specified type
// @Param	name	query	string	true	"Space Name"
// @Param	type	query	string	true	"Product Type"
// @Success 200 {object} models.Space
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router / [post]
func (this *ProductsController) Create() {
	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized to create product")
		this.CustomAbort(401, "Not authorized")
	}

	name := this.GetString("name")
	productType := this.GetString("type")

	var product interface{}
	var err error

	switch productType {
	case products.PRODUCT_TYPE_CO_WORKING:
		product, err = products.CreateCoWorkingProduct(name)
		break
	case products.PRODUCT_TYPE_SPACE:
		product, err = products.CreateSpace(name)
		break
	case products.PRODUCT_TYPE_TUTOR:
		product, err = products.CreateTutor(&products.Tutor{})
		break
	default:
		err = fmt.Errorf("unknown product type")
	}

	if err != nil {
		beego.Error("Failed to create product", err, " (productType=", productType, ")")
		this.CustomAbort(500, "Failed to create product")
	}

	this.Data["json"] = product
	this.ServeJson()
}

// @Title Get
// @Description Get product by ID
// @Param	id		path 	int	true		"Space ID"
// @Param	type	query	string	true	"Product Type"
// @Success 200 string
// @Failure	400	Bad Request
// @Failure	500	Failed to get product
// @router /:id [get]
func (this *ProductsController) Get() {
	id, err := this.GetInt64(":id")
	if err != nil {
		beego.Error("Failed to get read id")
		this.CustomAbort(400, "Failed to get space")
	}

	productType := this.GetString("type")

	var product interface{}

	switch productType {
	case products.PRODUCT_TYPE_CO_WORKING:
		product, err = products.GetCoWorkingProduct(id)
		break
	case products.PRODUCT_TYPE_SPACE:
		product, err = products.GetSpace(id)
		break
	case products.PRODUCT_TYPE_TUTOR:
		product, err = products.GetTutor(id)
		break
	default:
		err = fmt.Errorf("unknown product type")
	}

	if err != nil {
		beego.Error("Failed to get product:", err, " (productType=", productType, ")")
		this.CustomAbort(500, "Failed to get product")
	}

	this.Data["json"] = product
	this.ServeJson()
}

// @Title Get All
// @Description Get all products [of specified type, if param present]
// @Param	type	query	string	true	"Product Type"
// @Success 200
// @Failure	500	Failed to get products
// @Failure	401	Not authorized
// @router / [get]
func (this *ProductsController) GetAll() {
	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	productType := this.GetString("type")

	var ps interface{}
	var err error

	switch productType {
	case products.PRODUCT_TYPE_CO_WORKING:
		ps, err = products.GetAllCoWorkingProducts()
		break
	case products.PRODUCT_TYPE_SPACE:
		ps, err = products.GetAllSpaces()
		break
	case products.PRODUCT_TYPE_TUTOR:
		ps, err = products.GetAllTutors()
		break
	case "":
		ps, err = products.GetAllProducts()
	default:
		err = fmt.Errorf("unknown product type")
	}

	if err != nil {
		beego.Error("Failed to get all products:", err, " (productType=", productType, ")")
		this.CustomAbort(500, "Failed to get all products")
	}

	this.Data["json"] = ps

	this.ServeJson()
}

// @Title Put
// @Description Update product
// @Param	type	query	string	true	"Product Type"
// @Success 201 string
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:rid [put]
func (this *ProductsController) Put() {
	if !this.IsAdmin() {
		beego.Error("Unauthorized attempt to update product")
		this.CustomAbort(401, "Unauthorized")
	}

	productType := this.GetString("type")

	var response interface{}
	var err error

	switch productType {
	case products.PRODUCT_TYPE_CO_WORKING:
		table := &products.CoWorkingProduct{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(table); err != nil {
			beego.Error("Failed to decode json:", err)
			this.CustomAbort(400, "Failed to update table")
		}

		if err = products.UpdateCoWorkingProduct(table); err == nil {
			response = table
		}
		break
	case products.PRODUCT_TYPE_SPACE:
		space := &products.Space{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(space); err != nil {
			beego.Error("Failed to decode json:", err)
			this.CustomAbort(400, "Failed to update space")
		}

		if err = products.UpdateSpace(space); err == nil {
			response = space
		}
		break
	case products.PRODUCT_TYPE_TUTOR:
		tutor := &products.Tutor{}
		dec := json.NewDecoder(this.Ctx.Request.Body)
		defer this.Ctx.Request.Body.Close()
		if err := dec.Decode(tutor); err != nil {
			beego.Error("Failed to decode json:", err)
			this.CustomAbort(400, "Failed to update tutor")
		}

		if err = products.UpdateTutor(tutor); err == nil {
			response = tutor
		}
		break
	default:
		err = fmt.Errorf("unknown product type")
	}

	if err != nil {
		beego.Error("Failed to update product:", err, " (productType=", productType, ")")
		this.CustomAbort(500, "Failed to update product")
	}

	this.Data["json"] = response
	this.ServeJson()
}

// @Title Archive Product
// @Description Archive product
// @Param	productId	query	int	true	"Product ID"
// @Success 200 string
// @Failure	400	Bad Request
// @Failure	401	Unauthorized
// @Failure 500 Internal Server Error
// @router /:productId/archive [put]
func (this *ProductsController) ArchiveProduct() {
	if !this.IsAdmin() {
		beego.Error("Unauthorized attempt to archvie product")
		this.CustomAbort(401, "Unauthorized")
	}

	productId, err := this.GetInt64(":productId")
	if err != nil {
		beego.Error("Failed to get :productId variable")
		this.CustomAbort(400, "Incorrect productId")
	}

	var product *products.Product
	product, err = products.GetProduct(productId)
	if err != nil {
		beego.Error("Failed to get product")
		this.CustomAbort(500, "Failed to get product")
	}

	err = products.ArchiveProduct(product)
	if err != nil {
		beego.Error("Failed to archive product")
		this.CustomAbort(500, "Failed to archive product")
	}

	this.ServeJson()
}
