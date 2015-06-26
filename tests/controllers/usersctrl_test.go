package controllerTest

import (
	"bytes"
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
		Convey("Testing /users/login/", func() {
			Convey("Try to log in without parameters, should return 403", func() {
				r, _ := http.NewRequest("POST", "/api/users/login", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to log in with wrong parameters, should return 403", func() {
				r, _ := http.NewRequest("POST", "/api/users/login?username=a&password=a", nil)
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
		Convey("Testing /users/loginuid/", func() {
			Convey("Try to log in without uid parameter, should return 403", func() {
				r, _ := http.NewRequest("POST", "/api/users/loginuid", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to log in with wrong parameters, should return 403", func() {
				r, _ := http.NewRequest("POST", "/api/users/loginuid?uid=a", nil)
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
				models.AuthUpdateNfcUid(uid, "123456")

				r, _ := http.NewRequest("POST", "/api/users/loginuid?uid=123456", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
		})
		Convey("Testing /users/logout", func() {
			Convey("Try to logout without being logged in, should return 200", func() {
				r, _ := http.NewRequest("GET", "/api/users/logout", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
			Convey("Try to logout after being logged in as a regular user, should return 200", func() {
				r, _ := http.NewRequest("GET", "/api/users/logout", nil)
				r.AddCookie(LoginAsRegular())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
			Convey("Try to logout after being logged in as an admin, should return 200", func() {
				r, _ := http.NewRequest("GET", "/api/users/logout", nil)
				r.AddCookie(LoginAsAdmin())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
		})
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
		Convey("Testing /users/signup/", func() {
			Convey("Try signup with empty body", func() {
				var jsonStr = []byte("")
				r, _ := http.NewRequest("POST", "/api/users/signup", bytes.NewBuffer(jsonStr))
				r.Header.Set("X-Custom-Header", "myvalue")
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 500)
			})
			Convey("Try signup with User only", func() {
				var jsonStr = []byte(`{"User": {"Username":"A"} }`)
				r, _ := http.NewRequest("POST", "/api/users/signup", bytes.NewBuffer(jsonStr))
				r.Header.Set("X-Custom-Header", "myvalue")
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
			Convey("Try signup with User and password", func() {
				var jsonStr = []byte(`{"User": {"Username":"A"}, "Password":"A" }`)
				r, _ := http.NewRequest("POST", "/api/users/signup", bytes.NewBuffer(jsonStr))
				r.Header.Set("X-Custom-Header", "myvalue")
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
		})
	})
}
