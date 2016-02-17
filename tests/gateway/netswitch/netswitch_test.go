package gatewayNetswitchTest

import (
	"github.com/FabLabBerlin/localmachines/tests/gateway/mocks"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNetswitch(t *testing.T) {
	Convey("Testing Sync", t, func() {
		Convey("A switch that is off and supposed to be off stays off", func() {
			mock := mocks.NewNetSwitch(mocks.DESIRED_OFF, mocks.RELAY_OFF)
			So(mock.NetSwitch.Sync(), ShouldBeNil)
			<-time.After(time.Second)
			So(mock.PollRequests, ShouldEqual, 1)
			So(mock.SwitchRequests, ShouldEqual, 0)
		})

		Convey("A switch that is on and supposed to be off goes off", func() {
			mock := mocks.NewNetSwitch(mocks.DESIRED_OFF, mocks.RELAY_ON)
			So(mock.NetSwitch.Sync(), ShouldBeNil)
			<-time.After(time.Second)
			So(mock.PollRequests, ShouldEqual, 1)
			So(mock.SwitchRequests, ShouldEqual, 1)
		})

		Convey("A switch that is off and supposed to be on goes on", func() {
			mock := mocks.NewNetSwitch(mocks.DESIRED_ON, mocks.RELAY_OFF)
			So(mock.NetSwitch.Sync(), ShouldBeNil)
			<-time.After(time.Second)
			So(mock.PollRequests, ShouldEqual, 1)
			So(mock.SwitchRequests, ShouldEqual, 1)
		})

		Convey("A switch that is on and supposed to be on stays on", func() {
			mock := mocks.NewNetSwitch(mocks.DESIRED_ON, mocks.RELAY_ON)
			So(mock.NetSwitch.Sync(), ShouldBeNil)
			<-time.After(time.Second)
			So(mock.PollRequests, ShouldEqual, 1)
			So(mock.SwitchRequests, ShouldEqual, 0)
		})
	})
}
