package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
)

type AppleTouchIconController struct {
	beego.Controller
}

// Redirect to default html interface
func (this *AppleTouchIconController) Get() {
	http.ServeFile(this.Ctx.ResponseWriter,
		this.Ctx.Request, "files/apple-touch-icon.png")
}
