package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kr15h/fabsmith/routers"
)

func init() {
	// Debug init
	beego.SetLogFuncCall(false)
	// beego.SetLevel(beego.LevelInfo)

	// Template init
	beego.TemplateLeft = "<<<"
	beego.TemplateRight = ">>>"

	// Database init
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "fabsmith:fabsmith@/fabsmith?charset=utf8")
}

func main() {
	beego.Run()
}
