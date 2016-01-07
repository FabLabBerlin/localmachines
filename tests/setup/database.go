package setup

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/FabLabBerlin/localmachines/database/connect"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/products"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var dbModels = []Model{
	&invoices.Invoice{},
	&models.Auth{},
	&models.Machine{},
	&models.Membership{},
	&models.NetSwitchMapping{},
	&models.Permission{},
	&models.User{},
	&products.Product{},
	&purchases.Purchase{},
	&models.UserMembership{},
}

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

type Model interface {
	TableName() string
}

// ResetDB : Reset the database after each test
func ResetDB() {
	o := orm.NewOrm()

	for _, dbModel := range dbModels {
		query := fmt.Sprintf("DELETE FROM %v", dbModel.TableName())
		if _, err := o.Raw(query).Exec(); err != nil {
			panic(err.Error())
		}
	}
}
