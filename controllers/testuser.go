package controllers

import (
	"github.com/kr15h/fabsmith/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

// Operations about Users
type TestUserController struct {
	beego.Controller
}

// @Title createUser
// @Description create users
// @Param	body		body 	models.TestUser	true		"body for user content"
// @Success 200 {int} models.TestUser.Id
// @Failure 403 body is empty
// @router / [post]
func (u *TestUserController) Post() {
	var user models.TestUser
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	uid := models.AddUser(user)
	u.Data["json"] = map[string]string{"uid": uid}
	u.ServeJson()
}

/*
// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.TestUser
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *TestUserController) Get() {
	uid := u.GetString(":uid")
	if uid != "" {
		user, err := models.GetUser(uid)
		if err != nil {
			u.Data["json"] = err
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJson()
}
*/

// @Title update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.TestUser	true		"body for user content"
// @Success 200 {object} models.TestUser
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *TestUserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.TestUser
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(uid, &user)
		if err != nil {
			u.Data["json"] = err
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJson()
}

// @Title delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *TestUserController) Delete() {
	uid := u.GetString(":uid")
	models.DeleteUser(uid)
	u.Data["json"] = "delete success!"
	u.ServeJson()
}

