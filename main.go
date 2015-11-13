package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/toolbox"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kr15h/fabsmith/database/connect"
	_ "github.com/kr15h/fabsmith/docs"
	"github.com/kr15h/fabsmith/models"
	_ "github.com/kr15h/fabsmith/routers"
)

func main() {
	beego.Info("beego.RunMode:", beego.RunMode)

	configClients()
	configDatabase()
	setupTasks()
	toolbox.StartTask()
	defer toolbox.StopTask()

	// Config automatic API docs
	beego.DirectoryIndex = true
	beego.StaticDir["/swagger"] = "swagger"

	// Config default files directory
	beego.StaticDir["/files"] = "files"

	// Routing https
	beego.InsertFilter("/", beego.BeforeRouter, RedirectHttp) // for http://mysite

	beego.Run()

}

var RedirectHttp = func(ctx *context.Context) {
	HttpsEnabled, err := beego.AppConfig.Bool("EnableHttpTLS")
	if HttpsEnabled && err == nil {
		if !ctx.Input.IsSecure() {
			url := "https://" + ctx.Input.Domain() + ":" + beego.AppConfig.String("HttpsPort") + ctx.Input.Uri()
			ctx.Redirect(302, url)
		}
	}
}

func configClients() {

	// Allow access index.html file
	beego.DirectoryIndex = true

	beego.Trace(beego.RunMode)

	// Serve self-contained Angular JS applications depending on runmode
	if beego.RunMode == "dev" {
		beego.SetStaticPath("/machines", "clients/machines/dev")
		beego.SetStaticPath("/admin", "clients/admin/dev")
		beego.SetStaticPath("/signup", "clients/signup/dev")
		beego.SetStaticPath("/user", "clients/user/dev")
	} else { // prod and any other runmode
		beego.SetStaticPath("/machines", "clients/machines/prod")
		beego.SetStaticPath("/admin", "clients/admin/prod")
		beego.SetStaticPath("/signup", "clients/signup/prod")
		beego.SetStaticPath("/user", "clients/user/prod")
	}
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

	connect.Connect(mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDb)
}

// Setup Beego toolbox tasks. They are kind of cron jobs.
func setupTasks() {
	fmt.Println("Starting tasks")
	extUsrMemberships := toolbox.NewTask("Extend User Memberships",
		"0/10 * * * * *",
		models.AutoExtendUserMemberships)
	toolbox.AddTask("Extend User Memberships", extUsrMemberships)

	dataLogSync := toolbox.NewTask("Sync Data Log",
		"0/10 * * * * *",
		models.DataLogSync)
	toolbox.AddTask("Sync Data Log", dataLogSync)
}
