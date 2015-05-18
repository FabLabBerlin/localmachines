package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
)

type CrossdomainController struct {
	beego.Controller
}

// Redirect to default html interface
func (this *CrossdomainController) Get() {
	http.ServeFile(this.Ctx.ResponseWriter,
		this.Ctx.Request, "crossdomain.xml")
}
