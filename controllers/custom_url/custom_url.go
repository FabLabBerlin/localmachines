package custom_url

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"log"
	"os"
	"strconv"
)

var locs []*locations.Location

func loadData() (err error) {
	if locs == nil {
		locs, err = locations.GetAll()
	}
	return
}

type Controller struct {
	controllers.Controller
}

// @Title Get
// @Description Serve custom url
// @Success 302 Redirect
// @Failure	404	Not found
// @router / [get]
func (c *Controller) GetAll() {
	if err := loadData(); err != nil {
		// This must work, otherwise the process must be restarted by
		// daemontools. Failing fast.
		log.Printf("Init: %v", err)
		os.Exit(42)
	}
	for _, l := range locs {
		url := "/" + l.Title
		if url == c.Ctx.Request.URL.Path {
			id := strconv.FormatInt(l.Id, 10)
			c.Redirect("/machines/?location="+id+"#/login", 302)
			return
		}
	}
	c.CustomAbort(404, "Not found")
}
