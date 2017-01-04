package filter

import (
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
)

func Invoices(all []*invutil.Invoice, from, to day.Day) (ivs []*invutil.Invoice) {
	ivs = make([]*invutil.Invoice, 0, len(all))

	for _, iv := range all {
		if iv.GetMonth().Before(from.Month()) {
			continue
		}

		if iv.GetMonth().After(to.Month()) {
			continue
		}

		if !iv.Canceled {
			ivs = append(ivs, iv)
		}
	}

	return
}
