package main

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/database/connect"
	_ "github.com/FabLabBerlin/localmachines/docs"
	_ "github.com/FabLabBerlin/localmachines/lib/log"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/user_memberships/auto_extend"
	"github.com/FabLabBerlin/localmachines/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/toolbox"
	_ "github.com/go-sql-driver/mysql"
)

var runMode string

func main() {
	runMode = beego.AppConfig.String("RunMode")
	beego.Info("beego RunMode:", runMode)
	configClients()
	configDatabase()
	routers.Init()
	setupTasks()
	toolbox.StartTask()
	defer toolbox.StopTask()

	// Config automatic API docs
	beego.AppConfig.Set("DirectoryIndex", "true")
	if runMode == "dev" {
		beego.SetStaticPath("/swagger", "swagger")
	}

	// Config default files directory
	beego.SetStaticPath("/files", "files")

	// Routing https
	beego.InsertFilter("/", beego.BeforeRouter, RedirectHttp) // for http://mysite
	beego.Run()

}

var RedirectHttp = func(ctx *context.Context) {
	HttpsEnabled, err := beego.AppConfig.Bool("EnableHttpTLS")
	if HttpsEnabled && err == nil {
		if !ctx.Input.IsSecure() {
			url := "https://" + ctx.Input.Domain() + ":" + beego.AppConfig.String("HttpsPort") + ctx.Input.URI()
			ctx.Redirect(302, url)
		}
	}
}

func configClients() {

	// Allow access index.html file
	beego.AppConfig.Set("DirectoryIndex", "true")

	beego.Trace(runMode)

	// Serve self-contained Angular JS applications depending on runmode
	if runMode == "dev" {
		beego.SetStaticPath("/admin/assets", "clients/admin/dev/assets")
		beego.SetStaticPath("/admin/bower_components", "clients/admin/dev/bower_components")
		beego.SetStaticPath("/admin/ng-components", "clients/admin/dev/ng-components")
		beego.SetStaticPath("/admin/ng-main.js", "clients/admin/dev/ng-main.js")
		beego.SetStaticPath("/admin/ng-modules", "clients/admin/dev/ng-modules")
		beego.SetStaticPath("/machines/assets", "clients/machines/dev")
		beego.SetStaticPath("/signup", "clients/signup/dev")
		beego.SetStaticPath("/user", "clients/user/dev")
		beego.SetStaticPath("/landing", "../localmachines-web")
	} else { // prod and any other runmode
		beego.SetStaticPath("/admin/assets", "clients/admin/prod/assets")
		beego.SetStaticPath("/admin/ng-modules", "clients/admin/prod/ng-modules")
		beego.SetStaticPath("/machines/assets", "clients/machines/prod")
		beego.SetStaticPath("/signup", "clients/signup/prod")
		beego.SetStaticPath("/user", "clients/user/prod")
		beego.SetStaticPath("/landing", "../localmachines-web")
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
		"0 0/10 * * * *",
		auto_extend.RunTask)
	fetchLocalIps := toolbox.NewTask("Fetch Local IPs",
		"0 0/2 * * * *",
		machine.FetchLocalIpsTask)
	calculateTotals := toolbox.NewTask("Calculate Invoice Totals",
		" 0 0/50 * * * *",
		invutil.TaskCalculateTotals)
	fastbillSync := toolbox.NewTask("Sync Fastbill",
		" 0 0/59 * * * *",
		invutil.TaskFastbillSync)
	pingNetswitches := toolbox.NewTask("Ping Netswitches",
		" 0 0/10 * * * *",
		machine.TaskPingNetswitches)

	toolbox.AddTask("Calculate Invoice Totals", calculateTotals)
	toolbox.AddTask("Extend User Memberships", extUsrMemberships)
	toolbox.AddTask("Fetch Local IPs", fetchLocalIps)
	toolbox.AddTask("Sync Fastbill", fastbillSync)
	toolbox.AddTask("Ping Netswitches", pingNetswitches)
}
