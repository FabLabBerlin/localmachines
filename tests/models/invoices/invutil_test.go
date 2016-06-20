package invoices

import (
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestInvutilInvoices(t *testing.T) {
	Convey("Testing invutil.Invoice model", t, func() {
		Reset(setup.ResetDB)
		Convey("InvoiceOfMonth()", func() {
		})
	})
}
