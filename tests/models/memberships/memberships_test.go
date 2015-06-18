package membershipTests

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"github.com/kr15h/fabsmith/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	// find app.conf path
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator)+".."+string(filepath.Separator)+".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	// Setting log level
	beego.SetLevel(beego.LevelError)

	// Force Runmode to "test"
	beego.RunMode = "test"

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
	beego.Info("Lel: " + mysqlDb)
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

func TestDeleteMembership(t *testing.T) {
	Convey("Testing DeleteMembership", t, func() {
		membershipName := "My awesome membership program"
		Convey("Try to delete non-existing membership", func() {
			err := models.DeleteMembership(0)

			So(err, ShouldNotBeNil)
		})
		Convey("Creating a membership and delete it", func() {
			mid, _ := models.CreateMembership(membershipName)
			err := models.DeleteMembership(mid)

			So(err, ShouldBeNil)
		})
	})
}

func TestCreateMembership(t *testing.T) {
	Convey("Testing CreateMembership", t, func() {
		membershipName := "My awesome membership"
		Convey("Creating one membership", func() {
			mid, err := models.CreateMembership(membershipName)
			defer models.DeleteMembership(mid)

			So(err, ShouldBeNil)
		})
	})
}
