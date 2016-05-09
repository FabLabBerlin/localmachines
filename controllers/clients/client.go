package clients

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/lib/cache"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"os"
)

var runMode = beego.AppConfig.String("RunMode")

func get(c *controllers.Controller, clientName string) {
	if c.Ctx.Request.URL.Path == "/"+clientName {
		c.Redirect("/"+clientName+"/", 302)
		return
	}

	var fn string
	if runMode == "dev" {
		fn = "clients/" + clientName + "/dev/index.html"
	} else {
		fn = "clients/" + clientName + "/prod/index.html"
	}
	locId, ok := c.GetSessionLocationId()
	if !ok {
		clientIp := c.ClientIp()
		// Try to get locId based on IP
		if locs, err := locations.GetAll(); err == nil {
			for _, loc := range locs {
				if loc.LocalIp == clientIp {
					locId = loc.Id
					c.SetSessionLocationId(locId)
					break
				}
			}
		}
	}
	c.Ctx.Output.ContentType("html")
	if runMode == "dev" {
		f, err := os.Open(fn)
		if err != nil {
			beego.Error("cannot open ", fn, ":", err)
		}
		defer f.Close()
		if _, err := io.Copy(c.Ctx.ResponseWriter, f); err != nil {
			beego.Error("io copy:", err)
		}
	} else {
		if o, ok := cache.Get(fn); ok {
			beego.Info("Cache hit")
			if html, ok := o.(string); ok {
				c.Ctx.WriteString(html)
				c.Finish()
			} else {
				beego.Error(clientName + " client: Get: Cannot cast to string")
				c.CustomAbort(500, "Internal Server Error")
			}
		} else {
			beego.Info("Cache miss")
			bs, err := ioutil.ReadFile(fn)
			if err != nil {
				beego.Error(clientName+" client: Get: Error reading file:", err)
				c.CustomAbort(500, "Internal Server Error")
			}
			s := string(bs)
			cache.Put(fn, s)
			c.Ctx.WriteString(s)
			c.Finish()
		}
	}
	c.Finish()
}
