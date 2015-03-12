package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kr15h/fabsmith/routers"
	"os"
)

func main() {
	configRunmode()
	configClients()
	configDatabase()
	beego.Run()
}

func configClients() {

	// Allow access index.html file
	beego.DirectoryIndex = true

	// Serve self-contained Angular JS applications depending on runmode
	if beego.RunMode == "dev" {
		beego.SetStaticPath("/machines", "clients/machines/dev")
		beego.SetStaticPath("/admin", "clients/admin/dev")
	} else { // prod and any other runmode
		beego.SetStaticPath("/machines", "clients/machines/prod")
		beego.SetStaticPath("/admin", "clients/admin/prod")
	}
}

func configRunmode() {

	// Set Beego runmode from FabSmith env variables
	runmode := os.Getenv("FABSMITH_RUNMODE")
	if runmode != "" {
		beego.RunMode = runmode
	}

	// Print FABSMITH_RUNMODE environment variable
	beego.Trace("FABSMITH_RUNMODE:", os.Getenv("FABSMITH_RUNMODE"))
	beego.Trace("beego.RunMode:", beego.RunMode)
}

func configDatabase() {

	// Get MySQL config from environment variables
	mysqlUser := beego.AppConfig.String("mysqluser")
	if mysqlUser == "" {
		panic("Please set mysqluser in app.conf")
	}

	mysqlPass := beego.AppConfig.String("mysqlpass")
	if mysqlPass == "" {
		panic("Please set mysqlpass in app.conf")
	}

	mysqlHost := beego.AppConfig.String("mysqlhost")
	if mysqlHost == "" {
		mysqlHost = "localhost"
	}

	mysqlPort := beego.AppConfig.String("mysqlport")
	if mysqlPort == "" {
		mysqlPort = "3306"
	}

	mysqlDb := beego.AppConfig.String("mysqldb")
	if mysqlDb == "" {
		panic("Please set mysqldb in app.conf")
	}

	// Build MySQL connection string out of the config variables
	mysqlConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDb)

	// Register MySQL driver and default database for beego ORM
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", mysqlConnString)
}
