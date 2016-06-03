package auto_extend

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices"
	"github.com/astaxie/beego"
	"time"
)

// Automatically extend user membership end date if auto_extend for the specific
// membership is true and the end_date is before current date.
func AutoExtendUserMemberships() (err error) {

	beego.Info("Running AutoExtendUserMemberships Task")

	if err = autoExtendUserMemberships(); err != nil {
		beego.Error("Failed to get all user memberships:", err)
	}

	return
}

func autoExtendUserMemberships() (err error) {
	ls, err := locations.GetAll()
	if err != nil {
		return fmt.Errorf("get all locations: %v", err)
	}

	for _, l := range ls {
		if err := extendUserMembershipsAt(l.Id); err != nil {
			return fmt.Errorf("extend userMemberships at %v: %v", l.Id, err)
		}
	}

	return
}

func extendUserMembershipsAt(locId int64) (err error) {
	ums, err := memberships.GetAllUserMembershipsAt(locId)
	if err != nil {
		return fmt.Errorf("get all user memberships: %v", err)
	}

	for _, um := range ums {
		if !um.AutoExtend || um.EndDate.After(time.Now()) {
			continue
		}

		m, err := memberships.GetMembership(um.MembershipId)
		if err != nil {
			return fmt.Errorf("get membership: %v", err)
		}

		inv, err := invoices.Get(um.InvoiceId)
		if err != nil {
			return fmt.Errorf("get invoice: %v", err)
		}

		if inv.Month != int(time.Now().Month()) ||
			inv.Year != time.Now().Year() {
			continue
		}

		beego.Trace("Extending user membership with Id", um.Id)

		um.EndDate = um.EndDate.AddDate(0, int(m.AutoExtendDurationMonths), 0)
		if err = um.Update(); err != nil {
			return fmt.Errorf("Failed to update user membership end date")
		}
	}

	return
}
