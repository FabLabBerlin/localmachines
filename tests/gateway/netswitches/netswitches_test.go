package gatewayNetswitchesTest

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitch"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/tests/gateway/mocks"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNetswitches(t *testing.T) {
	netSwitch1 := mocks.NewNetSwitch(mocks.DESIRED_OFF, mocks.RELAY_OFF)
	defer netSwitch1.Close()
	netSwitch2 := mocks.NewNetSwitch(mocks.DESIRED_OFF, mocks.RELAY_OFF)
	defer netSwitch2.Close()
	netSwitch3 := mocks.NewNetSwitch(mocks.DESIRED_OFF, mocks.RELAY_OFF)
	defer netSwitch3.Close()

	lmApi := mocks.NewLmApi()
	defer lmApi.Close()
	lmApi.AddMapping(machine.Machine{
		Id:                  11,
		NetswitchHost:       netSwitch1.Host(),
		NetswitchSensorPort: 1,
		NetswitchXmpp:       true,
	})
	lmApi.AddMapping(machine.Machine{
		Id:            22,
		NetswitchHost: netSwitch2.Host(),
		NetswitchXmpp: false,
	})

	global.Cfg.API.Url = lmApi.URL()
	global.Cfg.Main.StateFile = "foo.state.test"

	netSwitches := netswitches.New()

	Convey("Testing Load", t, func() {
		client := &http.Client{}
		err := netSwitches.Load(client)
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
			So(ns.Id, ShouldEqual, 11)
			So(ns.NetswitchHost, ShouldEqual, netSwitch1.Host())
			So(ns.NetswitchSensorPort, ShouldEqual, 1)
			So(ns.NetswitchXmpp, ShouldBeTrue)
		})

		Convey("It fails when there are duplicate combinations hosts+sensor port that could lead to contradictory state which in turn could lead to switches turning on and off every 30 seconds", func() {
			lmApi.AddMapping(machine.Machine{
				Id:                  44,
				NetswitchHost:       netSwitch1.Host(),
				NetswitchSensorPort: 1,
				NetswitchXmpp:       true,
			})
			err := netSwitches.Load(client)
			So(err, ShouldNotBeNil)
			So(strings.Contains(err.Error(), netswitches.ErrDuplicateCombinationHostSensorPort.Error()), ShouldBeTrue)
			lmApi.DeleteMapping(44)
			err = netSwitches.Load(client)
			So(err, ShouldBeNil)
		})
	})

	Convey("Testing Sync", t, func() {
		Convey("It should poll the Xmpp switch", func() {
			before := netSwitch1.PollRequests
			err := netSwitches.Sync(11)
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			after := netSwitch1.PollRequests
			So(after-before, ShouldEqual, 1)
		})
	})

	Convey("Testing SyncAll", t, func() {
		Convey("It should poll the only registered Xmpp switch", func() {
			before := netSwitch1.PollRequests
			err := netSwitches.SyncAll()
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			after := netSwitch1.PollRequests
			So(after-before, ShouldEqual, 1)
		})

		Convey("It should poll the two registered Xmpp switches", func() {
			lmApi.AddMapping(machine.Machine{
				Id:                  33,
				NetswitchHost:       netSwitch3.Host(),
				NetswitchSensorPort: 1,
				NetswitchXmpp:       true,
			})
			before1 := netSwitch1.PollRequests
			client := &http.Client{}
			err := netSwitches.Load(client)
			So(err, ShouldBeNil)
			err = netSwitches.SyncAll()
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			after1 := netSwitch1.PollRequests
			So(after1-before1, ShouldEqual, 1)
			So(netSwitch3.PollRequests, ShouldEqual, 1)
		})

		Convey("It should only poll the one remaining Xmpp switch", func() {
			lmApi.DeleteMapping(33)
			before1 := netSwitch1.PollRequests
			before3 := netSwitch3.PollRequests
			client := &http.Client{}
			err := netSwitches.Load(client)
			So(err, ShouldBeNil)
			err = netSwitches.SyncAll()
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			after1 := netSwitch1.PollRequests
			after3 := netSwitch3.PollRequests
			So(after1-before1, ShouldEqual, 1)
			So(after3-before3, ShouldEqual, 0)
			So(netSwitch3.SwitchRequests, ShouldEqual, 0)
		})

		Convey("It discards a switch when it's changed to non-xmpp", func() {
			ns := mocks.NewNetSwitch(mocks.DESIRED_OFF, mocks.RELAY_OFF)
			defer ns.Close()
			lmApi.AddMapping(machine.Machine{
				Id:                  55,
				NetswitchHost:       ns.Host(),
				NetswitchSensorPort: 123,
				NetswitchXmpp:       true,
			})
			client := &http.Client{}
			err := netSwitches.Load(client)
			So(err, ShouldBeNil)
			err = netSwitches.SyncAll()
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			So(ns.PollRequests, ShouldEqual, 1)
			lmApi.UpdateMapping(55, machine.Machine{
				Id:                  55,
				NetswitchHost:       ns.Host(),
				NetswitchSensorPort: 123,
				NetswitchXmpp:       false,
			})
			err = netSwitches.Load(client)
			So(err, ShouldBeNil)
			err = netSwitches.SyncAll()
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			So(ns.PollRequests, ShouldEqual, 1)
		})
	})

	// SetOn affects synchronization. So it's best to do this test *after*
	// playing around with Sync methods.
	Convey("Testing SetOn", t, func() {
		pollBefore := netSwitch1.PollRequests
		switchBefore := netSwitch1.SwitchRequests
		netSwitches.SetOn(11, true)
		<-time.After(time.Second)

		Convey("It should trigger one poll request", func() {
			after := netSwitch1.PollRequests
			So(after-pollBefore, ShouldEqual, 1)
		})

		Convey("It should trigger one switch request", func() {
			after := netSwitch1.SwitchRequests
			So(after-switchBefore, ShouldEqual, 1)
		})
	})

	Convey("netswitch 2 should not have received any requests because it's not on XMPP", t, func() {
		So(netSwitch2.PollRequests, ShouldEqual, 0)
		So(netSwitch2.SwitchRequests, ShouldEqual, 0)
	})

	Convey("netswitch 3 should have received only one poll request and no switch requests", t, func() {
		So(netSwitch3.PollRequests, ShouldEqual, 1)
		So(netSwitch3.SwitchRequests, ShouldEqual, 0)
	})
}
