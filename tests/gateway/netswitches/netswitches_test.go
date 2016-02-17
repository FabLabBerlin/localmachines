package gatewayNetswitchesTest

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitch"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	modelsNetswitch "github.com/FabLabBerlin/localmachines/models/netswitch"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNetswitches(t *testing.T) {
	pollRequests := 0
	switchRequests := 0

	netSwitch := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			pollRequests++
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"sensors":[{"output":0,"power":0.0,"energy":0.0,"enabled":0,"current":0.0,"voltage":0.0,"powerfactor":0.0,"relay":0,"lock":0}],"status":"success"}`)
		} else {
			switchRequests++
		}
	}))
	defer netSwitch.Close()

	url, err := url.Parse(netSwitch.URL)
	if err != nil {
		panic(err.Error())
	}

	easylabApi := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := []modelsNetswitch.Mapping{
			modelsNetswitch.Mapping{
				Id:         1,
				MachineId:  11,
				Host:       url.Host,
				SensorPort: 1,
				Xmpp:       true,
			},
			modelsNetswitch.Mapping{
				Id:        2,
				MachineId: 22,
				Host:      url.Host,
				Xmpp:      false,
			},
		}
		enc := json.NewEncoder(w)
		w.WriteHeader(http.StatusOK)
		enc.Encode(state)
	}))
	defer easylabApi.Close()

	global.Cfg.API.Url = easylabApi.URL
	global.Cfg.Main.StateFile = "foo.state.test"

	netSwitches := netswitches.New()

	Convey("Testing Load", t, func() {
		client := &http.Client{}
		err = netSwitches.Load(client)
		So(err, ShouldBeNil)
		Convey("It should load the Xmpp switches and discard the others", func() {
			netSwitches.Save()
			var nss []netswitch.NetSwitch
			f, err := os.Open(global.Cfg.Main.StateFile)
			if err != nil {
				panic(err.Error())
			}
			defer f.Close()
			dec := json.NewDecoder(f)
			dec.Decode(&nss)
			So(len(nss), ShouldEqual, 1)
			ns := nss[0]
			So(ns.Id, ShouldEqual, 1)
			So(ns.MachineId, ShouldEqual, 11)
			So(ns.Host, ShouldEqual, url.Host)
			So(ns.SensorPort, ShouldEqual, 1)
			So(ns.Xmpp, ShouldBeTrue)
		})
	})

	Convey("Testing Sync", t, func() {
		Convey("It should poll the Xmpp switch", func() {
			before := pollRequests
			err := netSwitches.Sync(11)
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			after := pollRequests
			So(after-before, ShouldEqual, 1)
		})
	})

	Convey("Testing SyncAll", t, func() {
		Convey("It should poll the Xmpp switches", func() {
			before := pollRequests
			err := netSwitches.SyncAll()
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			after := pollRequests
			So(after-before, ShouldEqual, 1)
		})
	})

	// SetOn affects synchronization. So it's best to do this test *after*
	// playing around with Sync methods.
	Convey("Testing SetOn", t, func() {
		pollBefore := pollRequests
		switchBefore := switchRequests
		netSwitches.SetOn(11, true)
		<-time.After(time.Second)

		Convey("It should trigger one poll request", func() {
			after := pollRequests
			So(after-pollBefore, ShouldEqual, 1)
		})

		Convey("It should trigger one switch request", func() {
			after := switchRequests
			So(after-switchBefore, ShouldEqual, 1)
		})
	})
}
