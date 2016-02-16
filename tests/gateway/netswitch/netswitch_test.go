package gatewayNetswitchTest

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNetswitch(t *testing.T) {
	Convey("Testing Sync", t, func() {
		Convey("A switch that is off and supposed to be off stays off", func() {
			mock := NewMockNetSwitch(DESIRED_OFF, RELAY_OFF)
			So(mock.NetSwitch.Sync(), ShouldBeNil)
			<-time.After(time.Second)
			So(mock.UrlCalled, ShouldEqual, 1)
		})

		Convey("A switch that is on and supposed to be off goes off", func() {
			mock := NewMockNetSwitch(DESIRED_OFF, RELAY_ON)
			So(mock.NetSwitch.Sync(), ShouldBeNil)
			<-time.After(time.Second)
			So(mock.UrlCalled, ShouldEqual, 2)
		})

		Convey("A switch that is off and supposed to be on goes on", func() {
			mock := NewMockNetSwitch(DESIRED_ON, RELAY_OFF)
			So(mock.NetSwitch.Sync(), ShouldBeNil)
			<-time.After(time.Second)
			So(mock.UrlCalled, ShouldEqual, 2)
		})

		Convey("A switch that is on and supposed to be on stays on", func() {
			mock := NewMockNetSwitch(DESIRED_ON, RELAY_ON)
			So(mock.NetSwitch.Sync(), ShouldBeNil)
			<-time.After(time.Second)
			So(mock.UrlCalled, ShouldEqual, 1)
		})
	})
}
