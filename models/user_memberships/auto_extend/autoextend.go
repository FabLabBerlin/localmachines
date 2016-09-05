package auto_extend

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"sync"
	"time"
)

var mu sync.Mutex

func Lock() {
	mu.Lock()
}

func Unlock() {
	mu.Unlock()
}

// Automatically extend user membership end date if auto_extend for the specific
// membership is true and the end_date is before current date.
func RunTask() (err error) {

	mu.Lock()
	defer mu.Unlock()

	beego.Info("Running AutoExtendUserMemberships Task")

	if err = AutoExtendUserMemberships(time.Now()); err != nil {
		beego.Error("Failed to get all user memberships:", err)
	}

	return
}

func AutoExtendUserMemberships(minimumTime time.Time) (err error) {
	ls, err := locations.GetAll()
	if err != nil {
		return fmt.Errorf("get all locations: %v", err)
	}

	for _, l := range ls {
		if err := extendUserMembershipsAt(l.Id, minimumTime); err != nil {
			return fmt.Errorf("extend userMemberships at %v: %v", l.Id, err)
		}
	}

	return
}

func extendUserMembershipsAt(locId int64, minimumTime time.Time) (err error) {
	beego.Info("for all users invutil.AssureUserHasDraftFor", locId, "begin")
	y := minimumTime.Year()
	m := minimumTime.Month()

	us, err := users.GetAllUsersAt(locId)
	if err != nil {
		return fmt.Errorf("GetAllUsersAt: %v", err)
	}

	for _, u := range us {
		if err := invutil.AssureUserHasDraftFor(locId, u, y, m); err != nil {
			return fmt.Errorf("AssureUserHasDraftFor loc %v: %v", locId, err)
		}
	}
	beego.Info("invutil.AssureUserHasDraftFor done")

	ums, err := user_memberships.GetAllAt(locId)
	if err != nil {
		return fmt.Errorf("get all user memberships: %v", err)
	}

	for _, um := range ums {
		if !um.AutoExtend || um.EndDate.After(minimumTime) {
			continue
		}

		m, err := memberships.Get(um.MembershipId)
		if err != nil {
			return fmt.Errorf("get membership: %v", err)
		}

		inv, err := invoices.Get(um.InvoiceId)
		if err != nil {
			return fmt.Errorf("get invoice: %v", err)
		}

		if inv.Month != int(minimumTime.Month()) ||
			inv.Year != minimumTime.Year() {
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
