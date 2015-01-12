package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kr15h/fabsmith/routers"
)

func init() {

	// Set debug level for our app depending on the runmode set
	if v, err := beego.GetConfig("string", "runmode"); err == nil && v == "dev" {

		// Show all messages while in dev mode
		beego.SetLevel(beego.LevelDebug)
	} else {

		// Show only errors when in prod mode
		beego.SetLevel(beego.LevelError)
	}

	// Template init, we replace the default template tags
	// as AngularJS uses the same ones as GoLang
	beego.TemplateLeft = "<<<"
	beego.TemplateRight = ">>>"

	// Get MySQL database variables from config file
	mysqlUser := beego.AppConfig.String("mysqluser")
	mysqlPass := beego.AppConfig.String("mysqlpass")
	mysqlHost := beego.AppConfig.String("mysqlhost")
	mysqlPort := beego.AppConfig.String("mysqlport")
	mysqlDb := beego.AppConfig.String("mysqldb")

	// If MySQL port is empty, replace with default value
	if mysqlPort == "" {
		mysqlPort = "3306"
	}

	// Build MySQL connection string out of the config variables
	mysqlConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDb)

	// Register MySQL driver and default database for beego ORM
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", mysqlConnString)
}

func main() {
	beego.Run()
}
