package clients

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/lib/cache"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"os"
)

var runMode = beego.AppConfig.String("RunMode")

type Machines struct {
	controllers.Controller
}

// @Title Get
// @Description Get Machines Page
// @Success 200 Redirect
// @Failure	500	Internal Server Error
// @router / [get]
func (c *Machines) Get() {
	var fn string
	if runMode == "dev" {
		fn = "clients/machines/dev/index.html"
	} else {
		fn = "clients/machines/prod/index.html"
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
				beego.Error("Machines: Get: Cannot cast to string")
				c.CustomAbort(500, "Internal Server Error")
			}
		} else {
			beego.Info("Cache miss")
			bs, err := ioutil.ReadFile(fn)
			if err != nil {
				beego.Error("Machines: Get: Error reading file:", err)
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
