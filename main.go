package main

import (
	"github.com/astaxie/beego"
	_ "github.com/kr15h/fabsmith/routers"
)

func init() {
	beego.TemplateLeft = "<<<"
	beego.TemplateRight = ">>>"
}

func main() {
	beego.Run()
}
