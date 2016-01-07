package purchases

import (
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/tests/assert"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestCoWorkingPurchases(t *testing.T) {
	Convey("Testing CoWorking purchase model", t, func() {

		Reset(setup.ResetDB)

		Convey("CreateCoWorking and GetCoWorking", func() {
			c := &purchases.CoWorking{
				Purchase: purchases.Purchase{
					UserId: 234,
				},
			}
			id, err1 := purchases.CreateCoWorking(c)
			c, err2 := purchases.GetCoWorking(id)
			assert.NoErrors(err1, err2)
			So(c.Purchase.UserId, ShouldEqual, 234)
		})

		Convey("GetAllCoWorkings", func() {
			id1, err1 := purchases.CreateCoWorking(&purchases.CoWorking{})
			id2, err2 := purchases.CreateCoWorking(&purchases.CoWorking{})
			cs, err := purchases.GetAllCoWorking()
			assert.NoErrors(err1, err2, err)
			So(len(cs), ShouldEqual, 2)
			So(id1, ShouldEqual, cs[0].Id)
			So(id2, ShouldEqual, cs[1].Id)
		})

		Convey("Update", func() {
			c := &purchases.CoWorking{
				Purchase: purchases.Purchase{
					UserId: 123,
				},
			}
			id, err := purchases.CreateCoWorking(c)
			c.Purchase.UserId = 456
			err2 := c.Update()
			c, err3 := purchases.GetCoWorking(id)
			assert.NoErrors(err, err2, err3)
			So(456, ShouldEqual, c.Purchase.UserId)
		})

	})
}
