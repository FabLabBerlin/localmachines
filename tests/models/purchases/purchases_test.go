package purchases

import (
	"github.com/kr15h/fabsmith/models/purchases"
	"github.com/kr15h/fabsmith/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestPurchases(t *testing.T) {
	Convey("Testing general purchase func", t, func() {

		Reset(setup.ResetDB)

		Convey("Creating a purchase", func() {
			purchase := purchases.Purchase{}
			id, err := purchases.Create(&purchase)
			Convey("should return no error", func() {
				So(err, ShouldEqual, nil)
			})
			Convey("should return a valid id", func() {
				So(id, ShouldBeGreaterThan, 0)
			})
			Convey("should set the purchase id equal to the returned id", func() {
				So(purchase.Id, ShouldEqual, id)
			})
		})
		/*
			Convey("When archiving a purchase", func() {

				purchase := Purchase{}

				// I am adding all of the following fields because we should not lose
				// lose the values during archiving.
				purchase.Type = TYPE_TUTOR
				purchase.ProductId = 2
				purchase.Created = time.Now()
				purchase.UserId = 1
				purchase.TimeStart = time.Now()
				purchase.TimeEnd = time.Now()
				purchase.Quantity = 2
				purchase.PricePerUnit = 23
				purchase.PriceUnit = "dolla"
				purchase.Vat = 3
				purchase.Cancelled = false

				// Create a purchase so we can test archiving
				models.
					Convey("")
				s := purchases.NewSpace()
				s.UserId = 234
				err1 := s.Save()
				s, err2 := purchases.GetSpace(s.Id)
				assert.NoErrors(err1, err2)
				So(s.Purchase.UserId, ShouldEqual, 234)
			})
		*/

	})
}
