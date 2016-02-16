package modelTest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestNetswitch(t *testing.T) {
	Convey("Testing Netswitch model", t, func() {
		Reset(setup.ResetDB)
		Convey("Testing CreateNetswitchMapping", func() {
			Convey("Creating a netswitch mapping regulary", func() {
				nid, err := models.CreateNetSwitchMapping(0)

				So(err, ShouldBeNil)
				So(nid, ShouldBeGreaterThan, 0)
			})
			Convey("Creating two netswitch mapping on the same machine", func() {
				nid1, err1 := models.CreateNetSwitchMapping(0)
				nid2, err2 := models.CreateNetSwitchMapping(0)

				So(err1, ShouldBeNil)
				So(nid1, ShouldBeGreaterThan, 0)
				So(err2, ShouldBeNil)
				So(nid2, ShouldBeGreaterThan, 0)
			})
		})
		Convey("Testing GetNetSwitchMapping", func() {
			Convey("Getting netswith on non-existing machine", func() {
				_, err := models.GetNetSwitchMapping(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Creating netswith and getting it ", func() {
				nid, _ := models.CreateNetSwitchMapping(0)
				netswitch, err := models.GetNetSwitchMapping(0)

				So(err, ShouldBeNil)
				So(netswitch.Id, ShouldEqual, nid)
			})
		})
		Convey("Testing DeleteNetSwitchMapping", func() {
			Convey("Creating a netswitch and deleting it", func() {
				models.CreateNetSwitchMapping(0)
				err := models.DeleteNetSwitchMapping(0)

				So(err, ShouldBeNil)
			})
			Convey("Deleting a non-existing netswitch mapping", func() {
				err := models.DeleteNetSwitchMapping(0)

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing UpdateNetSwitchMapping", func() {
			urlOffTest := "QQ"
			SkipConvey("Creating a netswitch and updating it", func() {
				nid, _ := models.CreateNetSwitchMapping(0)
				netswitch, _ := models.GetNetSwitchMapping(nid)
				netswitch.UrlOff = urlOffTest
				err := netswitch.Update()
				netswitch, _ = models.GetNetSwitchMapping(nid)

				So(err, ShouldBeNil)
				So(netswitch.UrlOff, ShouldEqual, urlOffTest)
			})
		})
		Convey("Netswitch url on/off 200 with should give no error", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			ns := models.NetSwitchMapping{
				UrlOn:  ts.URL + "?method=on",
				UrlOff: ts.URL + "?method=off",
				Xmpp:   false,
			}
			So(ns.On(), ShouldBeNil)
			So(ns.Off(), ShouldBeNil)
		})
		Convey("Netswitch url on/off 500 with should give an error", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer ts.Close()

			ns := models.NetSwitchMapping{
				UrlOn:  ts.URL + "?method=on",
				UrlOff: ts.URL + "?method=off",
				Xmpp:   false,
			}
			So(ns.On(), ShouldNotBeNil)
			So(ns.Off(), ShouldNotBeNil)
		})
	})
}
