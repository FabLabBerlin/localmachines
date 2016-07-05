package purchases

import (
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/tests/assert"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func init() {
	setup.ConfigDB()
}

func TestTutoringPurchases(t *testing.T) {
	Convey("Testing Tutoring purchase model", t, func() {

		Reset(setup.ResetDB)

		inv := &invutil.Invoice{}
		inv.LocationId = 1
		inv.UserId = 123
		inv.Month = int(time.Now().Month())
		inv.Year = time.Now().Year()
		inv.Status = "draft"
		if _, err := invoices.Create(&inv.Invoice); err != nil {
			panic(err.Error())
		}

		Convey("CreateTutoring and GetTutoring", func() {
			t := &purchases.Tutoring{
				Purchase: purchases.Purchase{
					UserId:     123,
					LocationId: 1,
					InvoiceId:  inv.Id,
				},
			}
			id, err1 := purchases.CreateTutoring(t)
			t, err2 := purchases.GetTutoring(id)
			assert.NoErrors(err1, err2)
			So(t.Purchase.UserId, ShouldEqual, 123)
		})

		Convey("GetAllTutorings", func() {
			id1, err1 := purchases.CreateTutoring(&purchases.Tutoring{
				Purchase: purchases.Purchase{
					UserId:     123,
					LocationId: 1,
					InvoiceId:  inv.Id,
				},
			})
			id2, err2 := purchases.CreateTutoring(&purchases.Tutoring{
				Purchase: purchases.Purchase{
					UserId:     123,
					LocationId: 1,
					InvoiceId:  inv.Id,
				},
			})
			ts, err := purchases.GetAllTutorings()
			assert.NoErrors(err1, err2, err)
			So(len(ts.Data), ShouldEqual, 2)
			So(id1, ShouldEqual, ts.Data[0].Id)
			So(id2, ShouldEqual, ts.Data[1].Id)
		})

		Convey("Update", func() {
			t := &purchases.Tutoring{
				Purchase: purchases.Purchase{
					UserId:     123,
					LocationId: 1,
					InvoiceId:  inv.Id,
				},
			}
			id, err := purchases.CreateTutoring(t)
			t.Purchase.UserId = 456
			err2 := t.Update()
			t, err3 := purchases.GetTutoring(id)
			assert.NoErrors(err, err2, err3)
			So(456, ShouldEqual, t.Purchase.UserId)
		})

	})
}
