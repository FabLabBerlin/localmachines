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
	// Create response
	var response struct{ Users []PublicUser }
	var err error
	if this.isAdmin() || this.isStaff() {
		response.Users, err = this.getAllUsers()
		if err != nil {
			if beego.AppConfig.String("runmode") == "dev" {
				panic("Could not get all users")
			}
			this.serveErrorResponse("Could not get users")
		}
	} else {
		response.Users, err = this.getSessionUser()
		if err != nil {
			if beego.AppConfig.String("runmode") == "dev" {
				panic("Could not get session user")
			}
			this.serveErrorResponse("Could not get session user")
		}
	}
	this.Data["json"] = &response
	this.ServeJson()
}

// If user role is admin or staff, return all users
func (this *UsersController) getAllUsers() ([]PublicUser, error) {
	beego.Trace("Attempt to get all users from DB")
	users := []models.User{}
	o := orm.NewOrm()
	num, err := o.Raw("Select * FROM user").QueryRows(&users)
	if err != nil {
		beego.Error("Could not get all users:", err)
		return nil, err
	} else {
		beego.Trace("Got", num, "users")
	}
	// Loop through users and cherry pick values
	arrPubUsers := []PublicUser{}
	for i := range users {
		userFullName := fmt.Sprintf("%s %s", users[i].FirstName, users[i].LastName)
		user := PublicUser{users[i].Id, userFullName, users[i].Email}
		arrPubUsers = append(arrPubUsers, user)
	}
	return arrPubUsers, nil
}

// If user role is NOT admin or staff, return session user
func (this *UsersController) getSessionUser() ([]PublicUser, error) {
	userModelInterface, err := this.getUser()
	if err != nil {
		return nil, err
	}
	userModel := userModelInterface.(models.User)
	// Interpret it as PublicUser
	pubUser := PublicUser{}
	pubUser.Id = userModel.Id
	pubUser.Name = fmt.Sprintf("%s %s", userModel.FirstName, userModel.LastName)
	pubUser.Email = userModel.Email
	arr := []PublicUser{pubUser}
	return arr, nil
}
