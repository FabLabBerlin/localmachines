package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
)

func TaskFastbillSync() (err error) {
	beego.Info("Running TaskFastbillSync")

	var locId int64 = 1

	us, err := users.GetAllUsersAt(locId)
	if err != nil {
		return fmt.Errorf("get all users at: %v", err)
	}

	for _, u := range us {
		if u.ClientId == 0 {
			continue
		}

		if err := FastbillSync(locId, u); err != nil {
			beego.Error("sync invoice of user %v: %v", u.Id, err)
		}
	}

	return
}

func TaskCalculateTotals() (err error) {
	beego.Info("Running TaskCalculateTotals")

	ls, err := locations.GetAll()
	if err != nil {
		return fmt.Errorf("get all locations: %v", err)
	}
	for _, l := range ls {
		invs, err := GetAllAt(l.Id)
		if err != nil {
			return fmt.Errorf("get all invoices @ %v: %v", l.Id, err)
		}
		for _, inv := range invs {
			if err := inv.CalculateTotals(); err != nil {
				return fmt.Errorf("calculate totals for %v: %v", inv.Id, err)
			}
		}
	}

	return
}
