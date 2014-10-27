package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kr15h/fabsmith/routers"
)

func init() {
	// Debug init
	// Shows file names and line numbers in debug output when beego.Info("asd"),
	// beego.Error("omg") etc. debug functions are used
	beego.SetLogFuncCall(true)
	// Set log level:
	// beego.SetLevel(beego.LevelInfo)
	// Available options are
	// beego.LevelDebug
	// beego.LevelInformational
	// beego.LevelWarning
	// beego.LevelError
	// beego.LevelCritical
	// See more in Beego docs or source

	// Template init
	beego.TemplateLeft = "<<<"
	beego.TemplateRight = ">>>"

	// Database init
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@/fabsmith?charset=utf8")
}

func main() {
	beego.Run()
}
