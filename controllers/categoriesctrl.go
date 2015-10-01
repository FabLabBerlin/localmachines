package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type CategoriesController struct {
	Controller
}

// @router / [get]
func (this *CategoriesController) GetAll() {

	// This is admin and staff only
	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var categories []models.Category
	var err error
	categories, err = models.GetAllCategories()
	if err != nil {
		beego.Error("Failed to get all categories", err)
		this.CustomAbort(403, "Failed to get all categories")
	}

	this.Data["json"] = categories
	this.ServeJson()
}

// @router / [post]
func (this *CategoriesController) Create() {
	name := this.GetString("cname")
	beego.Trace(name)

	if !this.IsAdmin() && !this.IsStaff() {
		beego.Error("Not authorized to create category")
		this.CustomAbort(401, "Not authorized")
	}

	// All clear - create machine in the database
	var categoryId int64
	var err error
	categoryId, err = models.CreateCategory(name)
	if err != nil {
		beego.Error("Failed to create category", err)
		this.CustomAbort(403, "Failed to create category")
	}

	// Success - we even got an ID!!!
	this.Data["json"] = map[string]int64{
		"id": categoryId,
	}
	this.ServeJson()
}

// @router /:cid [put]
func (this *CategoriesController) Update() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error

	// Attempt to decode passed json
	dec := json.NewDecoder(this.Ctx.Request.Body)
	req := models.Category{}
	if err = dec.Decode(&req); err != nil {
		beego.Error("Failed to decode json:", err)
		this.CustomAbort(403, "Failed to update category")
	}

	var cid int64

	// Get cid and check if it matches with the machine model ID
	cid, err = this.GetInt64(":cid")
	if err != nil {
		beego.Error("Could not get :cid:", err)
		this.CustomAbort(403, "Failed to update category")
	}
	if cid != req.Id {
		beego.Error("cid and model ID does not match:", err)
		this.CustomAbort(403, "Failed to update category")
	}

	// Update the database
	err = models.UpdateCategory(&req)
	if err != nil {
		beego.Error("Failed updating category:", err)
		this.CustomAbort(403, "Failed to update category")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}

// @router /:cid [delete]
func (this *CategoriesController) Delete() {

	if !this.IsAdmin() {
		beego.Error("Not authorized")
		this.CustomAbort(401, "Not authorized")
	}

	var err error
	var cid int64

	cid, err = this.GetInt64(":cid")
	if err != nil {
		beego.Error("Failed to get cid:", err)
		this.CustomAbort(403, "Failed to delete category")
	}

	err = models.DeleteCategory(cid)
	if err != nil {
		beego.Error("Failed to delete category:", err)
		this.CustomAbort(403, "Failed to delete category")
	}

	this.Data["json"] = "ok"
	this.ServeJson()
}



