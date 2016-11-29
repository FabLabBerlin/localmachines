// /api/metrics
package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/invoices/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/metrics"
	"github.com/FabLabBerlin/localmachines/models/metrics/machine_capacity"
	"github.com/FabLabBerlin/localmachines/models/metrics/machine_earnings"
	"github.com/FabLabBerlin/localmachines/models/metrics/retention"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

type Controller struct {
	controllers.Controller
}

// @Title Get All
// @Description Get all metrics
// @Success 200
// @Failure	500	Failed to get metrics
// @Failure	401	Not authorized
// @router / [get]
func (c *Controller) GetAll() {
	// Only local admin can use this API call
	locId, authorized := c.GetLocIdAdmin()
	if !authorized {
		c.CustomAbort(401, "Not authorized")
	}

	data, err := metrics.FetchData(locId)
	if err != nil {
		beego.Error("Failed to get metrics data:", err)
		c.CustomAbort(500, "Failed to get metrics data")
	}

	resp, err := metrics.NewResponse(data)
	if err != nil {
		beego.Error("Failed to get metrics response:", err)
		c.CustomAbort(500, "Failed to get metrics response")
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title Get All Activations fully denormalized
// @Description Get all activations
// @Success 200
// @Failure	500	Failed to get activations
// @Failure	401	Not authorized
// @router /activations [get]
func (c *Controller) GetActivations() {
	locId, authorized := c.GetLocIdAdmin()
	if !authorized {
		c.CustomAbort(401, "Not authorized")
	}

	endTime := time.Now()

	interval := lib.Interval{
		MonthFrom: int(time.August),
		YearFrom:  2015,
		MonthTo:   int(endTime.Month()),
		YearTo:    endTime.Year(),
	}

	monthlyEarning, err := monthly_earning.New(locId, interval)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}

	buf := bytes.NewBufferString("")
	w := csv.NewWriter(buf)
	w.Comma = ';'
	w.Write([]string{
		"Time Start",
		"Name",
		"E-Mail",
		"Machine",
		"Duration (minutes)",
		"Membership",
		"Billed Price (EUR)",
	})
	for _, inv := range monthlyEarning.Invoices {
		for _, p := range inv.Purchases {
			row := make([]string, 0, 20)
			if p.Type != purchases.TYPE_ACTIVATION {
				continue
			}
			row = append(row, p.TimeStart.String())
			row = append(row, p.User.FirstName+" "+p.User.LastName, p.User.Email)
			row = append(row, p.Machine.Name)
			mins := fmt.Sprintf("%v", p.TimeEnd().Sub(p.TimeStart).Minutes())
			mins = strings.Replace(mins, ".", ",", -1)
			row = append(row, mins)
			membershipStr, err := p.MembershipStr()
			if err != nil {
				c.CustomAbort(500, err.Error())
			}
			row = append(row, membershipStr)
			discTotal := fmt.Sprintf("%v", p.DiscountedTotal)
			discTotal = strings.Replace(discTotal, ".", ",", -1)
			row = append(row, discTotal)
			w.Write(row)
		}
	}
	w.Flush()
	c.Ctx.WriteString(buf.String())
	c.ServeJSON()
}

// @Title Get machine capacities
// @Description Get all Machine Capacities
// @Success 200
// @Failure	500	Failed to get machine capacities
// @Failure	401	Not authorized
// @router /machine_capacities [get]
func (c *Controller) GetMachineCapacities() {
	locId, authorized := c.GetLocIdAdmin()
	if !authorized {
		c.CustomAbort(401, "Not authorized")
	}

	machines, err := machine.GetAllAt(locId)
	if err != nil {
		beego.Error("Failed to get machines:", err)
		c.Abort("500")
	}

	invs, err := invutil.GetAllAt(locId)
	if err != nil {
		c.Fail(500, "Failed to get invoices")
	}

	resp := make([]interface{}, 0, 40)

	for _, machine := range machines {
		res := make(map[string]interface{})

		mc := machine_capacity.New(
			machine,
			month.New(1, 2015),
			month.New(12, 2017),
			invs,
		)

		res["Machine"] = machine
		res["Capacity"] = mc.CapacityCached().Hours() / 24
		res["Hours"] = mc.UsageCached().Hours() / 24
		res["Utilization"] = mc.Utilization()
		resp = append(resp, res)
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title Get machine earnings
// @Description Get all Machine Earnings
// @Success 200
// @Failure	500	Failed to get machine earnings
// @Failure	401	Not authorized
// @router /machine_earnings [get]
func (c *Controller) GetMachineEarnings() {
	locId, authorized := c.GetLocIdAdmin()
	if !authorized {
		c.CustomAbort(401, "Not authorized")
	}

	machines, err := machine.GetAllAt(locId)
	if err != nil {
		beego.Error("Failed to get machines:", err)
		c.Abort("500")
	}

	invs, err := invutil.GetAllAt(locId)
	if err != nil {
		c.Fail(500, "Failed to get invoices")
	}

	resp := make([]interface{}, 0, 40)

	for _, machine := range machines {
		res := make(map[string]interface{})

		me := machine_earnings.New(
			machine,
			month.New(1, 2015),
			month.New(12, 2017),
			invs,
		)

		res["Machine"] = machine
		res["Memberships"] = me.MembershipsCached()
		res["PayAsYouGo"] = me.PayAsYouGoCached()
		resp = append(resp, res)
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title Get machine revenue
// @Description Get all activations
// @Success 200
// @Failure	500	Failed to get machine revenue
// @Failure	401	Not authorized
// @router /machine_revenue [get]
func (c *Controller) GetMachineRevenue() {
	if !c.IsAdminAt(1) {
		c.Abort("401")
	}

	s, err := metrics.NewTrotecStats()
	if err != nil {
		beego.Error("trotec stats:", err)
		c.Abort("500")
	}

	c.Data["json"] = s
	c.ServeJSON()
}

// @Title Get user retention
// @Description Get user retention
// @Success 200
// @Failure	500	Failed to get user retention
// @Failure	401	Not authorized
// @router /retention [get]
func (c *Controller) GetRetention() {
	locId, authorized := c.GetLocIdAdmin()
	if !authorized {
		c.CustomAbort(401, "Not authorized")
	}

	excludeNeverActive := c.GetString("excludeNeverActive") == "true"

	allUsers, err := users.GetAllUsersAt(locId)
	if err != nil {
		beego.Error("Failed to get users:", err)
		c.Abort("500")
	}

	uls, err := user_locations.GetAllForLocation(locId)
	if err != nil {
		beego.Error("Failed to get user locations:", err)
		c.Abort("500")
	}

	rolesByUid := make(map[int64]user_roles.Role)
	for _, ul := range uls {
		rolesByUid[ul.UserId] = ul.GetRole()
	}

	everActiveUid := make(map[int64]struct{})

	invs, err := invutil.GetAllAt(locId)
	if err != nil {
		c.Fail(500, "Failed to get invoices")
	}

	for _, inv := range invs {
		for _, p := range inv.Purchases {
			everActiveUid[p.UserId] = struct{}{}
		}
		for _, ium := range inv.InvUserMemberships {
			everActiveUid[ium.UserId] = struct{}{}
		}
	}

	us := make([]*users.User, 0, len(allUsers))
	for _, u := range allUsers {
		if excludeNeverActive {
			if _, everActive := everActiveUid[u.Id]; !everActive {
				continue
			}
		}
		if r, ok := rolesByUid[u.Id]; !ok || r == user_roles.MEMBER {
			us = append(us, u)
		}
	}

	r := retention.New(
		locId,
		30,
		day.New(2015, 7, 1),
		day.Now(),
		invs,
		us,
	)

	c.Data["json"] = r.Calculate()
	c.ServeJSON()
}
