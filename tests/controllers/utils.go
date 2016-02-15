package controllerTest

import (
	"net/http"
	"net/http/httptest"

	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/astaxie/beego"
)

var AdminUID int64
var RegularUID int64

// LoginAsAdmin : Create an admin user and login
func LoginAsAdmin() *http.Cookie {
	u := models.User{
		Username: "admin",
		Email:    "admin@easylab.io",
		UserRole: user_roles.ADMIN.String(),
	}
	uid, err := models.CreateUser(&u)
	if err != nil {
		panic(err.Error())
	}
	AdminUID = uid
	models.AuthSetPassword(uid, "admin")
	ul := user_locations.UserLocation{
		LocationId: 1,
		UserId:     uid,
		UserRole:   user_roles.ADMIN.String(),
	}
	if _, err := user_locations.Create(&ul); err != nil {
		panic(err.Error())
	}

	r, _ := http.NewRequest("POST", "/api/users/login?username=admin&password=admin&location=1", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	cookie, _ := r.Cookie("localmachines")

	return cookie
}

// LoginAsRegular : Create an admin user and login
func LoginAsRegular() *http.Cookie {
	u := models.User{
		Username: "user",
		Email:    "user@easylab.io",
	}
	uid, _ := models.CreateUser(&u)
	RegularUID = uid
	models.AuthSetPassword(uid, "user")

	r, _ := http.NewRequest("POST", "/api/users/login?username=user&password=user&location=1", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	cookie, _ := r.Cookie("localmachines")

	return cookie
}
