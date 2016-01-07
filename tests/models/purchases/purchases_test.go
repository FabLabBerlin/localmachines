package purchases

import (
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
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
				So(err, ShouldBeNil)
			})
			Convey("should return a valid id", func() {
				So(id, ShouldBeGreaterThan, 0)
			})
			Convey("should set the purchase id equal to the returned id", func() {
				So(purchase.Id, ShouldEqual, id)
			})
		})

		Convey("Archiving a purchase", func() {

			purchase := purchases.Purchase{}

			// I am adding all of the following fields because we should not lose
			// lose the values during archiving.
			purchase.Type = purchases.TYPE_TUTOR
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

			id, _ := purchases.Create(&purchase) //
			err := purchases.Archive(&purchase)  // Archive purchase
			ap, aperr := purchases.Get(id)       // Get archived purchase
			Convey("should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("it should be possible to read the archived purchase back via Get", func() {
				So(aperr, ShouldBeNil)
			})
			Convey("the archived purchase should not be included in GetAllOfType", func() {
				/*
				  TYPE_ACTIVATION  = "activation"
				  TYPE_CO_WORKING  = "co-working"
				  TYPE_RESERVATION = "reservation"
				  TYPE_SPACE       = "space"
				  TYPE_TUTOR       = "tutor"
				*/
				var p []*purchases.Purchase
				var err error
				p, err = purchases.GetAllOfType(purchases.TYPE_ACTIVATION)
				So(len(p), ShouldBeZeroValue)
				So(err, ShouldBeNil)
				p, err = purchases.GetAllOfType(purchases.TYPE_CO_WORKING)
				So(len(p), ShouldBeZeroValue)
				So(err, ShouldBeNil)
				p, err = purchases.GetAllOfType(purchases.TYPE_RESERVATION)
				So(len(p), ShouldBeZeroValue)
				So(err, ShouldBeNil)
				p, err = purchases.GetAllOfType(purchases.TYPE_SPACE)
				So(len(p), ShouldBeZeroValue)
				So(err, ShouldBeNil)
				p, err = purchases.GetAllOfType(purchases.TYPE_TUTOR)
				So(len(p), ShouldBeZeroValue)
				So(err, ShouldBeNil)
			})
			Convey("the archived purchase should contain previously stored data", func() {
				So(ap.Type, ShouldEqual, purchase.Type)
				So(ap.ProductId, ShouldEqual, purchase.ProductId)
				So(ap.Created, ShouldHappenWithin, time.Duration(1)*time.Second, purchase.Created)
				So(ap.UserId, ShouldEqual, purchase.UserId)
				So(ap.TimeStart, ShouldHappenWithin, time.Duration(1)*time.Second, purchase.TimeStart)
				So(ap.TimeEnd, ShouldHappenWithin, time.Duration(1)*time.Second, purchase.TimeEnd)
				So(ap.Quantity, ShouldEqual, purchase.Quantity)
				So(ap.PricePerUnit, ShouldEqual, purchase.PricePerUnit)
				So(ap.PriceUnit, ShouldEqual, purchase.PriceUnit)
				So(ap.Vat, ShouldEqual, purchase.Vat)
				So(ap.Cancelled, ShouldEqual, purchase.Cancelled)
			})

		})

	})
}
