package invutil

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/users"
)

func SyncFastBillAttributesTask() (err error) {
	var locId int64 = 1

	us, err := users.GetAllUsersAt(locId)
	if err != nil {
		return fmt.Errorf("get all users at: %v", err)
	}

	for _, u := range us {
		if err := SyncFastbillInvoices(locId, u); err != nil {
			return fmt.Errorf("sync invoice of user %v", u.Id)
		}
	}

	return
}
