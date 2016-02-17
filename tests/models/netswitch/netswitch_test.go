package netswitchTest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FabLabBerlin/localmachines/models/netswitch"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestNetswitch(t *testing.T) {
	Convey("Testing Netswitch model", t, func() {
		Reset(setup.ResetDB)
		Convey("Testing CreateMapping", func() {
			Convey("Creating a netswitch mapping regulary", func() {
				nid, err := netswitch.CreateMapping(0)

				So(err, ShouldBeNil)
				So(nid, ShouldBeGreaterThan, 0)
			})
			Convey("Creating two netswitch mapping on the same machine", func() {
				nid1, err1 := netswitch.CreateMapping(0)
				nid2, err2 := netswitch.CreateMapping(0)

				So(err1, ShouldBeNil)
				So(nid1, ShouldBeGreaterThan, 0)
				So(err2, ShouldBeNil)
				So(nid2, ShouldBeGreaterThan, 0)
			})
			Convey("Creating same netswitch mapping for two machines", func() {
				_, err := netswitch.CreateMapping(1)
				if err != nil {
					panic(err.Error())
				}
				_, err = netswitch.CreateMapping(2)
				if err != nil {
					panic(err.Error())
				}
				n1, err := netswitch.GetMapping(1)
				if err != nil {
					panic(err.Error())
				}
				n2, err := netswitch.GetMapping(2)
				if err != nil {
					panic(err.Error())
				}
				n1.Host = "example.com"
				err = n1.Update()
				So(err, ShouldBeNil)
				n2.Host = "example.com"
				err = n2.Update()
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing GetMapping", func() {
			Convey("Getting netswith on non-existing machine", func() {
				_, err := netswitch.GetMapping(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Creating netswith and getting it ", func() {
				nid, _ := netswitch.CreateMapping(0)
				netswitch, err := netswitch.GetMapping(0)

				So(err, ShouldBeNil)
				So(netswitch.Id, ShouldEqual, nid)
			})
		})
		Convey("Testing DeleteMapping", func() {
			Convey("Creating a netswitch and deleting it", func() {
				netswitch.CreateMapping(0)
				err := netswitch.DeleteMapping(0)

				So(err, ShouldBeNil)
			})
			Convey("Deleting a non-existing netswitch mapping", func() {
				err := netswitch.DeleteMapping(0)

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing UpdateMapping", func() {
			urlOffTest := "QQ"
			SkipConvey("Creating a netswitch and updating it", func() {
				nid, _ := netswitch.CreateMapping(0)
				mapping, _ := netswitch.GetMapping(nid)
				mapping.UrlOff = urlOffTest
				err := mapping.Update()
				mapping, _ = netswitch.GetMapping(nid)

				So(err, ShouldBeNil)
				So(mapping.UrlOff, ShouldEqual, urlOffTest)
			})
		})
		Convey("Netswitch url on/off 200 with should give no error", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			ns := netswitch.Mapping{
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

			ns := netswitch.Mapping{
				UrlOn:  ts.URL + "?method=on",
				UrlOff: ts.URL + "?method=off",
				Xmpp:   false,
			}
			So(ns.On(), ShouldNotBeNil)
			So(ns.Off(), ShouldNotBeNil)
		})
	})
}
