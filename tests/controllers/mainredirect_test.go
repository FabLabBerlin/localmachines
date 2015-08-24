package controllerTest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

// TestMain is a sample to run an endpoint test
func TestRedirect(t *testing.T) {
	Convey("Test Station Endpoint", t, func() {
		Convey("Status Code Should Be 302", func() {
			r, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				panic(err.Error())
			}
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, 302)
		})
	})
}
