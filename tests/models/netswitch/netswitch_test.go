package userTests

import (
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"github.com/kr15h/fabsmith/models"
	. "github.com/kr15h/fabsmith/tests/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	ConfigDB()
}

func TestUsers(t *testing.T) {
	Convey("Testing Netswitch model", t, func() {
		Reset(func() {
			o := orm.NewOrm()
			var netswitches []models.NetSwitchMapping
			o.QueryTable("netswitch").All(&netswitches)
			for _, item := range netswitches {
				o.Delete(&item)
			}
		})
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
				err := models.UpdateNetSwitchMapping(netswitch)
				netswitch, _ = models.GetNetSwitchMapping(nid)

				So(err, ShouldBeNil)
				So(netswitch.UrlOff, ShouldEqual, urlOffTest)
			})
			Convey("Trying to call the function with null parameter", func() {
				panicFunc := func() {
					models.UpdateNetSwitchMapping(nil)
				}

				So(panicFunc, ShouldPanic)
			})
		})
	})
}
