package gatewayNetswitchesTest

import (
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitch"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	modelsNetswitch "github.com/FabLabBerlin/localmachines/models/netswitch"
	"github.com/FabLabBerlin/localmachines/tests/gateway/mocks"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNetswitches(t *testing.T) {
	netSwitch := mocks.NewNetSwitch(mocks.DESIRED_OFF, mocks.RELAY_OFF)
	defer netSwitch.Close()

	lmApi := mocks.NewLmApi()
	defer lmApi.Close()
	lmApi.AddMapping(modelsNetswitch.Mapping{
		Id:         1,
		MachineId:  11,
		Host:       netSwitch.Host(),
		SensorPort: 1,
		Xmpp:       true,
	})
	lmApi.AddMapping(modelsNetswitch.Mapping{
		Id:        2,
		MachineId: 22,
		Host:      netSwitch.Host(),
		Xmpp:      false,
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
			So(ns.Id, ShouldEqual, 1)
			So(ns.MachineId, ShouldEqual, 11)
			So(ns.Host, ShouldEqual, netSwitch.Host())
			So(ns.SensorPort, ShouldEqual, 1)
			So(ns.Xmpp, ShouldBeTrue)
		})
	})

	Convey("Testing Sync", t, func() {
		Convey("It should poll the Xmpp switch", func() {
			before := netSwitch.PollRequests
			err := netSwitches.Sync(11)
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			after := netSwitch.PollRequests
			So(after-before, ShouldEqual, 1)
		})
	})

	Convey("Testing SyncAll", t, func() {
		Convey("It should poll the Xmpp switches", func() {
			before := netSwitch.PollRequests
			err := netSwitches.SyncAll()
			So(err, ShouldBeNil)
			<-time.After(time.Second)
			after := netSwitch.PollRequests
			So(after-before, ShouldEqual, 1)
		})
	})

	// SetOn affects synchronization. So it's best to do this test *after*
	// playing around with Sync methods.
	Convey("Testing SetOn", t, func() {
		pollBefore := netSwitch.PollRequests
		switchBefore := netSwitch.SwitchRequests
		netSwitches.SetOn(11, true)
		<-time.After(time.Second)

		Convey("It should trigger one poll request", func() {
			after := netSwitch.PollRequests
			So(after-pollBefore, ShouldEqual, 1)
		})

		Convey("It should trigger one switch request", func() {
			after := netSwitch.SwitchRequests
			So(after-switchBefore, ShouldEqual, 1)
		})
	})
}
