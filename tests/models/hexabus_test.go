package modelTest

import (
	"testing"

	"github.com/kr15h/fabsmith/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

func TestHexabus(t *testing.T) {
	Convey("Testing Hexabus model", t, func() {
		Reset(ResetDB)
		Convey("Testing CreateHexabusMapping", func() {
			Convey("Creating a Hexabus mapping regulary", func() {
				nid, err := models.CreateHexabusMapping(0)

				So(err, ShouldBeNil)
				So(nid, ShouldBeGreaterThan, 0)
			})
			Convey("Creating two Hexabus mapping on the same machine", func() {
				nid1, err1 := models.CreateHexabusMapping(0)
				nid2, err2 := models.CreateHexabusMapping(0)

				So(err1, ShouldBeNil)
				So(nid1, ShouldBeGreaterThan, 0)
				So(err2, ShouldBeNil)
				So(nid2, ShouldBeGreaterThan, 0)
			})
		})
		Convey("Testing GetHexabusMapping", func() {
			Convey("Getting netswith on non-existing machine", func() {
				_, err := models.GetHexabusMapping(0)

				So(err, ShouldNotBeNil)
			})
			Convey("Creating netswith and getting it ", func() {
				nid, _ := models.CreateHexabusMapping(0)
				Hexabus, err := models.GetHexabusMapping(0)

				So(err, ShouldBeNil)
				So(Hexabus.Id, ShouldEqual, nid)
			})
		})
		Convey("Testing DeleteHexabusMapping", func() {
			Convey("Creating a Hexabus and deleting it", func() {
				models.CreateHexabusMapping(0)
				err := models.DeleteHexabusMapping(0)

				So(err, ShouldBeNil)
			})
			Convey("Deleting a non-existing Hexabus mapping", func() {
				err := models.DeleteHexabusMapping(0)

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing UpdateHexabusMapping", func() {
			ipTest := "QQ"
			SkipConvey("Creating a Hexabus and updating it", func() {
				nid, _ := models.CreateHexabusMapping(0)
				hexabus, _ := models.GetHexabusMapping(nid)
				hexabus.SwitchIp = ipTest
				err := models.UpdateHexabusMapping(hexabus)
				hexabus, _ = models.GetHexabusMapping(nid)

				So(err, ShouldBeNil)
				So(hexabus.SwitchIp, ShouldEqual, ipTest)
			})
			Convey("Trying to call the function with null parameter", func() {
				panicFunc := func() {
					models.UpdateHexabusMapping(nil)
				}

				So(panicFunc, ShouldPanic)
			})
		})
	})
}
