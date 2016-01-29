package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/metrics"
	"github.com/FabLabBerlin/localmachines/models/purchases"
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
	// Only admin can use this API call
	if !c.IsAdmin() {
		beego.Error("Not authorized")
		c.CustomAbort(401, "Not authorized")
	}

	data, err := metrics.FetchData()
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
	if !c.IsAdmin() {
		c.CustomAbort(401, "Not authorized")
	}

	endTime := time.Now()
	startTime := time.Date(2015, time.August, 1, 0, 0, 0, 0, time.UTC)

	invoice, err := invoices.CalculateSummary(startTime, endTime)
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
	for _, summary := range invoice.UserSummaries {
		for _, p := range summary.Purchases.Data {
			row := make([]string, 0, 20)
			if p.Type != purchases.TYPE_ACTIVATION {
				continue
			}
			row = append(row, p.TimeStart.String())
			row = append(row, p.User.FirstName+" "+p.User.LastName, p.User.Email)
			row = append(row, p.Machine.Name)
			if p.TimeEnd.IsZero() {
				row = append(row, "0")
			} else {
				mins := fmt.Sprintf("%v", p.TimeEnd.Sub(p.TimeStart).Minutes())
				mins = strings.Replace(mins, ".", ",", -1)
				row = append(row, mins)
			}
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
