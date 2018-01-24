package invoices

import (
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestInvoices(t *testing.T) {
	Convey("Testing invoices.Invoice model", t, func() {
		Reset(setup.ResetDB)
		Convey("InvoiceOfMonth()", func() {
			Convey("Auto creates new invoice for current month", func() {
				iv1, err := invoices.GetDraft(1, 123, time.Now())
				So(err, ShouldBeNil)

				iv2, err := invoices.GetDraft(1, 123, time.Now())
				So(err, ShouldBeNil)

				So(iv1.Current, ShouldBeTrue)
				So(iv1.Id, ShouldNotEqual, 0)
				So(iv1.Id, ShouldEqual, iv2.Id)
			})

			Convey("Not auto creates new invoice for past month because in that case there are no purchases", func() {
				t := time.Now().AddDate(0, -1, -1)
				_, err := invoices.GetDraft(1, 123, t)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("GetByFastbillId()", func() {
			Convey("Invoice with fastbill_id 342 does not exist. should return nil", func() {
				iv, err := invoices.GetByFastbillId(342)
				So(err, ShouldBeNil)
				So(iv, ShouldBeNil)
			})
			Convey("Invoice with fastbill_id 323 does exist. should return nil", func() {
				m, err := invoices.GetDraft(1, 123, time.Now())
				m.FastbillId = 323
				m.Save()

				iv, err := invoices.GetByFastbillId(323)
				So(err, ShouldBeNil)
				So(iv, ShouldNotBeNil)
				So(iv.Id, ShouldEqual, m.Id)
			})
		})
	})
}
