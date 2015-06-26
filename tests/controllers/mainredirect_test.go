package controllerTest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/kr15h/fabsmith/routers"

	. "github.com/kr15h/fabsmith/tests/models"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

// TestMain is a sample to run an endpoint test
func TestRedirect(t *testing.T) {
	SkipConvey("Test Station Endpoint", t, func() {
		Convey("Status Code Should Be 302", func() {
			r, _ := http.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			So(w.Code, ShouldEqual, 302)
		})
	})
}
