package products

import (
	"github.com/FabLabBerlin/localmachines/models/products"
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
			c, err := products.CreateCoWorking(1, "foo")
			c2, err2 := products.GetCoWorking(c.Product.Id)
			assert.NoErrors(err, err2)
			So(c2.Product.Name, ShouldEqual, "foo")
		})

		Convey("GetAllCoWorkings", func() {
			c1, err1 := products.CreateCoWorking(1, "foo")
			c2, err2 := products.CreateCoWorking(1, "bar")
			cs, err := products.GetAllCoWorking()
			assert.NoErrors(err1, err2, err)
			So(len(cs), ShouldEqual, 2)
			So(c1.Product.Id, ShouldEqual, cs[0].Product.Id)
			So(c2.Product.Id, ShouldEqual, cs[1].Product.Id)
		})

		Convey("Update", func() {
			c, err := products.CreateCoWorking(1, "foobar")
			c.Product.Name = "foobaz"
			err2 := c.Update()
			cc, err3 := products.GetCoWorking(c.Product.Id)
			assert.NoErrors(err, err2, err3)
			So("foobaz", ShouldEqual, cc.Product.Name)
		})
	})
}
