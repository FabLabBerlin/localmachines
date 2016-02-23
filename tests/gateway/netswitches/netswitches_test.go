package gatewayNetswitchesTest

import (
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/tests/gateway/mocks"
	"net/http"
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

	netSwitches := netswitches.New()

	Convey("Testing Load", t, func() {
		client := &http.Client{}
		err := netSwitches.Load(client)
		So(err, ShouldBeNil)
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

	// SetOn affects synchronization. So it's best to do this test *after*
	// playing around with Sync methods.
	Convey("Testing SetOn", t, func() {
		switchBefore := netSwitch1.SwitchRequests
		netSwitches.SetOn(11, true)
		<-time.After(time.Second)

		Convey("It should trigger one switch request", func() {
			after := netSwitch1.SwitchRequests
			So(after-switchBefore, ShouldEqual, 1)
		})
	})

	Convey("netswitch 2 should not have received any requests because it's not on XMPP", t, func() {
		So(netSwitch2.SwitchRequests, ShouldEqual, 0)
	})

	Convey("netswitch 3 should have received only one poll request and no switch requests", t, func() {
		So(netSwitch3.SwitchRequests, ShouldEqual, 0)
	})
}
