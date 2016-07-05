package invoices

import (
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func init() {
	setup.ConfigDB()
}

func TestInvoices(t *testing.T) {
	Convey("Testing invoices.Invoice model", t, func() {
		Reset(setup.ResetDB)
		Convey("InvoiceOfMonth()", func() {
			Convey("Auto creates new invoice for current month", func() {
				y := time.Now().Year()
				m := time.Now().Month()
				iv1, err := invoices.InvoiceOfMonth(1, 123, y, m)
				So(err, ShouldBeNil)

				iv2, err := invoices.InvoiceOfMonth(1, 123, y, m)
				So(err, ShouldBeNil)

				So(iv1.Current, ShouldBeTrue)
				So(iv1.Id, ShouldNotEqual, 0)
				So(iv1.Id, ShouldEqual, iv2.Id)
			})

			Convey("Not auto creates new invoice for past month because in that case there are no purchases", func() {
				t := time.Now().AddDate(0, -1, 0)
				y := t.Year()
				m := t.Month()
				_, err := invoices.InvoiceOfMonth(1, 123, y, m)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
