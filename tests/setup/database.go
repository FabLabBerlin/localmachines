package setup

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/database/connect"
	"github.com/kr15h/fabsmith/models"
	"github.com/kr15h/fabsmith/models/invoices"
	"github.com/kr15h/fabsmith/models/purchases"
)

// ConfigDB : Configure database for tests
func ConfigDB() {
	beego.SetLevel(beego.LevelError)
	//beego.SetLevel(beego.LevelDebug)

	runmodetest, err := beego.AppConfig.Bool("runmodtest")
	if !runmodetest || err != nil {
		fmt.Println(err)
		panic("Your configuration file is wrong for testing, see app.example.conf")
	}

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

// ResetDB : Reset the database after each test
func ResetDB() {
	o := orm.NewOrm()

	// TODO: Use model.TableName() for table names

	var machines []models.Machine
	o.QueryTable("machines").All(&machines)
	for _, item := range machines {
		o.Delete(&item)
	}

	var memberships []models.Membership
	o.QueryTable("membership").All(&memberships)
	for _, item := range memberships {
		o.Delete(&item)
	}

	var invoices []invoices.Invoice
	o.QueryTable("invoices").All(&invoices)
	for _, item := range invoices {
		o.Delete(&item)
	}

	var netswitches []models.NetSwitchMapping
	o.QueryTable("netswitch").All(&netswitches)
	for _, item := range netswitches {
		o.Delete(&item)
	}

	var users []models.User
	o.QueryTable("user").All(&users)
	for _, item := range users {
		o.Delete(&item)
	}

	var user_memberships []models.UserMembership
	o.QueryTable("user_membership").All(&user_memberships)
	for _, item := range user_memberships {
		o.Delete(&item)
	}

	var permissions []models.Permission
	o.QueryTable("permission").All(&permissions)
	for _, item := range permissions {
		o.Delete(&item)
	}

	var auths []models.Auth
	o.QueryTable("auth").All(&auths)
	for _, item := range auths {
		o.Delete(&item)
	}

	var purchases []purchases.Purchase
	o.QueryTable("purchases").All(&purchases)
	for _, item := range purchases {
		o.Delete(&item)
	}

}
