package modelTest

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

// ConfigDB : Configure database for tests
func ConfigDB() {
	_, file, _, _ := runtime.Caller(1)

	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator)+".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	beego.SetLevel(beego.LevelError)

	beego.RunMode = "test"

	runmodetest, err := beego.AppConfig.Bool("runmodtest")
	if !runmodetest || err != nil {
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
	beego.Info("Lel: " + mysqlDb)
	if mysqlDb == "" {
		panic("Please set mysqldb in app.conf")
	}

	mysqlConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDb)

	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", mysqlConnString)
}

// ResetDB : Reset the database after each test
func ResetDB() {
	o := orm.NewOrm()

	var hexabuses []models.HexabusMapping
	o.QueryTable("hexaswitch").All(&hexabuses)
	for _, item := range hexabuses {
		o.Delete(&item)
	}

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

	var activations []models.Activation
	o.QueryTable("activations").All(&activations)
	for _, item := range activations {
		o.Delete(&item)
	}

}
