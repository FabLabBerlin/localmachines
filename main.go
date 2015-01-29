package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kr15h/fabsmith/routers"
	"os"
	"strconv"
)

func main() {
	configServer()
	configTemplate()
	configRunmode()
	configDatabase()

	beego.Run()
}

func configServer() {
	port, err := strconv.Atoi(os.Getenv("FABSMITH_HTTP_PORT"))
	if err == nil {
		beego.HttpPort = port
	}
}

func configTemplate() {

	// Template init, we replace the default template tags
	// as AngularJS uses the same ones as GoLang
	beego.TemplateLeft = "<<<"
	beego.TemplateRight = ">>>"
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
	mysqlUser := os.Getenv("FABSMITH_MYSQL_USER")
	if mysqlUser == "" {
		panic("Please set FABSMITH_MYSQL_USER environment variable")
	}

	mysqlPass := os.Getenv("FABSMITH_MYSQL_PASS")
	if mysqlPass == "" {
		panic("Please set FABSMITH_MYSQL_PASS environment variable")
	}

	mysqlHost := os.Getenv("FABSMITH_MYSQL_HOST")
	if mysqlHost == "" {
		mysqlHost = "localhost"
	}

	mysqlPort := os.Getenv("FABSMITH_MYSQL_PORT")
	if mysqlPort == "" {
		mysqlPort = "3306"
	}

	mysqlDb := os.Getenv("FABSMITH_MYSQL_DB")
	if mysqlDb == "" {
		panic("Please set FABSMITH_MYSQL_DB environment variable")
	}

	// Build MySQL connection string out of the config variables
	mysqlConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDb)

	// Register MySQL driver and default database for beego ORM
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", mysqlConnString)
}
