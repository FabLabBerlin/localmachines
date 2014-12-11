package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kr15h/fabsmith/routers"
)

func init() {
	// Debug init
	beego.SetLogFuncCall(false)
	// beego.SetLevel(beego.LevelInfo)

	// Template init, we replace the default template tags
	// as AngularJS uses the same ones as GoLang
	beego.TemplateLeft = "<<<"
	beego.TemplateRight = ">>>"

	// Setup database
	mysqlUser := beego.AppConfig.String("mysqluser")
	mysqlPass := beego.AppConfig.String("mysqlpass")
	mysqlHost := beego.AppConfig.String("mysqlhost")
	mysqlDb := beego.AppConfig.String("mysqldb")

	mysqlConnString := fmt.Sprintf("%s:%s@%s/%s?charset=utf8",
		mysqlUser, mysqlPass, mysqlHost, mysqlDb)

	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", mysqlConnString)
}

func main() {
	beego.Run()
}
