package controllerTest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/kr15h/fabsmith/routers"

	"github.com/kr15h/fabsmith/models"
	. "github.com/kr15h/fabsmith/tests/models"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

// TestMain is a sample to run an endpoint test
func TestUsersAPI(t *testing.T) {

	Convey("Test users API", t, func() {
		Reset(ResetDB)
		Convey("Testing /users/", func() {
			Convey("Try to get users without being logged in, should return 401", func() {
				r, _ := http.NewRequest("GET", "/api/users/", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to get users being logged in as a regular user, should return 401", func() {
				r, _ := http.NewRequest("GET", "/api/users/", nil)
				r.AddCookie(LoginAsRegular())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to get users being logged in as an admin, should return 200", func() {
				r, _ := http.NewRequest("GET", "/api/users/", nil)
				r.AddCookie(LoginAsAdmin())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
		})
		Convey("Testing /users/login/", func() {
			Convey("Try to log in without parameters, should return 403", func() {
				r, _ := http.NewRequest("POST", "/api/users/login/", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to log in with wrong parameters, should return 403", func() {
				r, _ := http.NewRequest("POST", "/api/users/login/?username=a&password=a", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to login with good parameters, should return 200", func() {
				u := models.User{
					Username: "aze",
				}
				uid, _ := models.CreateUser(&u)
				models.AuthSetPassword(uid, "aze")

				r, _ := http.NewRequest("POST", "/api/users/login?username=aze&password=aze", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
		})
	})
}
