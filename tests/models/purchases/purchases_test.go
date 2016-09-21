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
			purchase := purchases.Purchase{
				LocationId: 1,
				InvoiceId:  11,
			}
			err := purchases.Create(&purchase)
			Convey("should return no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("should return a valid id", func() {
				So(purchase.Id, ShouldBeGreaterThan, 0)
			})
		})

		Convey("Archiving a purchase", func() {

			// I am adding all of the following fields because we should not
			// lose the values during archiving.
			purchase := purchases.Purchase{
				LocationId:   1,
				Type:         purchases.TYPE_TUTOR,
				InvoiceId:    111,
				ProductId:    2,
				Created:      time.Now(),
				UserId:       1,
				TimeStart:    time.Now(),
				TimeEnd:      time.Now(),
				Quantity:     2,
				PricePerUnit: 23,
				PriceUnit:    "dolla",
				Vat:          3,
				Cancelled:    false,
			}

			if err := purchases.Create(&purchase); err != nil {
				panic(err.Error())
			}
			err := purchases.Archive(&purchase)     // Archive purchase
			ap, aperr := purchases.Get(purchase.Id) // Get archived purchase
			Convey("should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("it should be possible to read the archived purchase back via Get", func() {
				So(aperr, ShouldBeNil)
			})
			Convey("the archived purchase should not be included in GetAllOfType", func() {
				p, err := purchases.GetAllOfType(purchases.TYPE_ACTIVATION)
				So(len(p), ShouldBeZeroValue)
				So(err, ShouldBeNil)
				p, err = purchases.GetAllOfType(purchases.TYPE_RESERVATION)
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
