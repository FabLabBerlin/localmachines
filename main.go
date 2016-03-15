package main

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/database/connect"
	_ "github.com/FabLabBerlin/localmachines/docs"
	"github.com/FabLabBerlin/localmachines/models"
	_ "github.com/FabLabBerlin/localmachines/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/toolbox"
	_ "github.com/go-sql-driver/mysql"
)

var runMode string

func init() {
	beego.AppConfig.Set("DirectoryIndex", "true")
	beego.SetStaticPath("/swagger", "swagger")

	// Config default files directory
	beego.SetStaticPath("/files", "files")
}

func main() {
	runMode = beego.AppConfig.String("RunMode")
	beego.Info("beego RunMode:", runMode)
	configClients()
	configDatabase()
	setupTasks()
	toolbox.StartTask()
	defer toolbox.StopTask()

	// Config automatic API docs
	beego.AppConfig.Set("DirectoryIndex", "true")
	beego.SetStaticPath("/swagger", "swagger")

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
		beego.SetStaticPath("/machines", "clients/machines/dev")
		beego.SetStaticPath("/admin", "clients/admin/dev")
		beego.SetStaticPath("/signup", "clients/signup/dev")
		beego.SetStaticPath("/user", "clients/user/dev")
		beego.SetStaticPath("/landing", "../localmachines-web")
	} else { // prod and any other runmode
		beego.SetStaticPath("/machines", "clients/machines/prod")
		beego.SetStaticPath("/admin", "clients/admin/prod")
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
		"0/10 * * * * *",
		models.AutoExtendUserMemberships)
	toolbox.AddTask("Extend User Memberships", extUsrMemberships)
}
