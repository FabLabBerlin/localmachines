package controllers

import (
	"github.com/kr15h/fabsmith/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

// Operations about object
type TestObjectController struct {
	beego.Controller
}

// @Title create
// @Description create object
// @Param	body		body 	models.TestObject	true		"The object content"
// @Success 200 {string} models.TestObject.Id
// @Failure 403 body is empty
// @router / [post]
func (this *TestObjectController) Post() {
	var ob models.TestObject
	json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
	objectid := models.AddOne(ob)
	this.Data["json"] = map[string]string{"ObjectId": objectid}
	this.ServeJson()
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.TestObject
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (this *TestObjectController) Get() {
	objectId := this.Ctx.Input.Params[":objectId"]
	if objectId != "" {
		ob, err := models.GetOne(objectId)
		if err != nil {
			this.Data["json"] = err
		} else {
			this.Data["json"] = ob
		}
	}
	this.ServeJson()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.TestObject
// @Failure 403 :objectId is empty
// @router / [get]
func (this *TestObjectController) GetAll() {
	obs := models.GetAll()
	this.Data["json"] = obs
	this.ServeJson()
}

// @Title update
// @Description update the object
// @Param	objectId		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.TestObject	true		"The body"
// @Success 200 {object} models.TestObject
// @Failure 403 :objectId is empty
// @router /:objectId [put]
func (this *TestObjectController) Put() {
	objectId := this.Ctx.Input.Params[":objectId"]
	var ob models.TestObject
	json.Unmarshal(this.Ctx.Input.RequestBody, &ob)

	err := models.Update(objectId, ob.Score)
	if err != nil {
		this.Data["json"] = err
	} else {
		this.Data["json"] = "update success!"
	}
	this.ServeJson()
}

// @Title delete
// @Description delete the object
// @Param	objectId		path 	string	true		"The objectId you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /:objectId [delete]
func (this *TestObjectController) Delete() {
	objectId := this.Ctx.Input.Params[":objectId"]
	models.Delete(objectId)
	this.Data["json"] = "delete success!"
	this.ServeJson()
}

