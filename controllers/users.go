package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

// Inherits from fabsmith root controler that contains ton of cool features
type UsersController struct {
	Controller
}

// Contains only the data that should be passed
// to JSON output of this controller
type PublicUser struct {
	Id    int
	Name  string
	Email string
}

// Output JSON user list for API /users request
func (this *UsersController) GetUsers() {
	var isAdmin bool = false
	userRoles := this.getSessionUserRoles()
	if userRoles.Admin || userRoles.Staff {
		isAdmin = true
		beego.Info("Admin or staff user detected")
	} else {
		beego.Info("Regular user detected")
	}
	// Create response
	var response struct{ Users []PublicUser }
	if isAdmin {
		response.Users = this.getAllUsers()
	} else {
		response.Users = this.getSessionUser()
	}
	this.Data["json"] = &response
	this.ServeJson()
}

// If user role is admin or staff, return all users
func (this *UsersController) getAllUsers() []PublicUser {
	o := orm.NewOrm()
	var users []models.User
	num, err := o.Raw("Select * FROM user").QueryRows(&users)
	if err != nil {
		beego.Error(err)
	} else {
		beego.Info("Got", num, "users")
	}
	// Loop through users and cherry pick values
	var arrPubUsers []PublicUser
	for i := range users {
		userFullName := fmt.Sprintf("%s %s", users[i].FirstName, users[i].LastName)
		user := PublicUser{users[i].Id, userFullName, users[i].Email}
		arrPubUsers = append(arrPubUsers, user)
	}
	return arrPubUsers
}

// If user role is NOT admin or staff, return session user
func (this *UsersController) getSessionUser() []PublicUser {
	var userModel *models.User
	userModel = this.getSessionUserData()
	// Interpret it as PublicUser
	var pubUser PublicUser
	pubUser.Id = userModel.Id
	pubUser.Name = fmt.Sprintf("%s %s", userModel.FirstName, userModel.LastName)
	pubUser.Email = userModel.Email
	arr := []PublicUser{pubUser}
	return arr
}
