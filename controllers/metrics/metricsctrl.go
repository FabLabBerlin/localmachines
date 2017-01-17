// /api/metrics
package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/invoices/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/metrics"
	"github.com/FabLabBerlin/localmachines/models/metrics/bin"
	"github.com/FabLabBerlin/localmachines/models/metrics/filter"
	"github.com/FabLabBerlin/localmachines/models/metrics/heatmap"
	"github.com/FabLabBerlin/localmachines/models/metrics/machine_capacity"
	"github.com/FabLabBerlin/localmachines/models/metrics/machine_earnings"
	"github.com/FabLabBerlin/localmachines/models/metrics/memberstats"
	"github.com/FabLabBerlin/localmachines/models/metrics/retention"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

type Controller struct {
	controllers.Controller
}

func (c *Controller) FromTo() (from, to day.Day, err error) {
	from = day.New(2015, 8, 1)
	to = day.Now().AddDate(0, 0, -1)

	if fromString := c.GetString("from"); fromString != "" {
		if from, err = day.NewString(fromString); err != nil {
			return
		}
	}

	if toString := c.GetString("to"); toString != "" {
		if to, err = day.NewString(toString); err != nil {
			return
		}
	}

	return
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

	from, to, err := c.FromTo()
	if err != nil {
		c.Fail(400, fmt.Sprintf("from/to: %v", err))
	}

	interval := lib.Interval{
		MonthFrom: int(from.Month().Month()),
		YearFrom:  from.Year(),
		MonthTo:   int(to.Month().Month()),
		YearTo:    to.Year(),
	}

	binWidth := bin.NewWidth(bin.Unit(c.GetString("binwidth")))

	data, err := metrics.FetchData(locId, interval, binWidth)
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

type HeatmapDataPoint struct {
	Coordinate heatmap.Coordinate
	UserId     int64
}

// @Title Get Heatmap
// @Description Get some nice Heatmap data
// @Success 200
// @Failure	500	Failed to get heat map data
// @Failure	401	Not authorized
// @router /heatmap [get]
func (c *Controller) GetHeatmap() {
	locId, authorized := c.GetLocIdAdmin()
	if !authorized {
		c.CustomAbort(401, "Not authorized")
	}

	us, err := users.GetAllUsersAt(locId)
	if err != nil {
		c.Fail(500, err.Error())
	}

	ps := make([]HeatmapDataPoint, 0, len(us))
	for _, u := range us {
		var c *heatmap.Coordinate
		redis.Get(fmt.Sprintf("geocode(%v)", u.Id), &c)
		if c != nil {
			p := HeatmapDataPoint{
				Coordinate: *c,
				UserId:     u.Id,
			}
			ps = append(ps, p)
		}
	}

	c.Data["json"] = ps
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

	from, to, err := c.FromTo()
	if err != nil {
		c.Fail(400, fmt.Sprintf("from/to: %v", err))
	}

	for _, machine := range machines {
		res := make(map[string]interface{})

		mc := machine_capacity.New(
			machine,
			from.Month(),
			to.Month(),
			invs,
		)

		res["Machine"] = machine
		capacity, err := mc.CapacityCached()
		if err != nil {
			c.Fail(500, err.Error())
		}
		res["Capacity"] = capacity.Hours() / 24
		usage, err := mc.UsageCached()
		if err != nil {
			c.Fail(500, err.Error())
		}
		res["Hours"] = usage.Hours() / 24
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

	from, to, err := c.FromTo()
	if err != nil {
		c.Fail(400, fmt.Sprintf("from/to: %v", err))
	}

	for _, machine := range machines {
		res := make(map[string]interface{})

		me := machine_earnings.New(
			machine,
			from.Month(),
			to.Month(),
			invs,
		)

		res["Machine"] = machine
		res["Memberships"] = me.MembershipsCached()
		payg, err := me.PayAsYouGoCached()
		if err != nil {
			c.Fail(500, err.Error())
		}
		res["PayAsYouGo"] = payg
		resp = append(resp, res)
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title Get user memberships
// @Description Get user memberships
// @Success 200
// @Failure	500	Failed to get user memberships
// @Failure	401	Not authorized
// @router /memberships [get]
func (c *Controller) GetMemberships() {
	locId, authorized := c.GetLocIdAdmin()
	if !authorized {
		c.CustomAbort(401, "Not authorized")
	}

	from, to, err := c.FromTo()
	if err != nil {
		c.Fail(400, fmt.Sprintf("from/to: %v", err))
	}

	invs, err := invutil.GetAllAt(locId)
	if err != nil {
		c.Fail(500, "Failed to get invoices")
	}
	fmt.Printf("1\n")
	ms := memberstats.New(
		from.Month(),
		to.Month(),
		invs,
	)
	fmt.Printf("11\n")

	bins /*, err*/ := ms.Bins /*Cached*/ ()
	/*if err != nil {
		c.Fail(500, err.Error())
	}*/
	fmt.Printf("111\n")
	c.Data["json"] = bins
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

	allInvs, err := invutil.GetAllAt(locId)
	if err != nil {
		c.Fail(500, "Failed to get invoices")
	}

	invs := filter.InvoicesByUsers(locId, allInvs, uls, true, excludeNeverActive)

	us := make([]*users.User, 0, len(allUsers))

	for _, inv := range invs {
		us = append(us, inv.User)
	}

	from, to, err := c.FromTo()
	if err != nil {
		c.Fail(400, fmt.Sprintf("from/to: %v", err))
	}

	r := retention.New(
		locId,
		30,
		from,
		to,
		invs,
		us,
	)

	c.Data["json"] = r.Calculate()
	c.ServeJSON()
}
