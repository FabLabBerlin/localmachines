package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
)

type FaviconController struct {
	beego.Controller
}

// Redirect to default html interface
func (this *FaviconController) Get() {
	http.ServeFile(this.Ctx.ResponseWriter,
		this.Ctx.Request, "files/favicon.png")
}
