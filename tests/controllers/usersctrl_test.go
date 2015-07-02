package controllerTest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
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
		Convey("Testing POST /users/login/", func() {
			Convey("Try to log in without parameters, should return 401", func() {
				r, _ := http.NewRequest("POST", "/api/users/login", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to log in with wrong parameters, should return 401", func() {
				r, _ := http.NewRequest("POST", "/api/users/login?username=a&password=a", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to login with good parameters, should return 200", func() {
				u := models.User{
					Username: "aze",
					Email:    "aze@easylab.io",
				}
				uid, _ := models.CreateUser(&u)
				models.AuthSetPassword(uid, "aze")

				r, _ := http.NewRequest("POST", "/api/users/login?username=aze&password=aze", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
		})
		Convey("Testing POST /users/loginuid/", func() {
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
					Email:    "aze@easylab.io",
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
		Convey("Testing GET /users/logout", func() {
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
		Convey("Testing GET /users/", func() {
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
		Convey("Testing POST /users/signup/", func() {
			Convey("Try signup with empty body, should return 500", func() {
				var jsonStr = []byte("{}")
				r, _ := http.NewRequest("POST", "/api/users/signup", bytes.NewBuffer(jsonStr))
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 500)
			})
			Convey("Try signup with User only, should return 500", func() {
				var jsonStr = []byte(`{"User": {"Username":"A"} }`)
				r, _ := http.NewRequest("POST", "/api/users/signup", bytes.NewBuffer(jsonStr))
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 500)
			})
			Convey("Try signup with User and password, should return 200", func() {
				var jsonStr = []byte(`{"User": {"Username":"A", "Email": "a@easylab.io"}, "Password":"A" }`)
				r, _ := http.NewRequest("POST", "/api/users/signup", bytes.NewBuffer(jsonStr))
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
		})
		Convey("Testing POST /users/", func() {
			Convey("Try creating user without being logged in, should return 401", func() {
				r, _ := http.NewRequest("POST", "/api/users/", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try creating user without parameters, should return 500", func() {
				r, _ := http.NewRequest("POST", "/api/users/", nil)
				r.AddCookie(LoginAsAdmin())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 500)
			})
			Convey("Try creating user with email, should return 200", func() {
				r, _ := http.NewRequest("POST", "/api/users/?email=a@easylab.io", nil)
				r.AddCookie(LoginAsAdmin())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
			Convey("Try creating user as a regular user, should return 401", func() {
				r, _ := http.NewRequest("POST", "/api/users/?email=a@easylab.io", nil)
				r.AddCookie(LoginAsRegular())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
		})
		Convey("Testing GET /users/:uid", func() {
			Convey("Try to get user without being logged in, should return 401", func() {
				r, _ := http.NewRequest("GET", "/api/users/1", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to get non-existing user, should return 403", func() {
				r, _ := http.NewRequest("GET", "/api/users/0", nil)
				r.AddCookie(LoginAsAdmin())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 403)
			})
			Convey("Try to get existing user, should return 200", func() {
				u := &models.User{
					Username: "A",
					Email:    "a@easylab.io",
				}
				mid, _ := models.CreateUser(u)
				r, _ := http.NewRequest("GET", "/api/users/"+strconv.FormatInt(mid, 10), nil)
				r.AddCookie(LoginAsAdmin())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
			Convey("Try to get existing user as a regular one, should return 401", func() {
				u := &models.User{
					Username: "A",
					Email:    "a@easylab.io",
				}
				mid, _ := models.CreateUser(u)
				r, _ := http.NewRequest("GET", "/api/users/"+strconv.FormatInt(mid, 10), nil)
				r.AddCookie(LoginAsRegular())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
		})
		Convey("Testing DELETE /users/:uid", func() {
			Convey("Try to delete user without being logged in, should return 401", func() {
				r, _ := http.NewRequest("DELETE", "/api/users/0", nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to delete user as a regular user, should return 401", func() {
				r, _ := http.NewRequest("DELETE", "/api/users/0", nil)
				r.AddCookie(LoginAsRegular())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to delete a non-existing user, should return 403", func() {
				r, _ := http.NewRequest("DELETE", "/api/users/0", nil)
				r.AddCookie(LoginAsAdmin())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 403)
			})
			Convey("Create a user and delete it, should return 200", func() {
				u := &models.User{
					Username: "A",
					Email:    "a@easylab.io",
				}
				mid, _ := models.CreateUser(u)
				r, _ := http.NewRequest("DELETE", "/api/users/"+strconv.FormatInt(mid, 10), nil)
				r.AddCookie(LoginAsAdmin())
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
		})
		Convey("Testing PUT /users/:uid", func() {
			Convey("Try to modify a user without being connected, should return 500", func() {
				var jsonStr = []byte("{}")
				r, _ := http.NewRequest("PUT", "/api/users/0", bytes.NewBuffer(jsonStr))
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 500)
			})
			Convey("Try to modify a user as a regular user, should return 401", func() {
				var jsonStr = []byte("{}")
				r, _ := http.NewRequest("PUT", "/api/users/0", bytes.NewBuffer(jsonStr))
				r.AddCookie(LoginAsRegular())
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 401)
			})
			Convey("Try to modify self userRole as a regular user, should return 500", func() {
				cookie := LoginAsRegular()
				var jsonStr = []byte(`{"User": {"Id": ` + strconv.FormatInt(RegularUID, 10) + `, "Email": "raaaaaaaaaaaaaaaaadom@easylab.io", "UserRole": "` + models.ADMIN + `"}}`)
				r, _ := http.NewRequest("PUT", "/api/users/"+strconv.FormatInt(RegularUID, 10), bytes.NewBuffer(jsonStr))
				r.AddCookie(cookie)
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 500)
			})
			Convey("Try to modify self user as a regular user, should return 200", func() {
				cookie := LoginAsRegular()
				var jsonStr = []byte(`{"User": {"Id": ` + strconv.FormatInt(RegularUID, 10) + `, "Email": "raaaaaaaaaaaaaaaaadom@easylab.io"}}`)
				r, _ := http.NewRequest("PUT", "/api/users/"+strconv.FormatInt(RegularUID, 10), bytes.NewBuffer(jsonStr))
				r.AddCookie(cookie)
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
			Convey("Try to modify a user that doesn't exists as an admin, should return 500", func() {
				var jsonStr = []byte(`{"User": {"Id": 0, "Email": "raaaaaaaaaaaaaaaaadom@easylab.io"}}`)
				r, _ := http.NewRequest("PUT", "/api/users/0", bytes.NewBuffer(jsonStr))
				r.AddCookie(LoginAsAdmin())
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 500)
			})
			Convey("Try to modify a user as an admin, should return 200", func() {
				u := models.User{
					Username: "lel",
					Email:    "lel@easylab.io",
				}
				uid, _ := models.CreateUser(&u)
				var jsonStr = []byte(`{"User": {"Id": ` + strconv.FormatInt(uid, 10) + `, "Email": "raaaaaaaaaaaaaaaaadom@easylab.io"}}`)
				r, _ := http.NewRequest("PUT", "/api/users/"+strconv.FormatInt(uid, 10), bytes.NewBuffer(jsonStr))
				r.AddCookie(LoginAsAdmin())
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				So(w.Code, ShouldEqual, 200)
			})
		})
	})
}
