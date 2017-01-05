package filter

import (
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
)

func Invoices(all []*invutil.Invoice, from, to month.Month) (ivs []*invutil.Invoice) {
	ivs = make([]*invutil.Invoice, 0, len(all))

	for _, iv := range all {
		if iv.GetMonth().Before(from) {
			continue
		}

		if iv.GetMonth().After(to) {
			continue
		}

		if !iv.Canceled {
			ivs = append(ivs, iv)
		}
	}

	return
}
